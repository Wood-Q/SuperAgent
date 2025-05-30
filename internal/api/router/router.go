package router

import (
	"MoonAgent/cmd/di"
	"MoonAgent/internal/api/handler"

	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/hertz-contrib/cors"
)

func InitRouter(h *server.Hertz, app *di.Application) {
	// 添加CORS中间件
	h.Use(cors.New(
		cors.Config{
			AllowAllOrigins:  true,
			AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
			AllowHeaders:     []string{"Content-Type", "Authorization", "X-Requested-With", "Accept", "Origin"},
			ExposeHeaders:    []string{"Content-Length", "Access-Control-Allow-Origin", "Access-Control-Allow-Headers"},
			AllowCredentials: true,
		},
	))

	ChatHandler := handler.NewChatHandler(app)
	v1 := h.Group("/api")
	v1.POST("/chat", ChatHandler.ChatWithModel)
	v1.POST("/chat/stream", ChatHandler.StreamChatWithModel)
}
