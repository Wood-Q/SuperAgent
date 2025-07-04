package pipeline

import (
	"context"

	"github.com/cloudwego/eino/compose"

	"MoonAgent/cmd/di"

	"github.com/cloudwego/eino/schema"
)

func BuildAssitant(ctx context.Context, app *di.Application) (r compose.Runnable[string, *schema.Message], err error) {
	const (
		Lambda3       = "Lambda3"
		ChatTemplate2 = "ChatTemplate2"
		Retriever4    = "Retriever4"
		Lambda5       = "Lambda5"
	)
	// 构建图
	g := compose.NewGraph[string, *schema.Message]()
	// 构建Lambda3
	lambda3KeyOfLambda, err := newLambda(ctx, app)
	if err != nil {
		return nil, err
	}
	_ = g.AddLambdaNode(Lambda3, lambda3KeyOfLambda)
	// 构建ChatTemplate2
	chatTemplate2KeyOfChatTemplate, err := newChatTemplate(ctx)
	if err != nil {
		return nil, err
	}
	_ = g.AddChatTemplateNode(ChatTemplate2, chatTemplate2KeyOfChatTemplate)
	// 构建Retriever4
	retriever4KeyOfRetriever := app.Retriever
	if err != nil {
		return nil, err
	}
	_ = g.AddRetrieverNode(Retriever4, retriever4KeyOfRetriever)
	_ = g.AddLambdaNode(Lambda5, compose.InvokableLambda(newLambda1))
	_ = g.AddEdge(compose.START, Retriever4)
	_ = g.AddEdge(Lambda3, compose.END)
	_ = g.AddEdge(ChatTemplate2, Lambda3)
	_ = g.AddEdge(Lambda5, ChatTemplate2)
	_ = g.AddEdge(Retriever4, Lambda5)
	r, err = g.Compile(ctx, compose.WithGraphName("Assitant"), compose.WithNodeTriggerMode(compose.AnyPredecessor))
	if err != nil {
		return nil, err
	}
	return r, err
}
