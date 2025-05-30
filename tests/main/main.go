package main

import (
	"MoonAgent/cmd/di"
	"MoonAgent/internal/pipeline"
	"context"
	"fmt"
)

func main() {
	app, clear, err := di.InitializeApplication()
	if err != nil {
		panic(err)
	}
	defer clear()
	userInput := "你好，可以告诉我缪尔赛思的源石技艺适应性是什么吗，并且跳转到prts缪尔赛思相关页面"
	ctx := context.WithValue(context.Background(), "user_input", userInput)

	runnable, err := pipeline.BuildAssitant(ctx, app)
	if err != nil {
		panic(err)
	}
	out, err := runnable.Stream(ctx, userInput)
	if err != nil {
		panic(err)
	}
	for {
		chunk, err := out.Recv()
		if err != nil {
			panic(err)
		}
		fmt.Println(chunk)
	}
}
