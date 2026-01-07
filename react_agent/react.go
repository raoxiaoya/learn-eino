/*
 * Copyright 2024 CloudWeGo Authors
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

package react_agent

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"time"

	"github.com/cloudwego/eino/callbacks"
	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/compose"
	"github.com/cloudwego/eino/flow/agent"
	"github.com/cloudwego/eino/flow/agent/react"
	"github.com/cloudwego/eino/schema"

	"learn-eino/react_agent/tools"
	"learn-eino/util"

	"learn-eino/react_agent/visualize"
)

func Main2() {
	// arkApiKey := os.Getenv("ARK_API_KEY")
	// arkModelName := os.Getenv("ARK_MODEL_NAME")
	// cozeloopApiToken := os.Getenv("COZELOOP_API_TOKEN")
	// cozeloopWorkspaceID := os.Getenv("COZELOOP_WORKSPACE_ID") // use cozeloop trace, from https://loop.coze.cn/open/docs/cozeloop/go-sdk#4a8c980e

	ctx := context.Background()
	// var handlers []callbacks.Handler
	// if cozeloopApiToken != "" && cozeloopWorkspaceID != "" {
	// 	client, err := cozeloop.NewClient(
	// 		cozeloop.WithAPIToken(cozeloopApiToken),
	// 		cozeloop.WithWorkspaceID(cozeloopWorkspaceID),
	// 	)
	// 	if err != nil {
	// 		panic(err)
	// 	}
	// 	defer client.Close(ctx)
	// 	handlers = append(handlers, clc.NewLoopHandler(client))
	// }
	// callbacks.AppendGlobalHandlers(handlers...)

	// minimal: we will export graph via API when available and compile a mermaid diagram

	// Create a new cached ark chat model.
	//arkModel, err = NewCachedARKChatModel(ctx, config)

	
	arkModel, err := util.GetChatModel(ctx)
	if err != nil {
		fmt.Printf("failed to create chat model: %v\n", err)
		return
	}

	// prepare tools
	restaurantTool := tools.GetRestaurantTool()
	dishTool := tools.GetDishTool()

	// prepare persona (system prompt) (optional)
	persona := `# Character:
你是一个帮助用户推荐餐厅和菜品的助手，根据用户的需要，查询餐厅信息并推荐，查询餐厅的菜品并推荐。
`

	// replace tool call checker with a custom one: check all trunks until you get a tool call
	// because some models(claude or doubao 1.5-pro 32k) do not return tool call in the first response
	// uncomment the following code to enable it
	/*toolCallChecker := func(ctx context.Context, sr *schema.StreamReader[*schema.Message]) (bool, error) {
		defer sr.Close()
		for {
			msg, err := sr.Recv()
			if err != nil {
				if errors.Is(err, io.EOF) {
					// finish
					break
				}

				return false, err
			}

			if len(msg.ToolCalls) > 0 {
				return true, nil
			}
		}
		return false, nil
	}*/

	rAgent, err := react.NewAgent(ctx, &react.AgentConfig{
		ToolCallingModel: arkModel,
		ToolsConfig: compose.ToolsNodeConfig{
			Tools: []tool.BaseTool{restaurantTool, dishTool},
		},
		// StreamToolCallChecker: toolCallChecker, // uncomment it to replace the default tool call checker with custom one
	})
	if err != nil {
		fmt.Printf("failed to create agent: %v\n", err)
		return
	}

	// if you want ping/pong, use Generate
	// msg, err := agent.Generate(ctx, []*schema.Message{
	// 	{
	// 		Role:    schema.User,
	// 		Content: "我在北京，给我推荐一些菜，需要有口味辣一点的菜，至少推荐有 2 家餐厅",
	// 	},
	// }, react.WithCallbacks(&myCallback{}))
	// if err != nil {
	// 	log.Printf("failed to generate: %v\n", err)
	// 	return
	// }
	// fmt.Println(msg.String())

	// If you want to use ark caching in react, call ark.WithCache()
	//cacheOption := &ark.CacheOption{
	//	APIType: ark.ResponsesAPI,
	//	SessionCache: &ark.SessionCacheConfig{
	//		EnableCache: true,
	//		TTL:         3600,
	//	},
	//}

	opt := []agent.AgentOption{
		agent.WithComposeOptions(compose.WithCallbacks(&LoggerCallback{})),
		//react.WithChatModelOptions(ark.WithCache(cacheOption)),
	}

	// Export graph and compile with mermaid (non-critical path)
	{
		anyG, opts := rAgent.ExportGraph()
		gen := visualize.NewMermaidGenerator("flow/agent/react")
		g := compose.NewGraph[[]*schema.Message, *schema.Message]()
		_ = g.AddGraphNode("react_agent", anyG, opts...)
		_ = g.AddEdge(compose.START, "react_agent")
		_ = g.AddEdge("react_agent", compose.END)
		_, _ = g.Compile(context.Background(), compose.WithGraphCompileCallbacks(gen))
	}

	sr, err := rAgent.Stream(ctx, []*schema.Message{
		{
			Role:    schema.System,
			Content: persona,
		},
		{
			Role:    schema.User,
			Content: "我在北京，给我推荐一些菜，需要有口味辣一点的菜，至少推荐有 2 家餐厅",
		},
	}, opt...)
	if err != nil {
		fmt.Printf("failed to stream: %v\n", err)
		return
	}

	defer sr.Close() // remember to close the stream

	fmt.Println("\n\n===== start streaming =====\n\n")

	for {
		msg, err := sr.Recv()
		if err != nil {
			if errors.Is(err, io.EOF) {
				// finish
				break
			}
			// error
			fmt.Printf("failed to recv: %v\n", err)
			return
		}

		// 打字机打印
		fmt.Printf("%v", msg.Content)
	}

	fmt.Println("\n\n===== finished =====\n")
	time.Sleep(2 * time.Second)
}

