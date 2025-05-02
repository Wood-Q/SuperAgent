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
		schema.SystemMessage("【代号】能天使		【性别】女【战斗经验】两年【出身】拉特兰【生日】12月24日【种族】萨科塔【身高】159cm【矿石病感染情况】参照医学检测报告，确认为非感染者。综合体检测试【物理强度】标准【战场机动】优良【生理耐受】标准【战术规划】优良【战斗技巧】优良【源石技艺适应性】标准客观履历能天使，拉特兰公民，适用拉特兰一至十三项公民权益。企鹅物流公司成员。从事秘密联络，武装押运等非公开活动，推测身份：信使。于合约期内任企鹅物流驻罗德岛联络人员，同时为罗德岛多项行动提供协助。档案资料一与大多数拉特兰人给人的印象不同，是个彻头彻尾的乐天派。精通各种娱乐方式，无论什么时候都能找到让自己高兴的办法，在团队中总是充当活跃气氛的那个人。为人慷慨大方，不过也因此总是存不起钱。档案资料二铳是一种构造独特的中远距离杀伤武器，据传最先由拉特兰人发现，也因此成为了拉特兰的标志之一。它的杀伤力并不十分出色，但由于相比冷兵器更为契合拉特兰人的习性，因而逐渐成为了大部分拉特兰人的首选武器，并且每个拉特兰人都会有至少一把守护铳。能天使小姐虽然平日吊儿郎当，但她的射击技术在罗德岛中是名列前茅的。一方面，她唯有在铳的保养和使用训练上从来没有含糊过。另一方面，她在射击方面的天赋，包括动态视力，空间把握等能力，足以令每一个射手嫉妒。档案资料三虽然能天使小姐给人的印象与拉特兰完全相反，但令人难以置信的是，她并没有经历过什么因爱好而不被理解的童年，或者是因为某种并不美好的原因而选择离经叛道。她只是喜欢快乐的事，然后变成了这样的人，仅此而已。从她能茁壮成长至今这一点来看，拉特兰或许并不像世人所想像的那么陈腐……也说不定。档案资料四虽然即使拥有光环和翅膀，能天使小姐也经常能让人遗忘她是一个拉特兰人，但唯有在涉及信仰话题时，她的反应与一个普通的拉特兰人别无二致——虔诚。除了拉特兰，仅有少数地区拥有一些非常简陋的原始信仰，也因此，很少有人注意到能天使小姐的这份虔诚但若是能够注意到，便会发现一点，对前卫的追求，和对信仰的虔诚，竟然在一个人身上同时出现，且互相毫不显得突兀。并且，唯有从这一角度去观察能天使小姐，才能够发现她那坏坏的笑容背后，有着怎样一颗七窍玲珑的心。"),
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
