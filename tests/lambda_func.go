package tests

import (
	"context"

	"github.com/cloudwego/eino/schema"
)

// newLambda1 component initialization function of node 'Lambda5' in graph 'Assitant'
func newLambda1(ctx context.Context, input []*schema.Document) (output map[string]any, err error) {
	output = make(map[string]any)
	contentString := ""
	for _, doc := range input {
		contentString += doc.Content
	}
	output["retrieve_result"] = contentString
	return output, nil
}
