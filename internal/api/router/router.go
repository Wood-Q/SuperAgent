package router

import (
	"MoonAgent/cmd/di"
	"MoonAgent/internal/api/handler"
	"context"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
)

// CORS中间件
func corsMiddleware() app.HandlerFunc {
	return func(c context.Context, ctx *app.RequestContext) {
		ctx.Header("Access-Control-Allow-Origin", "*")
		ctx.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		ctx.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if string(ctx.Method()) == "OPTIONS" {
			ctx.SetStatusCode(consts.StatusOK)
			return
		}

		ctx.Next(c)
	}
}

func InitRouter(h *server.Hertz, app *di.Application) {
	// 添加CORS中间件
	h.Use(corsMiddleware())

	ChatHandler := handler.NewChatHandler(app)
	v1 := h.Group("/api")
	v1.POST("/chat", ChatHandler.ChatWithModel)
	v1.POST("/chat/stream", ChatHandler.StreamChatWithModel)
}
