package api

import (
	"SuperAgent/global"
	"SuperAgent/message_model"
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
	var req message_model.ChatRequest
	if err := ctx.BindAndValidate(&req); err != nil {
		ctx.JSON(400, message_model.ChatResponse{
			Code:    400,
			Message: "无效的请求参数",
		})
		return
	}
	messages := []*schema.Message{
		schema.SystemMessage("你是明日方舟游戏的角色阿米娅。罗德岛的公开领袖，在内部拥有最高执行权。虽然，从外表上看起来仅仅是个不成熟的少女，实际上，她却是深受大家信任的合格的领袖。兼具无私奉献的精神、乐观向上的态度和温暖真挚的同理心，同时又不失坚韧负责的领导力与青春活力。这样的助手既能在技术与信息上给予用户专业支持，又能在情感上给予及时的安慰与鼓励，让互动既高效又富有人情味。喜欢称呼我为“博士”。"),
		schema.UserMessage(req.Message),
	}
	// 调用本地模型进行对话
	resp, err := global.ChatModel.Generate(c, messages)
	if err != nil {
		ctx.JSON(500, message_model.ChatResponse{
			Code:    500,
			Message: "模型调用失败: " + err.Error(),
		})
		return
	}
	ctx.JSON(200, message_model.ChatResponse{
		Code: 200,
		Data: resp,
	})
}

func DoChatWithSchemaJSON(c context.Context, ctx *app.RequestContext) {
	var req message_model.ChatRequest
	if err := ctx.BindAndValidate(&req); err != nil {
		ctx.JSON(400, message_model.ChatResponse{
			Code:    400,
			Message: "无效的请求参数",
		})
		return
	}
	messages := []*schema.Message{
		schema.SystemMessage("你是一个专业的健康助理，请根据用户提供的身高体重信息，生成一份详细的身体健康报告，输出必须是严格的JSON格式"),
		schema.UserMessage(req.Message),
	}
	resp, err := global.ChatModel.Generate(c, messages)
	if err != nil {
		ctx.JSON(500, message_model.ChatResponse{
			Code:    500,
			Message: "模型调用失败: " + err.Error(),
		})
		return
	}
	ctx.JSON(200, message_model.ChatResponse{
		Code: 200,
		Data: resp,
	})
}
