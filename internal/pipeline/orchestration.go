package pipeline

import (
	"context"

	"MoonAgent/pkg/config"

	"github.com/cloudwego/eino/compose"
	"github.com/cloudwego/eino/schema"
)

func BuildAssitant(ctx context.Context, config *config.ServerConfig) (r compose.Runnable[[]*schema.Message, *schema.Message], err error) {
	const Lambda3 = "Lambda3"
	g := compose.NewGraph[[]*schema.Message, *schema.Message]()
	lambda3KeyOfLambda, err := newLambda(ctx, config)
	if err != nil {
		return nil, err
	}
	_ = g.AddLambdaNode(Lambda3, lambda3KeyOfLambda)
	_ = g.AddEdge(compose.START, Lambda3)
	_ = g.AddEdge(Lambda3, compose.END)
	r, err = g.Compile(ctx, compose.WithGraphName("Assitant"), compose.WithNodeTriggerMode(compose.AnyPredecessor))
	if err != nil {
		return nil, err
	}
	return r, err
}
