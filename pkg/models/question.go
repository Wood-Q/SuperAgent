package model

// QuestionRequest 问题请求
type QuestionRequest struct {
	Question string `json:"question"`
}

// QuestionResponse 问题响应
type QuestionResponse struct {
	Answer string `json:"answer"`
}
