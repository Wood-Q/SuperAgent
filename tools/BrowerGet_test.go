package tools

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBrowerGetDocument(t *testing.T) {
	// 设置测试网址，这里使用明日方舟Wiki的某个干员页面作为示例
	// 您可以更换为其他任何页面进行测试
	url := "https://prts.wiki/w/缪尔赛思"

	// 调用 BrowerGet 函数
	result := BrowerGetDocument(url)

	// 打印结果
	fmt.Println("结果输出:")
	fmt.Println(result)

	require.NotEmpty(t, result)
}
