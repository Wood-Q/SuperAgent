/*
 * Copyright 2024 CloudWeGo Authors
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package es8

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/cloudwego/eino/callbacks"
	"github.com/cloudwego/eino/components"
	"github.com/cloudwego/eino/components/embedding"
	"github.com/cloudwego/eino/components/indexer"
	"github.com/cloudwego/eino/schema"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esutil"
	"github.com/elastic/go-elasticsearch/v8/typedapi/indices/getmapping"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types"
	densevectorsimilarity "github.com/elastic/go-elasticsearch/v8/typedapi/types/enums/densevectorsimilarity"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types/enums/dynamicmapping"
)

const defaultBatchSize = 5
const typ = "es8_indexer"

// MappingValidationMode defines how to handle mapping validation failures
type MappingValidationMode int

const (
	// ValidationModeError raises error when validation fails
	ValidationModeError MappingValidationMode = iota
	// ValidationModeWarn logs warning when validation fails
	ValidationModeWarn
	// ValidationModeSkip skips validation entirely
	ValidationModeSkip
)

type IndexerConfig struct {
	Client *elasticsearch.Client `json:"client"`

	Index string `json:"index"`
	// BatchSize controls max texts size for embedding.
	// Default is 5.
	BatchSize int `json:"batch_size"`
	// FieldMapping supports customize es fields from eino document.
	// Each key - FieldValue.Value from field2Value will be saved, and
	// vector of FieldValue.Value will be saved if FieldValue.EmbedKey is not empty.
	DocumentToFields func(ctx context.Context, doc *schema.Document) (field2Value map[string]FieldValue, err error)
	// Embedding vectorization method, must provide in two cases
	// 1. VectorFields contains fields except doc Content
	// 2. VectorFields contains doc Content and vector not provided in doc extra (see Document.Vector method)
	Embedding embedding.Embedder

	// LocalMapping defines the expected index structure using official types.TypeMapping
	LocalMapping *types.TypeMapping `json:"local_mapping"`
	// Dynamic setting for the index mapping (optional)
	Dynamic *dynamicmapping.DynamicMapping `json:"dynamic"`
	// ValidationMode controls how mapping validation failures are handled
	// Default is ValidationModeError
	ValidationMode MappingValidationMode `json:"validation_mode"`
	// EnableSchemaCheck enables document schema validation before indexing
	EnableSchemaCheck bool `json:"enable_schema_check"`
}

type FieldValue struct {
	// Value original Value
	Value any
	// EmbedKey if set, Value will be vectorized and saved to es.
	// If Stringify method is provided, Embedding input text will be Stringify(Value).
	// If Stringify method not set, retriever will try to assert Value as string.
	EmbedKey string
	// Stringify converts Value to string
	Stringify func(val any) (string, error)
}

type Indexer struct {
	client *elasticsearch.Client
	config *IndexerConfig
}

// getDefaultMapping returns the default index mapping using types.TypeMapping
func getDefaultMapping() *types.TypeMapping {
	dims := 2560
	index := true
	similarity := densevectorsimilarity.Cosine

	return &types.TypeMapping{
		Properties: map[string]types.Property{
			"content":        types.NewTextProperty(),
			"extra_location": types.NewTextProperty(),
			"content_dense_vector": &types.DenseVectorProperty{
				Dims:       &dims,
				Index:      &index,
				Similarity: &similarity,
			},
		},
	}
}

// getRemoteMapping retrieves mapping from Elasticsearch for the specified index
func (i *Indexer) getRemoteMapping(ctx context.Context) (map[string]types.Property, *dynamicmapping.DynamicMapping, error) {
	indexName := i.config.Index

	req := getmapping.New(i.client.Transport)
	req.Index(indexName)
	res, err := req.Do(ctx)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get mapping for index %s: %w", indexName, err)
	}

	indexMapping, exists := res[indexName]
	if !exists {
		return nil, nil, fmt.Errorf("index %s not found in mapping response", indexName)
	}

	properties := indexMapping.Mappings.Properties
	dynamic := indexMapping.Mappings.Dynamic

	fmt.Printf("远程索引 %s mapping 获取成功，字段数量: %d\n", indexName, len(properties))
	if dynamic != nil {
		fmt.Printf("远程索引 dynamic 设置: %v\n", *dynamic)
	}

	return properties, dynamic, nil
}

// ensureIndex ensures the index exists, creates it if not, and validates mapping consistency
func (i *Indexer) ensureIndex(ctx context.Context) error {
	indexName := i.config.Index

	// 1. Check if index exists
	res, err := i.client.Indices.Exists([]string{indexName})
	if err != nil {
		return fmt.Errorf("failed to check if index exists: %w", err)
	}

	if res.StatusCode == 404 {
		// 2. Index doesn't exist, create it
		fmt.Printf("索引 %s 不存在，正在创建...\n", indexName)
		return i.createIndexWithLocalMapping(ctx)
	}

	if res.StatusCode != 200 {
		return fmt.Errorf("unexpected response when checking index existence: %s", res.String())
	}

	// 3. Index exists, validate mapping if needed
	fmt.Printf("索引 %s 已存在\n", indexName)

	if i.config.ValidationMode != ValidationModeSkip && i.config.LocalMapping != nil {
		if err := i.validateMappingConsistency(ctx); err != nil {
			switch i.config.ValidationMode {
			case ValidationModeError:
				return fmt.Errorf("mapping validation failed: %w", err)
			case ValidationModeWarn:
				log.Printf("警告: mapping validation failed: %v", err)
			}
		} else {
			fmt.Printf("✓ 索引 %s mapping 验证通过\n", indexName)
		}
	}

	return nil
}

// createIndexWithLocalMapping creates index using local mapping definition
func (i *Indexer) createIndexWithLocalMapping(ctx context.Context) error {
	indexName := i.config.Index

	// Use default mapping if no local mapping is provided
	mapping := i.config.LocalMapping
	if mapping == nil {
		mapping = getDefaultMapping()
		fmt.Println("使用默认 mapping 创建索引")
	}

	// Prepare index settings
	indexBody := map[string]interface{}{
		"mappings": map[string]interface{}{
			"properties": mapping.Properties,
		},
	}

	// Add dynamic setting if configured
	if i.config.Dynamic != nil {
		indexBody["mappings"].(map[string]interface{})["dynamic"] = *i.config.Dynamic
		fmt.Printf("设置 dynamic 模式: %v\n", *i.config.Dynamic)
	}

	// Serialize to JSON
	jsonBody, err := json.Marshal(indexBody)
	if err != nil {
		return fmt.Errorf("failed to marshal index mapping: %w", err)
	}

	// Create index
	createRes, err := i.client.Indices.Create(
		indexName,
		i.client.Indices.Create.WithBody(strings.NewReader(string(jsonBody))),
		i.client.Indices.Create.WithContext(ctx),
	)
	if err != nil {
		return fmt.Errorf("failed to create index %s: %w", indexName, err)
	}
	defer createRes.Body.Close()

	if createRes.IsError() {
		var createError map[string]interface{}
		if err := json.NewDecoder(createRes.Body).Decode(&createError); err == nil {
			return fmt.Errorf("failed to create index %s: %v", indexName, createError)
		}
		return fmt.Errorf("failed to create index %s: %s", indexName, createRes.String())
	}

	fmt.Printf("✓ 成功创建索引 %s\n", indexName)
	return nil
}

// validateMappingConsistency validates consistency between local and remote mappings
func (i *Indexer) validateMappingConsistency(ctx context.Context) error {
	remoteProperties, remoteDynamic, err := i.getRemoteMapping(ctx)
	if err != nil {
		return fmt.Errorf("failed to get remote mapping: %w", err)
	}

	localProperties := i.config.LocalMapping.Properties

	fmt.Printf("开始对比 mapping 一致性...\n")
	fmt.Printf("本地字段数量: %d, 远程字段数量: %d\n", len(localProperties), len(remoteProperties))

	// Validate field consistency
	if err := i.validateFieldConsistency(localProperties, remoteProperties); err != nil {
		return err
	}

	// Validate dynamic setting consistency
	if err := i.validateDynamicConsistency(remoteDynamic); err != nil {
		return err
	}

	return nil
}

// validateFieldConsistency compares local and remote field definitions
func (i *Indexer) validateFieldConsistency(localProperties, remoteProperties map[string]types.Property) error {
	// Check if all local fields exist in remote mapping with correct types
	for fieldName, localProp := range localProperties {
		remoteProp, exists := remoteProperties[fieldName]
		if !exists {
			return fmt.Errorf("字段 '%s' 在远程索引中不存在", fieldName)
		}

		if err := i.comparePropertyTypes(fieldName, localProp, remoteProp); err != nil {
			return err
		}
	}

	// Check for extra fields in remote mapping (warning only)
	for fieldName := range remoteProperties {
		if _, exists := localProperties[fieldName]; !exists {
			fmt.Printf("警告: 远程索引包含本地未定义的字段 '%s'\n", fieldName)
		}
	}

	return nil
}

// comparePropertyTypes compares property types between local and remote
func (i *Indexer) comparePropertyTypes(fieldName string, localProp, remoteProp types.Property) error {
	// Handle different property types
	switch local := localProp.(type) {
	case *types.TextProperty:
		if _, ok := remoteProp.(*types.TextProperty); !ok {
			return fmt.Errorf("字段 '%s' 类型不匹配: 期望 text, 远程为其他类型", fieldName)
		}
	case *types.DenseVectorProperty:
		remote, ok := remoteProp.(*types.DenseVectorProperty)
		if !ok {
			return fmt.Errorf("字段 '%s' 类型不匹配: 期望 dense_vector, 远程为其他类型", fieldName)
		}
		// Compare dimensions
		if local.Dims != nil && remote.Dims != nil && *local.Dims != *remote.Dims {
			return fmt.Errorf("字段 '%s' 向量维度不匹配: 期望 %d, 远程为 %d",
				fieldName, *local.Dims, *remote.Dims)
		}
		// Compare similarity
		if local.Similarity != nil && remote.Similarity != nil && *local.Similarity != *remote.Similarity {
			return fmt.Errorf("字段 '%s' 相似度算法不匹配: 期望 %v, 远程为 %v",
				fieldName, *local.Similarity, *remote.Similarity)
		}
	case *types.ObjectProperty:
		remote, ok := remoteProp.(*types.ObjectProperty)
		if !ok {
			return fmt.Errorf("字段 '%s' 类型不匹配: 期望 object, 远程为其他类型", fieldName)
		}
		// Recursively validate nested properties
		if local.Properties != nil && remote.Properties != nil {
			for nestedName, nestedLocal := range local.Properties {
				if nestedRemote, exists := remote.Properties[nestedName]; exists {
					if err := i.comparePropertyTypes(fmt.Sprintf("%s.%s", fieldName, nestedName),
						nestedLocal, nestedRemote); err != nil {
						return err
					}
				} else {
					return fmt.Errorf("嵌套字段 '%s.%s' 在远程索引中不存在", fieldName, nestedName)
				}
			}
		}
	case *types.LongNumberProperty:
		if _, ok := remoteProp.(*types.LongNumberProperty); !ok {
			return fmt.Errorf("字段 '%s' 类型不匹配: 期望 long, 远程为其他类型", fieldName)
		}
	case *types.DoubleNumberProperty:
		if _, ok := remoteProp.(*types.DoubleNumberProperty); !ok {
			return fmt.Errorf("字段 '%s' 类型不匹配: 期望 double, 远程为其他类型", fieldName)
		}
	case *types.BooleanProperty:
		if _, ok := remoteProp.(*types.BooleanProperty); !ok {
			return fmt.Errorf("字段 '%s' 类型不匹配: 期望 boolean, 远程为其他类型", fieldName)
		}
	// Add more property types as needed
	default:
		fmt.Printf("警告: 字段 '%s' 的类型验证尚未实现\n", fieldName)
	}

	return nil
}

// validateDynamicConsistency validates dynamic setting consistency
func (i *Indexer) validateDynamicConsistency(remoteDynamic *dynamicmapping.DynamicMapping) error {
	if i.config.Dynamic == nil {
		// Local doesn't specify dynamic, skip validation
		return nil
	}

	if remoteDynamic == nil {
		return fmt.Errorf("本地配置了 dynamic 设置为 %v，但远程索引未设置 dynamic", *i.config.Dynamic)
	}

	if *i.config.Dynamic != *remoteDynamic {
		return fmt.Errorf("dynamic 设置不一致: 期望 %v, 远程为 %v", *i.config.Dynamic, *remoteDynamic)
	}

	fmt.Printf("✓ dynamic 设置验证通过: %v\n", *i.config.Dynamic)
	return nil
}

// validateDocumentSchema validates document against expected schema (if enabled)
func (i *Indexer) validateDocumentSchema(ctx context.Context, doc *schema.Document) error {
	if !i.config.EnableSchemaCheck || i.config.LocalMapping == nil {
		return nil
	}

	fields, err := i.config.DocumentToFields(ctx, doc)
	if err != nil {
		return fmt.Errorf("failed to get document fields: %w", err)
	}

	expectedFields := i.config.LocalMapping.Properties

	// Check for missing required fields
	for expectedField := range expectedFields {
		found := false
		for _, fieldValue := range fields {
			if fieldValue.EmbedKey == expectedField {
				found = true
				break
			}
		}
		if !found {
			if _, directExists := fields[expectedField]; !directExists {
				fmt.Printf("警告: 文档 %s 缺少期望字段 '%s'\n", doc.ID, expectedField)
			}
		}
	}

	// Check for unexpected fields
	for fieldName := range fields {
		if _, expected := expectedFields[fieldName]; !expected {
			fmt.Printf("警告: 文档 %s 包含未定义字段 '%s'\n", doc.ID, fieldName)
		}
	}

	return nil
}

func NewIndexer(ctx context.Context, conf *IndexerConfig) (*Indexer, error) {
	if conf.Client == nil {
		return nil, fmt.Errorf("[NewIndexer] es client not provided")
	}

	if conf.DocumentToFields == nil {
		return nil, fmt.Errorf("[NewIndexer] DocumentToFields method not provided")
	}

	if conf.Index == "" {
		return nil, fmt.Errorf("[NewIndexer] index name cannot be empty")
	}

	if conf.BatchSize <= 0 {
		conf.BatchSize = defaultBatchSize
	}

	// Validate batch size is reasonable
	if conf.BatchSize > 1000 {
		return nil, fmt.Errorf("[NewIndexer] batch size %d is too large, maximum recommended is 1000", conf.BatchSize)
	}

	// Set default validation mode
	if conf.ValidationMode == 0 {
		conf.ValidationMode = ValidationModeError
	}

	// Test connection to Elasticsearch
	if _, err := conf.Client.Info(); err != nil {
		return nil, fmt.Errorf("[NewIndexer] failed to connect to Elasticsearch: %w", err)
	}

	indexer := &Indexer{
		client: conf.Client,
		config: conf,
	}

	// Ensure index exists and validate mapping
	if err := indexer.ensureIndex(ctx); err != nil {
		return nil, fmt.Errorf("[NewIndexer] index setup failed: %w", err)
	}

	return indexer, nil
}

func (i *Indexer) Store(ctx context.Context, docs []*schema.Document, opts ...indexer.Option) (ids []string, err error) {
	ctx = callbacks.EnsureRunInfo(ctx, i.GetType(), components.ComponentOfIndexer)
	ctx = callbacks.OnStart(ctx, &indexer.CallbackInput{Docs: docs})
	defer func() {
		if err != nil {
			callbacks.OnError(ctx, err)
		}
	}()

	options := indexer.GetCommonOptions(&indexer.Options{
		Embedding: i.config.Embedding,
	}, opts...)

	// Validate documents against schema if enabled
	if i.config.EnableSchemaCheck {
		for _, doc := range docs {
			if err := i.validateDocumentSchema(ctx, doc); err != nil {
				return nil, fmt.Errorf("document schema validation failed for %s: %w", doc.ID, err)
			}
		}
	}

	if err = i.bulkAdd(ctx, docs, options); err != nil {
		return nil, err
	}

	ids = iter(docs, func(t *schema.Document) string { return t.ID })

	callbacks.OnEnd(ctx, &indexer.CallbackOutput{IDs: ids})

	return ids, nil
}

func (i *Indexer) bulkAdd(ctx context.Context, docs []*schema.Document, options *indexer.Options) error {
	emb := options.Embedding
	bi, err := esutil.NewBulkIndexer(esutil.BulkIndexerConfig{
		Index:         i.config.Index,
		Client:        i.client,
		NumWorkers:    min(4, len(docs)),
		FlushBytes:    int(5e+6),
		FlushInterval: 30 * time.Second,
	})
	if err != nil {
		return err
	}

	var (
		tuples []tuple
		texts  []string
	)

	embAndAdd := func() error {
		var vectors [][]float64

		if len(texts) > 0 {
			if emb == nil {
				return fmt.Errorf("[bulkAdd] embedding method not provided")
			}

			vectors, err = emb.EmbedStrings(i.makeEmbeddingCtx(ctx, emb), texts)
			if err != nil {
				return fmt.Errorf("[bulkAdd] embedding failed, %w", err)
			}

			if len(vectors) != len(texts) {
				return fmt.Errorf("[bulkAdd] invalid vector length, expected=%d, got=%d", len(texts), len(vectors))
			}
		}

		for _, t := range tuples {
			fields := make(map[string]any)

			// Copy original fields
			for k, v := range t.fields {
				fields[k] = v
			}

			// Add vector fields
			for k, idx := range t.key2Idx {
				if idx < len(vectors) {
					fields[k] = vectors[idx]
				}
			}

			b, err := json.Marshal(fields)
			if err != nil {
				return fmt.Errorf("[bulkAdd] marshal bulk item failed, %w", err)
			}

			if err = bi.Add(ctx, esutil.BulkIndexerItem{
				Index:      i.config.Index,
				Action:     "index",
				DocumentID: t.id,
				Body:       bytes.NewReader(b),
				OnSuccess: func(ctx context.Context, item esutil.BulkIndexerItem, res esutil.BulkIndexerResponseItem) {
					fmt.Printf("✓ 文档 %s 索引成功\n", item.DocumentID)
				},
				OnFailure: func(ctx context.Context, item esutil.BulkIndexerItem, res esutil.BulkIndexerResponseItem, err error) {
					if err != nil {
						fmt.Printf("✗ 文档 %s 索引失败: %v\n", item.DocumentID, err)
					} else {
						fmt.Printf("✗ 文档 %s 索引失败: %s\n", item.DocumentID, res.Error.Reason)
					}
				},
			}); err != nil {
				return fmt.Errorf("[bulkAdd] failed to add document %s to bulk indexer: %w", t.id, err)
			}
		}

		tuples = tuples[:0]
		texts = texts[:0]

		return nil
	}

	for idx := range docs {
		doc := docs[idx]
		fields, err := i.config.DocumentToFields(ctx, doc)
		if err != nil {
			return fmt.Errorf("[bulkAdd] FieldMapping failed, %w", err)
		}

		rawFields := make(map[string]any)
		embSize := 0
		embedKeys := make(map[string]bool) // Track embed keys to prevent duplicates

		for k, v := range fields {
			rawFields[k] = v.Value
			if v.EmbedKey != "" {
				// Check if EmbedKey conflicts with existing field names
				if _, found := fields[v.EmbedKey]; found {
					return fmt.Errorf("[bulkAdd] embed key '%s' conflicts with existing field for document %s", v.EmbedKey, doc.ID)
				}

				// Check if EmbedKey conflicts with other embed keys
				if embedKeys[v.EmbedKey] {
					return fmt.Errorf("[bulkAdd] duplicate embed key '%s' found for document %s", v.EmbedKey, doc.ID)
				}
				embedKeys[v.EmbedKey] = true

				// Check if EmbedKey would overwrite a raw field
				if _, found := rawFields[v.EmbedKey]; found {
					return fmt.Errorf("[bulkAdd] embed key '%s' would overwrite existing field value for document %s", v.EmbedKey, doc.ID)
				}

				embSize++
			}
		}

		if embSize > i.config.BatchSize {
			return fmt.Errorf("[bulkAdd] needEmbeddingFields length over batch size, batch size=%d, got size=%d",
				i.config.BatchSize, embSize)
		}

		if len(texts)+embSize > i.config.BatchSize {
			if err = embAndAdd(); err != nil {
				return err
			}
		}

		key2Idx := make(map[string]int, embSize)
		for k, v := range fields {
			if v.EmbedKey != "" {
				if _, found := fields[v.EmbedKey]; found {
					return fmt.Errorf("[bulkAdd] duplicate key for origin key, key=%s", k)
				}

				if _, found := key2Idx[v.EmbedKey]; found {
					return fmt.Errorf("[bulkAdd] duplicate key from embed_key, key=%s", v.EmbedKey)
				}

				var text string
				if v.Stringify != nil {
					text, err = v.Stringify(v.Value)
					if err != nil {
						return err
					}
				} else {
					var ok bool
					text, ok = v.Value.(string)
					if !ok {
						return fmt.Errorf("[bulkAdd] assert value as string failed, key=%s, emb_key=%s", k, v.EmbedKey)
					}
				}

				key2Idx[v.EmbedKey] = len(texts)
				texts = append(texts, text)
			}
		}

		tuples = append(tuples, tuple{
			id:      doc.ID,
			fields:  rawFields,
			key2Idx: key2Idx,
		})
	}

	if len(tuples) > 0 {
		if err = embAndAdd(); err != nil {
			return err
		}
	}

	// Close bulk indexer and check statistics
	if err := bi.Close(ctx); err != nil {
		return fmt.Errorf("[bulkAdd] failed to close bulk indexer: %w", err)
	}

	// Check bulk indexer statistics
	stats := bi.Stats()
	fmt.Printf("批量索引完成: 成功=%d, 失败=%d\n", stats.NumIndexed, stats.NumFailed)

	if stats.NumFailed > 0 {
		return fmt.Errorf("[bulkAdd] %d documents failed to index", stats.NumFailed)
	}

	return nil
}

func (i *Indexer) makeEmbeddingCtx(ctx context.Context, emb embedding.Embedder) context.Context {
	runInfo := &callbacks.RunInfo{
		Component: components.ComponentOfEmbedding,
	}

	if embType, ok := components.GetType(emb); ok {
		runInfo.Type = embType
	}

	runInfo.Name = runInfo.Type + string(runInfo.Component)

	return callbacks.ReuseHandlers(ctx, runInfo)
}

func (i *Indexer) GetType() string {
	return typ
}

func (i *Indexer) IsCallbacksEnabled() bool {
	return true
}

type tuple struct {
	id      string
	fields  map[string]any
	key2Idx map[string]int
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func iter(docs []*schema.Document, fn func(*schema.Document) string) []string {
	result := make([]string, len(docs))
	for i, doc := range docs {
		result[i] = fn(doc)
	}
	return result
}
