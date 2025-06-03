package main

import (
	"MoonAgent/cmd/di"
	"MoonAgent/pkg/splitter"
	"context"
	"fmt"
	"os"

	"github.com/cloudwego/eino/schema"
)

func main() {
	ctx := context.Background()
	app, clear, err := di.InitializeApplication()
	if err != nil {
		panic(err)
	}
	defer clear()
	content, err := os.ReadFile("../../assets/documents/muelsyse copy.txt")
	if err != nil {
		panic(err)
	}
	document := []*schema.Document{
		{
			ID:      "muelsyse",
			Content: string(content),
		},
	}
	docs, err := splitter.SplitDocs(ctx, app.Embedder, document)
	if err != nil {
		panic(err)
	}
	_, err = app.Indexer.Store(ctx, docs)
	if err != nil {
		panic(err)
	}
	fmt.Println("Document stored successfully")

}
