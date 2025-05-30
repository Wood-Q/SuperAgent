package main

import (
	"MoonAgent/cmd/di"
	"MoonAgent/internal/api/router"

	"github.com/cloudwego/hertz/pkg/app/server"
)

func main() {
	app, clear, err := di.InitializeApplication()
	if err != nil {
		panic(err)
	}
	defer clear()
	H := server.Default()
	router.InitRouter(H, app)
	H.Spin()
}
