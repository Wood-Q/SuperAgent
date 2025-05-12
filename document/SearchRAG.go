package document

import (
	"SuperAgent/global"
	"context"

	"github.com/cloudwego/eino/schema"
	"go.uber.org/zap"
)

func SearchRAG(query string, ctx context.Context) []*schema.Document {

	// Retrieve documents
	documents, err := global.Retriever.Retrieve(ctx, query)
	if err != nil {
		zap.S().Error("Failed to retrieve: %v", zap.String("error", err.Error()))
	}
	return documents
}
