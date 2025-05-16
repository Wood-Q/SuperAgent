package main

import (
	"context"
	"errors"
	"io"
	"runtime/debug"
	"strings"
	"unicode/utf8"

	"github.com/cloudwego/eino/compose"
	"github.com/cloudwego/eino/schema"
	"github.com/cloudwego/eino/utils/safe"

	"github.com/cloudwego/eino-examples/internal/logs"
)

// state 函数演示了如何在Eino图计算引擎中使用状态管理
func state() {
	// 创建一个空的上下文，用于贯穿整个计算过程
	ctx := context.Background()

	// 定义图中节点的名称常量
	const (
		nodeOfL1 = "invokable"     // 同步调用节点
		nodeOfL2 = "streamable"    // 流式输出节点
		nodeOfL3 = "transformable" // 流式转换节点
	)

	// 定义测试状态结构体，用于在图计算过程中传递和维护状态
	type testState struct {
		ms []string // 用于存储消息历史
	}

	// 状态生成器函数，每次图执行时创建一个新的状态实例
	gen := func(ctx context.Context) *testState {
		return &testState{} // 返回一个空的状态实例
	}

	// 创建一个新的图，使用泛型参数 [string, string] 表示输入和输出类型都是字符串
	// WithGenLocalState 配置图使用上面定义的状态生成器
	sg := compose.NewGraph[string, string](compose.WithGenLocalState(gen))

	// ============== 第一个节点：同步处理节点 ==============

	// 创建一个InvokableLambda函数，接收字符串输入并返回处理后的字符串
	l1 := compose.InvokableLambda(func(ctx context.Context, in string) (out string, err error) {
		return "InvokableLambda: " + in, nil // 简单地在输入前添加前缀
	})

	// 节点1的状态前置处理器：在执行节点逻辑前处理状态
	l1StateToInput := func(ctx context.Context, in string, state *testState) (string, error) {
		state.ms = append(state.ms, in) // 将输入消息添加到状态历史中
		return in, nil                  // 返回原始输入，不做修改
	}

	// 节点1的状态后置处理器：在执行节点逻辑后处理状态
	l1StateToOutput := func(ctx context.Context, out string, state *testState) (string, error) {
		state.ms = append(state.ms, out) // 将输出消息添加到状态历史中
		return out, nil                  // 返回原始输出，不做修改
	}

	// 将节点1添加到图中，并绑定状态处理器
	_ = sg.AddLambdaNode(nodeOfL1, l1,
		compose.WithStatePreHandler(l1StateToInput),   // 前置处理器
		compose.WithStatePostHandler(l1StateToOutput), // 后置处理器
	)

	// ============== 第二个节点：流式输出节点 ==============

	// 创建一个StreamableLambda函数，接收字符串输入并返回字符串流
	l2 := compose.StreamableLambda(func(ctx context.Context, input string) (output *schema.StreamReader[string], err error) {
		outStr := "StreamableLambda: " + input // 构建输出字符串

		// 创建一个管道，用于流式传输数据
		// 管道容量设置为输出字符串的长度
		sr, sw := schema.Pipe[string](utf8.RuneCountInString(outStr))

		// 启动一个goroutine处理流式输出
		// nolint: byted_goroutine_recover
		go func() {
			// 将输出字符串分割成单词，并逐个发送到流中
			for _, field := range strings.Fields(outStr) {
				sw.Send(field+" ", nil) // 发送单词和空格
			}
			sw.Close() // 关闭写入器，表示流结束
		}()

		return sr, nil // 返回流读取器
	})

	// 节点2的状态后置处理器
	l2StateToOutput := func(ctx context.Context, out string, state *testState) (string, error) {
		state.ms = append(state.ms, out) // 将流输出的内容添加到状态历史中
		return out, nil
	}

	// 将节点2添加到图中，只绑定后置处理器
	_ = sg.AddLambdaNode(nodeOfL2, l2, compose.WithStatePostHandler(l2StateToOutput))

	// ============== 第三个节点：流式转换节点 ==============

	// 创建一个TransformableLambda函数，接收字符串流输入并返回转换后的字符串流
	l3 := compose.TransformableLambda(func(ctx context.Context, input *schema.StreamReader[string]) (
		output *schema.StreamReader[string], err error) {

		prefix := "TransformableLambda: " // 定义前缀
		sr, sw := schema.Pipe[string](20) // 创建容量为20的管道

		// 启动goroutine处理流转换
		go func() {
			// 添加错误恢复机制
			defer func() {
				panicErr := recover()
				if panicErr != nil {
					// 如果发生panic，捕获并记录错误
					err := safe.NewPanicErr(panicErr, debug.Stack())
					logs.Errorf("panic occurs: %v\n", err)
				}
			}()

			// 先发送前缀的每个单词
			for _, field := range strings.Fields(prefix) {
				sw.Send(field+" ", nil)
			}

			// 读取输入流并转发到输出流
			for {
				chunk, err := input.Recv() // 从输入流接收数据
				if err != nil {
					if err == io.EOF {
						break // 如果到达流结束，退出循环
					}
					// 如果发生其他错误，将错误传递到输出流
					sw.Send(chunk, err)
					break
				}

				sw.Send(chunk, nil) // 将接收到的块发送到输出流
			}
			sw.Close() // 关闭写入器，表示流结束
		}()

		return sr, nil // 返回流读取器
	})

	// 节点3的状态后置处理器
	l3StateToOutput := func(ctx context.Context, out string, state *testState) (string, error) {
		state.ms = append(state.ms, out) // 将转换后的输出添加到状态历史中

		// 打印状态历史
		logs.Infof("state result: ")
		for idx, m := range state.ms {
			logs.Infof("    %vth: %v", idx, m)
		}
		return out, nil
	}

	// 将节点3添加到图中，只绑定后置处理器
	_ = sg.AddLambdaNode(nodeOfL3, l3, compose.WithStatePostHandler(l3StateToOutput))

	// ============== 构建图的边，连接各个节点 ==============

	// 从起始节点连接到节点1
	_ = sg.AddEdge(compose.START, nodeOfL1)

	// 从节点1连接到节点2
	_ = sg.AddEdge(nodeOfL1, nodeOfL2)

	// 从节点2连接到节点3
	_ = sg.AddEdge(nodeOfL2, nodeOfL3)

	// 从节点3连接到结束节点
	_ = sg.AddEdge(nodeOfL3, compose.END)

	// ============== 编译并执行图 ==============

	// 编译图，生成可执行的运行器
	run, err := sg.Compile(ctx)
	if err != nil {
		logs.Errorf("sg.Compile failed, err=%v", err)
		return
	}

	// ============== 同步调用模式 ==============

	// 使用Invoke方法同步调用图
	out, err := run.Invoke(ctx, "how are you")
	if err != nil {
		logs.Errorf("run.Invoke failed, err=%v", err)
		return
	}
	logs.Infof("invoke result: %v", out) // 打印调用结果

	// ============== 流式调用模式 ==============

	// 使用Stream方法流式调用图
	stream, err := run.Stream(ctx, "how are you")
	if err != nil {
		logs.Errorf("run.Stream failed, err=%v", err)
		return
	}

	// 读取流式调用的结果
	for {
		chunk, err := stream.Recv() // 接收一块数据
		if err != nil {
			if errors.Is(err, io.EOF) {
				break // 如果到达流结束，退出循环
			}
			logs.Infof("stream.Recv() failed, err=%v", err)
			break
		}

		logs.Tokenf("%v", chunk) // 打印流式输出的每一块
	}
	stream.Close() // 关闭流

	// ============== 流式转换模式 ==============

	// 创建一个输入流
	sr, sw := schema.Pipe[string](1)
	sw.Send("how are you", nil) // 发送一条消息
	sw.Close()                  // 关闭写入器

	// 使用Transform方法流式转换调用图
	stream, err = run.Transform(ctx, sr)
	if err != nil {
		logs.Infof("run.Transform failed, err=%v", err)
		return
	}

	// 读取转换结果
	for {
		chunk, err := stream.Recv() // 接收一块转换后的数据
		if err != nil {
			if errors.Is(err, io.EOF) {
				break // 如果到达流结束，退出循环
			}
			logs.Infof("stream.Recv() failed, err=%v", err)
			break
		}

		logs.Infof("%v", chunk) // 打印转换后的每一块数据
	}
	stream.Close() // 关闭流
}
