package model

// ChatMessage 聊天消息模型
type ChatMessage struct {
	Message    string `json:"message"    form:"message"    binding:"required"` // 用户输入的消息
	Response   string `json:"response"   form:"response"`                      // AI的响应
	CreateTime int64  `json:"createTime" form:"createTime"`                    // 消息创建时间
}

// ChatResponse 聊天响应模型
type ChatResponse struct {
	Code    int         `json:"code"`    // 响应状态码
	Message string      `json:"message"` // 响应信息
	Data    interface{} `json:"data"`    // 响应数据
}
