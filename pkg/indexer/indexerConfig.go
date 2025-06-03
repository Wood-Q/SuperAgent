package indexer

import (
	"context"

	"MoonAgent/pkg/indexer/es8"

	"github.com/cloudwego/eino-ext/components/embedding/ark"
	"github.com/cloudwego/eino/schema"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types/enums/densevectorsimilarity"
)

const (
	indexName          = "my_rag"
	fieldContent       = "content"
	fieldContentVector = "content_vector"
)

func NewIndexerConfig(client *elasticsearch.Client, emb *ark.Embedder) *es8.IndexerConfig {
	dims := 2560
	similarity := densevectorsimilarity.Cosine
	index := true
	//创建本地mapping
	mapping := &types.TypeMapping{
		Properties: map[string]types.Property{
			"id":      types.NewTextProperty(),
			"content": types.NewTextProperty(),
			"content_dense_vector": &types.DenseVectorProperty{
				Dims:       &dims,
				Index:      &index,
				Similarity: &similarity,
			},
		},
	}
	return &es8.IndexerConfig{
		Client:       client,
		Index:        indexName,
		BatchSize:    5,
		LocalMapping: mapping,
		Embedding:    emb,
		DocumentToFields: func(ctx context.Context, doc *schema.Document) (field2Value map[string]es8.FieldValue, err error) {
			return map[string]es8.FieldValue{
				"id": {
					Value: doc.ID,
				},
				"content": {
					Value:    doc.Content,
					EmbedKey: "content_dense_vector",
				},
			}, nil
		},
		ValidationMode:    es8.ValidationModeWarn,
		EnableSchemaCheck: true,
	}
}
