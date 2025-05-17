package pipeline

import (
	"context"

	"github.com/cloudwego/eino/compose"
)

func BuildAssitant(ctx context.Context) (r compose.Runnable[string, string], err error) {
	const (
		ChatTemplate1 = "ChatTemplate1"
		Lambda1       = "Lambda1"
	)
	g := compose.NewGraph[string, string](compose.WithGenLocalState(func(ctx context.Context) (state any) {
		panic("implement me")
	}))
	chatTemplate1KeyOfChatTemplate, err := newChatTemplate(ctx)
	if err != nil {
		return nil, err
	}
	_ = g.AddChatTemplateNode(ChatTemplate1, chatTemplate1KeyOfChatTemplate)
	lambda1KeyOfLambda, err := newLambda(ctx)
	if err != nil {
		return nil, err
	}
	_ = g.AddLambdaNode(Lambda1, lambda1KeyOfLambda)
	_ = g.AddEdge(compose.START, ChatTemplate1)
	_ = g.AddEdge(Lambda1, compose.END)
	_ = g.AddEdge(ChatTemplate1, Lambda1)
	r, err = g.Compile(ctx, compose.WithGraphName("Assitant"), compose.WithNodeTriggerMode(compose.AnyPredecessor))
	if err != nil {
		return nil, err
	}
	return r, err
}
