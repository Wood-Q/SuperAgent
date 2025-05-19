package main

import (
	"MoonAgent/cmd/di"
	"MoonAgent/internal/pipeline"
	"MoonAgent/pkg/config"
	"context"
	"fmt"
)

func main() {
	cfg, err := config.NewConfig("../configs")
	if err != nil {
		panic(err)
	}
	app, clear, err := di.InitializeApplication()
	if err != nil {
		panic(err)
	}
	defer clear()
	userInput := "你好，可以告诉我缪尔赛思的源石技艺适应性是什么吗"
	ctx := context.WithValue(context.Background(), "user_input", userInput)

	runnable, err := pipeline.BuildAssitant(ctx, app, cfg)
	if err != nil {
		panic(err)
	}
	out, err := runnable.Invoke(ctx, userInput)
	if err != nil {
		panic(err)
	}
	fmt.Println(out.Content)
}
