// file: cmd/di/wire.go
//go:build wireinject
// +build wireinject

package di

import (
	"github.com/google/wire"
)

// InitializeApplication 是我们的 injector 函数
// 它声明了我们想构建 *Application，并列出了所有的 Provider
func InitializeApplication() (*Application, func(), error) { // 增加了清理函数
	// wire.Build 会分析 ProviderSet 中的函数，并生成代码来构建 Application
	// 清理函数用于释放资源，比如关闭 MilvusClient
	panic(wire.Build(
		ProviderSet,
	))
}
