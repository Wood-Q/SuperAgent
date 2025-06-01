package main

import (
	"MoonAgent/cmd/di"
	"MoonAgent/pkg/logger"
	"MoonAgent/pkg/splitter"
	"context"
	"fmt"
	"os"

	"github.com/cloudwego/eino/schema"
	"go.uber.org/zap"
)

func main() {
	// 初始化测试环境的logger
	logger.InitTestLogger()

	ctx := context.Background()
	zap.S().Info("开始初始化应用程序...")

	app, clear, err := di.InitializeApplication()
	if err != nil {
		zap.S().Fatalf("应用程序初始化失败: %v", err)
	}
	defer clear()

	zap.S().Info("应用程序初始化成功")
	zap.S().Info("开始读取文档文件...")

	content, err := os.ReadFile("E:\\Study\\FullStack\\Project\\SuperAgent\\assets\\documents\\muelsyse copy.txt")
	if err != nil {
		zap.S().Fatalf("文档文件读取失败: %v", err)
	}

	zap.S().Infof("文档文件读取成功，文件大小: %d 字节", len(content))

	document := []*schema.Document{
		{
			ID:      "muelsyse",
			Content: string(content),
		},
	}

	zap.S().Info("开始文档分割...")
	docs, err := splitter.SplitDocs(ctx, app.Embedder, document)
	if err != nil {
		zap.S().Fatalf("文档分割失败: %v", err)
	}

	zap.S().Infof("文档分割成功，共生成 %d 个文档块", len(docs))

	zap.S().Info("开始存储文档到ES8...")
	result, err := app.Indexer.Store(ctx, docs)
	if err != nil {
		zap.S().Fatalf("文档存储失败: %v", err)
	}

	zap.S().Infof("文档存储成功，存储ID: %v", result)
	fmt.Println("Document stored successfully", result)
}
