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

package main

import (
	"context"
	"fmt"
	"log"

	pgvector "github.com/Wood-Q/Eino-pgvector"
	"github.com/cloudwego/eino-ext/components/document/transformer/splitter/markdown"
	"github.com/cloudwego/eino-ext/components/embedding/tencentcloud"
	"github.com/cloudwego/eino/schema"
	_ "github.com/lib/pq"
)

func main() {
	ctx := context.Background()

	// 创建 embedder 配置
	cfg := &tencentcloud.EmbeddingConfig{
		SecretID:  "",
		SecretKey: "",
	}

	// 创建 embedder
	embedder, err := tencentcloud.NewEmbedder(ctx, cfg)
	if err != nil {
		panic(err)
	}

	// 创建 pgvector indexer
	indexer, err := pgvector.NewIndexer(ctx, &pgvector.IndexerConfig{
		Host:      "localhost",
		Port:      5433,
		User:      "postgres",
		Password:  "123456",
		DBName:    "vectorDB",
		SSLMode:   "disable",
		TableName: "documents",
		Dimension: 1024,
		IndexType: "hnsw",
		IndexOptions: map[string]interface{}{
			"m":               16,
			"ef_construction": 64,
		},
		Embedding: embedder,
	})
	if err != nil {
		log.Fatalf("创建indexer失败: %v", err)
	}
	defer indexer.Close()
	// 初始化分割器
	splitter, err := markdown.NewHeaderSplitter(ctx, &markdown.HeaderConfig{
		Headers: map[string]string{
			"##": "",
		},
		TrimHeaders: false,
	})
	if err != nil {
		panic(err)
	}
	markdownDoc := &schema.Document{
		ID:      "1", // 可以是任意唯一的I,
		Content: "## Title 1\nHello Word\n## Title 2\nWord Hello",
	}
	// 分割文档
	docs, err := splitter.Transform(ctx, []*schema.Document{markdownDoc})
	if err != nil {
		panic(err)
	}
	for _, doc := range docs {
		fmt.Printf("文档: %s\n,----------", doc.Content)
		fmt.Printf("元数据: %v\n\n", doc.MetaData)
	}
	// 存储文档
	ids, err := indexer.Store(ctx, docs)
	if err != nil {
		log.Fatalf("存储文档失败: %v", err)
	}

	fmt.Printf("成功存储 %d 个文档，ID: %v\n\n", len(ids), ids)
}
