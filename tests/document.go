package main

import (
	"MoonAgent/cmd/di"
	"context"
	"fmt"

	"github.com/cloudwego/eino/schema"
)

func main() {
	ctx := context.Background()
	app, clear, err := di.InitializeApplication()
	if err != nil {
		panic(err)
	}
	defer clear()

	// content, err := os.ReadFile("../assets/documents/muelsyse copy.txt")
	// if err != nil {
	// 	panic(err)
	// }
	// document := []*schema.Document{
	// 	{
	// 		ID:      "muelsyse",
	// 		Content: string(content),
	// 	},
	// }
	// docs, err := splitter.SplitDocs(ctx, app.Embedder, document)
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Println(docs)
	documents := []*schema.Document{
		{
			ID:       "test",
			Content:  "test",
			MetaData: map[string]any{"source": "test"},
		},
	}
	ids, err := app.Indexer.Store(ctx, documents)
	if err != nil {
		panic(err)
	}
	fmt.Println(ids)
}
