package api

import (
	"SuperAgent/global"
	"SuperAgent/model"
	"context"

	"github.com/cloudwego/eino/schema"
	"github.com/cloudwego/hertz/pkg/app"
)

// SendChat godoc
// @Summary 发送聊天消息
// @Description 处理用户发送的聊天消息并返回模型生成的响应
// @Tags Chat
// @Accept json
// @Produce json
// @Param chatRequest body model.ChatRequest true "聊天请求"
// @Success 200 {object} model.ChatResponse "成功响应"
// @Failure 400 {object} model.ChatResponse "无效的请求参数"
// @Failure 500 {object} model.ChatResponse "模型调用失败"
// @Router /chat/send [post]
func SendChat(c context.Context, ctx *app.RequestContext) {
	var req model.ChatRequest
	if err := ctx.BindAndValidate(&req); err != nil {
		ctx.JSON(400, model.ChatResponse{
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
		ctx.JSON(500, model.ChatResponse{
			Code:    500,
			Message: "模型调用失败: " + err.Error(),
		})
		return
	}
	ctx.JSON(200, model.ChatResponse{
		Code: 200,
		Data: resp,
	})
}
