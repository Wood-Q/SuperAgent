package constants

type AgentState string

const (
	//空闲状态
	AgentStateIdle AgentState = "idle"
	//运行中
	AgentStateRunning AgentState = "running"
	//成功完成状态
	AgentStateSuccess AgentState = "success"
	//失败状态
	AgentStateFailed AgentState = "failed"
	//错误状态
	AgentStateError AgentState = "error"
)
