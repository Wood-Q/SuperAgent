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
	content, err := os.ReadFile("../../assets/documents/muelsyse.txt")
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
	batchsize := 10
	for i := 0; i < len(docs); i += batchsize {
		end := i + batchsize
		if end > len(docs) {
			end = len(docs)
		}
		batch := docs[i:end]
		_, err = app.Indexer.Store(ctx, batch)
		if err != nil {
			panic(err)
		}
	fmt.Println("Document stored successfully")
	}
}
