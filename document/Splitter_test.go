package document

import (
	"SuperAgent/initialize"
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/cloudwego/eino/schema"
	"github.com/stretchr/testify/assert"
)

func TestSplit(t *testing.T) {
	initialize.InitConfig("..")
	initialize.InitClient()
	initialize.InitEmbedder()
	initialize.InitIndexer()
	// 2. 读取测试文件
	content, err := os.ReadFile("./test.txt")
	assert.NoError(t, err, "读取test.txt失败")
	fmt.Printf("文件内容长度: %d 字节\n", len(content))

	// 3. 创建原始文档
	ctx := context.Background()
	docs := []*schema.Document{
		{
			ID:      fmt.Sprintf("test-%d", time.Now().Unix()),
			Content: string(content),
			MetaData: map[string]any{
				"source": "test.txt",
				"type":   "text",
				"date":   time.Now().Format("2006-01-02"),
			},
		},
	}

	// 4. 分割文档
	fmt.Println("开始分割文档...")
	chunks, err := Split(ctx, docs)
	assert.NoError(t, err, "分割文档失败")
	assert.NotEmpty(t, chunks, "分割后的文档片段不应为空")
	for _, chunk := range chunks {
		fmt.Println("分割的文档id为", chunk.ID, "内容为", chunk.Content)
	}
}

func TestSplitAndStore(t *testing.T) {
	// 1. 初始化环境
	initialize.InitConfig("..")
	initialize.InitClient()
	initialize.InitEmbedder()
	initialize.InitIndexer()

	// 2. 读取测试文件
	content, err := os.ReadFile("./test.txt")
	assert.NoError(t, err, "读取test.txt失败")
	fmt.Printf("文件内容长度: %d 字节\n", len(content))

	// 3. 创建原始文档
	ctx := context.Background()
	docs := []*schema.Document{
		{
			ID:      fmt.Sprintf("test-%d", time.Now().Unix()),
			Content: string(content),
			MetaData: map[string]any{
				"source": "test.txt",
				"type":   "text",
				"date":   time.Now().Format("2006-01-02"),
			},
		},
	}

	// 4. 分割文档
	fmt.Println("开始分割文档...")
	chunks, err := Split(ctx, docs)
	assert.NoError(t, err, "分割文档失败")
	assert.NotEmpty(t, chunks, "分割后的文档片段不应为空")
	fmt.Printf("文档分割完成，共生成 %d 个片段\n", len(chunks))

	// 5. 存储到向量数据库
	fmt.Println("开始存储文档片段到向量数据库...")
	for _, chunk := range chunks {
		chunk.ID = fmt.Sprintf("test-%d", time.Now().Unix())
		err = BuildRAG(chunk.ID, chunk.Content, chunk.MetaData, ctx)
		assert.NoError(t, err, "存储文档片段失败")
	}
	fmt.Println("测试完成!")
}
