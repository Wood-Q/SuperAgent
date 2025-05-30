package router

import (
	"MoonAgent/cmd/di"
	"MoonAgent/internal/api/handler"

	"github.com/cloudwego/hertz/pkg/app/server"
)

func InitRouter(h *server.Hertz, app *di.Application) {
	ChatHandler := handler.NewChatHandler(app)
	v1 := h.Group("/api")
	v1.POST("/chat", ChatHandler.ChatWithModel)
}
