package main

import (
	"MoonAgent/internal/config"
	"MoonAgent/internal/pipeline"
	"context"
	"fmt"

	"github.com/cloudwego/eino/schema"
)

func main() {
	config.InitConfig("../configs")
	runnable, err := pipeline.BuildAssitant(context.Background())
	if err != nil {
		panic(err)
	}
	out, err := runnable.Invoke(context.Background(), map[string]any{
		"messages": []schema.Message{
			{
				Role:    "user",
				Content: "你好",
			},
		},
	})
	if err != nil {
		panic(err)
	}
	fmt.Println(out.Content)
}
