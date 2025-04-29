package api

import (
	"SuperAgent/global"
	"context"

	"github.com/cloudwego/eino/schema"
	"github.com/cloudwego/hertz/pkg/app"
)

type ChatRequest struct {
	Message string `json:"message" binding:"required"`
}

type ChatResponse struct {
	Code    int             `json:"code"`
	Message string          `json:"message,omitempty"`
	Data    *schema.Message `json:"data,omitempty"`
}

func SendChat(c context.Context, ctx *app.RequestContext) {
	var req ChatRequest
	if err := ctx.BindAndValidate(&req); err != nil {
		ctx.JSON(400, ChatResponse{
			Code:    400,
			Message: "无效的请求参数",
		})
		return
	}
	messages := []*schema.Message{
		schema.SystemMessage("你是一个助手"),
		schema.UserMessage(req.Message),
	}
	// 调用本地模型进行对话
	resp, err := global.ChatModel.Generate(c, messages)
	if err != nil {
		ctx.JSON(500, ChatResponse{
			Code:    500,
			Message: "模型调用失败: " + err.Error(),
		})
		return
	}
	ctx.JSON(200, ChatResponse{
		Code: 200,
		Data: resp,
	})
}
