package router

import (
	"SuperAgent/api"

	"github.com/cloudwego/hertz/pkg/app/server"
)

func InitRouter(rg *server.Hertz) {
	apiGroup := rg.Group("/api")
	modelGroup := apiGroup.Group("/chat")
	{
		modelGroup.POST("/send", api.SendChat)
		modelGroup.POST("/health", api.DoChatWithSchemaJSON)
	}
}
