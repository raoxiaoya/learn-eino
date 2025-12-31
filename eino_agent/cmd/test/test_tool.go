package main

import (
	"context"
	"encoding/json"
	"learn-eino/util/tool/task"
	"os"

	"github.com/cloudwego/eino-ext/components/model/ark"
	"github.com/cloudwego/eino/adk"
	"github.com/cloudwego/eino/components/model"
	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/compose"
)

var systemPrompt = `
# Role: Eino Expert Assistant

## Core Competencies
- knowledge of Eino framework and ecosystem
- Project scaffolding and best practices consultation
- Documentation navigation and implementation guidance
- Search web, clone github repo, open file/url, task management

## Interaction Guidelines
- Before responding, ensure you:
  • Fully understand the user's request and requirements, if there are any ambiguities, clarify with the user
  • Consider the most appropriate solution approach

- When providing assistance:
  • Be clear and concise
  • Include practical examples when relevant
  • Reference documentation when helpful
  • Suggest improvements or next steps if applicable

- If a request exceeds your capabilities:
  • Clearly communicate your limitations, suggest alternative approaches if possible

- If the question is compound or complex, you need to think step by step, avoiding giving low-quality answers directly.

`

func main() {
	fun1()
}

func fun1() {
	ctx := context.Background()
	taskTool, err := task.NewTaskTool(ctx, nil)
	if err != nil {
		panic(err)
	}

	str, _ := json.Marshal(task.TaskRequest{
		Action: task.ActionAdd,
		Task: &task.Task{
			Title:    "test task",
			Content:  "test content",
			Deadline: "2025-12-30 12:12:12",
		},
	})

	tt, ok := taskTool.(tool.InvokableTool)
	if !ok {
		panic("taskTool is not InvokableTool")
	}
	res, err := tt.InvokableRun(ctx, string(str))
	if err != nil {
		panic(err)
	}
	println(res)
}

func NewResumAnalysisAgent() adk.Agent {
	ctx := context.Background()

	taskTool, err := task.NewTaskTool(ctx, nil)
	if err != nil {
		panic(err)
	}

	chatModelIns, err := newChatModel(ctx)
	if err != nil {
		panic(err)
	}

	a, err := adk.NewChatModelAgent(ctx, &adk.ChatModelAgentConfig{
		Name:        "TaskToolAgent",
		Description: "调用task manager工具",
		Instruction: systemPrompt,
		Model:       chatModelIns,
		ToolsConfig: adk.ToolsConfig{
			ToolsNodeConfig: compose.ToolsNodeConfig{
				Tools: []tool.BaseTool{taskTool},
			},
		},
	})
	return a
}

func newChatModel(ctx context.Context) (cm model.ToolCallingChatModel, err error) {
	// TODO Modify component configuration here.
	config := &ark.ChatModelConfig{
		Model:  os.Getenv("ARK_CHAT_MODEL"),
		APIKey: os.Getenv("ARK_API_KEY"),
	}
	cm, err = ark.NewChatModel(ctx, config)
	if err != nil {
		return nil, err
	}
	return cm, nil
}
