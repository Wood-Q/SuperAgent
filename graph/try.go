/*
 * Copyright 2025 CloudWeGo Authors
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/cloudwego/eino-ext/components/model/ollama"
	"github.com/cloudwego/eino/components/model"
	"github.com/cloudwego/eino/components/prompt"
	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/components/tool/utils"
	"github.com/cloudwego/eino/compose"
	"github.com/cloudwego/eino/schema"
)

// 添加参数管理结构
type ChatParams struct {
	Input   string            `json:"input"`
	Role    string            `json:"role"`
	Context map[string]string `json:"context"`
}


// 添加上下文管理
type ChatContext struct {
	History []*schema.Message
	Params  map[string]any
	State   *myState
}

func (cc *ChatContext) UpdateParams(newParams map[string]any) {
	for k, v := range newParams {
		cc.Params[k] = v
	}
}

func try() {
	compose.RegisterSerializableType[myState]("state")

	ctx := context.Background()
	runner, err := composeGraph[map[string]any, *schema.Message](
		ctx,
		newChatTemplate(ctx),
		newChatModel(ctx),
		newToolsNode(ctx),
		newCheckPointStore(ctx),
	)
	if err != nil {
		log.Fatal(err)
	}

	var history []*schema.Message
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("开始对话（输入'exit'退出）:")

	for {
		// 获取用户输入
		fmt.Print("\n您: ")
		input, err := reader.ReadString('\n')
		if err != nil {
			log.Printf("读取输入错误: %v", err)
			continue
		}

		// 处理输入
		input = strings.TrimSpace(input)
		if strings.ToLower(input) == "exit" {
			fmt.Println("再见！")
			break
		}

		// 构建参数映射
		params := map[string]any{
			"input": input,
			"role":  "AI助手", // 可以根据需要设置角色
		}

		// 调用模型
		result, err := runner.Invoke(ctx, params,
			compose.WithCheckPointID("chat-"+fmt.Sprint(len(history))),
			compose.WithStateModifier(func(ctx context.Context, path compose.NodePath, state any) error {
				state.(*myState).history = history
				return nil
			}),
		)

		if err == nil {
			// 输出模型响应
			fmt.Printf("\nAI: %s\n", result.Content)
			history = append(history, result)
			continue
		}

		// 处理工具调用
		info, ok := compose.ExtractInterruptInfo(err)
		if !ok {
			log.Printf("错误: %v", err)
			continue
		}

		history = info.State.(*myState).history
		for i, tc := range history[len(history)-1].ToolCalls {
			fmt.Printf("\n将要调用工具: %s\n参数: %s\n",
				tc.Function.Name, tc.Function.Arguments)
			fmt.Print("参数是否正确? (y/n): ")

			var response string
			fmt.Scanln(&response)

			if strings.ToLower(response) == "n" {
				fmt.Print("请输入新的参数: ")
				scanner := bufio.NewScanner(os.Stdin)
				if scanner.Scan() {
					history[len(history)-1].ToolCalls[i].Function.Arguments = scanner.Text()
				}
			}
		}
	}
}

func newChatTemplate(_ context.Context) prompt.ChatTemplate {
	return prompt.FromMessages(schema.FString,
		schema.SystemMessage(`你是一个{role}。
			你可以进行对话并在需要时调用工具。
			请保持友好和专业的态度。
			如果用户询问订票相关问题，使用"BookTicket"工具进行预订。`),
		schema.UserMessage("{input}"),
	)
}

func newChatModel(ctx context.Context) model.BaseChatModel {
	cm, err := ollama.NewChatModel(ctx, &ollama.ChatModelConfig{
		Model:   "qwen3:8b",
		BaseURL: "http://localhost:11434",
	})
	if err != nil {
		log.Fatal(err)
	}
	tools := getTools()
	var toolsInfo []*schema.ToolInfo
	for _, t := range tools {
		info, err := t.Info(ctx)
		if err != nil {
			log.Fatal(err)
		}
		toolsInfo = append(toolsInfo, info)
	}

	err = cm.BindTools(toolsInfo)
	if err != nil {
		log.Fatal(err)
	}
	return cm
}

type bookInput struct {
	Location             string `json:"location"`
	PassengerName        string `json:"passenger_name"`
	PassengerPhoneNumber string `json:"passenger_phone_number"`
}

func newToolsNode(ctx context.Context) *compose.ToolsNode {
	tools := getTools()

	tn, err := compose.NewToolNode(ctx, &compose.ToolsNodeConfig{Tools: tools})
	if err != nil {
		log.Fatal(err)
	}
	return tn
}

func newCheckPointStore(ctx context.Context) compose.CheckPointStore {
	return &myStore{buf: make(map[string][]byte)}
}

type myState struct {
	history []*schema.Message
}

func composeGraph[I, O any](ctx context.Context, tpl prompt.ChatTemplate, cm model.BaseChatModel, tn *compose.ToolsNode, store compose.CheckPointStore) (compose.Runnable[I, O], error) {
	g := compose.NewGraph[I, O](compose.WithGenLocalState(func(ctx context.Context) *myState {
		return &myState{}
	}))
	err := g.AddChatTemplateNode(
		"ChatTemplate",
		tpl,
	)
	if err != nil {
		return nil, err
	}
	err = g.AddChatModelNode(
		"ChatModel",
		cm,
		compose.WithStatePreHandler(func(ctx context.Context, in []*schema.Message, state *myState) ([]*schema.Message, error) {
			state.history = append(state.history, in...)
			return state.history, nil
		}),
		compose.WithStatePostHandler(func(ctx context.Context, out *schema.Message, state *myState) (*schema.Message, error) {
			state.history = append(state.history, out)
			return out, nil
		}),
	)
	if err != nil {
		return nil, err
	}
	err = g.AddToolsNode("ToolsNode", tn, compose.WithStatePreHandler(func(ctx context.Context, in *schema.Message, state *myState) (*schema.Message, error) {
		return state.history[len(state.history)-1], nil
	}))
	if err != nil {
		return nil, err
	}

	err = g.AddEdge(compose.START, "ChatTemplate")
	if err != nil {
		return nil, err
	}
	err = g.AddEdge("ChatTemplate", "ChatModel")
	if err != nil {
		return nil, err
	}
	err = g.AddEdge("ToolsNode", "ChatModel")
	if err != nil {
		return nil, err
	}
	err = g.AddBranch("ChatModel", compose.NewGraphBranch(func(ctx context.Context, in *schema.Message) (endNode string, err error) {
		if len(in.ToolCalls) > 0 {
			return "ToolsNode", nil
		}
		return compose.END, nil
	}, map[string]bool{"ToolsNode": true, compose.END: true}))
	if err != nil {
		return nil, err
	}
	return g.Compile(
		ctx,
		compose.WithCheckPointStore(store),
		compose.WithInterruptBeforeNodes([]string{"ToolsNode"}),
	)
}

func getTools() []tool.BaseTool {
	getWeather, err := utils.InferTool("BookTicket", "this tool can book ticket of the specific location", func(ctx context.Context, input bookInput) (output string, err error) {
		return "success", nil
	})
	if err != nil {
		log.Fatal(err)
	}

	return []tool.BaseTool{
		getWeather,
	}
}

type myStore struct {
	buf map[string][]byte
}

func (m *myStore) Get(ctx context.Context, checkPointID string) ([]byte, bool, error) {
	data, ok := m.buf[checkPointID]
	return data, ok, nil
}

func (m *myStore) Set(ctx context.Context, checkPointID string, checkPoint []byte) error {
	m.buf[checkPointID] = checkPoint
	return nil
}
