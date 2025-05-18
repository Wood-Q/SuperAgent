package main

import (
	"MoonAgent/internal/pipeline"
	"MoonAgent/pkg/config"
	"context"
	"fmt"

	"github.com/cloudwego/eino/schema"
)

func main() {
	config, err := config.NewConfig("../configs")
	if err != nil {
		panic(err)
	}
	runnable, err := pipeline.BuildAssitant(context.Background(), config)
	if err != nil {
		panic(err)
	}
	out, err := runnable.Invoke(context.Background(), []*schema.Message{
		{
			Role:    "user",
			Content: "你好，可以帮我搜索一下明日方舟缪尔赛思，并且跳到prts对应的页面吗",
		},
	})
	if err != nil {
		panic(err)
	}
	fmt.Println(out.Content)
}