type LoggerCallback struct {
	callbacks.HandlerBuilder // 可以用 callbacks.HandlerBuilder 来辅助实现 callback
}

func (cb *LoggerCallback) OnStart(ctx context.Context, info *callbacks.RunInfo, input callbacks.CallbackInput) context.Context {
	fmt.Println("==================")
	inputStr, _ := json.MarshalIndent(input, "", "  ")
	fmt.Printf("[OnStart] %s\n", string(inputStr))
	return ctx
}

func (cb *LoggerCallback) OnEnd(ctx context.Context, info *callbacks.RunInfo, output callbacks.CallbackOutput) context.Context {
	fmt.Println("=========[OnEnd]=========")
	outputStr, _ := json.MarshalIndent(output, "", "  ")
	fmt.Println(string(outputStr))
	return ctx
}

func (cb *LoggerCallback) OnError(ctx context.Context, info *callbacks.RunInfo, err error) context.Context {
	fmt.Println("=========[OnError]=========")
	fmt.Println(err)
	return ctx
}

func (cb *LoggerCallback) OnEndWithStreamOutput(ctx context.Context, info *callbacks.RunInfo,
	output *schema.StreamReader[callbacks.CallbackOutput]) context.Context {

	var graphInfoName = react.GraphName

	go func() {
		defer func() {
			if err := recover(); err != nil {
				fmt.Println("[OnEndStream] panic err:", err)
			}
		}()

		defer output.Close() // remember to close the stream in defer

		fmt.Println("=========[OnEndStream]=========")
		for {
			frame, err := output.Recv()
			if errors.Is(err, io.EOF) {
				// finish
				break
			}
			if err != nil {
				fmt.Printf("internal error: %s\n", err)
				return
			}

			s, err := json.Marshal(frame)
			if err != nil {
				fmt.Printf("internal error: %s\n", err)
				return
			}

			if info.Name == graphInfoName { // 仅打印 graph 的输出, 否则每个 stream 节点的输出都会打印一遍
				fmt.Printf("%s: %s\n", info.Name, string(s))
			}
		}

	}()
	return ctx
}

func (cb *LoggerCallback) OnStartWithStreamInput(ctx context.Context, info *callbacks.RunInfo,
	input *schema.StreamReader[callbacks.CallbackInput]) context.Context {
	defer input.Close()
	return ctx
}