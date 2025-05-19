package tests

import (
	"MoonAgent/cmd/di"
	"MoonAgent/pkg/config"
	"context"
	"fmt"
)

func main() {
	config, err := config.NewConfig("../configs")
	if err != nil {
		panic(err)
	}
	app, clear, err := di.InitializeApplication()
	if err != nil {
		panic(err)
	}
	defer clear()
	// if err != nil {
	// 	panic(err)
	// }
	// runnable, err := pipeline.BuildAssitant(context.Background(), config)
	// if err != nil {
	// 	panic(err)
	// }
	// out, err := runnable.Invoke(context.Background(), []*schema.Message{
	// 	{
	// 		Role:    "user",
	// 		Content: "你好，可以为我搜索原神相关网站吗",
	// 	},
	// })
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Println(out.Content)

	runnable, err := BuildAssitant(context.Background(), app, config)
	if err != nil {
		panic(err)
	}
	out, err := runnable.Invoke(context.Background(), "你好，可以告诉我缪尔赛思是一个怎样的人吗，塞雷娅对她来说是什么样的存在")
	if err != nil {
		panic(err)
	}
	fmt.Println(out.Content)
}
