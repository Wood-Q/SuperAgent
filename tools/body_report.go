package tools

import (
	"github.com/cloudwego/eino/components/tool/utils"
	"github.com/cloudwego/eino/schema"
	"go.uber.org/zap"
)

// BodyReport 定义了最终要输出的身体报告字段
type BodyReport struct {
	Height  float64 `json:"height" jsonschema:"required,description=身高（cm）"`
	Weight  float64 `json:"weight" jsonschema:"required,description=体重（kg）"`
	BMI     float64 `json:"bmi" jsonschema:"description=身体质量指数"`
	FatRate float64 `json:"fat_rate" jsonschema:"description=体脂率"`
}

func BodyReportTool() []*schema.ToolInfo {
	// 生成 ParamsOneOf（JSON Schema）描述
	paramsOneOf, err := utils.GoStruct2ParamsOneOf[BodyReport]()
	if err != nil {
		zap.S().Errorf("生成 ParamsOneOf 失败：%v", err)
	}

	// 构造 ToolInfo
	toolInfo := &schema.ToolInfo{
		Name:        "generate_body_report",
		Desc:        "生成用户身体报告(JSON)",
		ParamsOneOf: paramsOneOf,
	}

	return []*schema.ToolInfo{toolInfo}
}
