package graph_demo

import (
	"context"
	"fmt"
	"learn-eino/util"

	"learn-eino/util/tool/task"

	"github.com/cloudwego/eino/callbacks"
	"github.com/cloudwego/eino/components/prompt"
	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/compose"
	"github.com/cloudwego/eino/schema"
)

func Main6() {
	ctx := context.Background()

	callbacks.AppendGlobalHandlers(&loggerCallbacks{})

	// 1. create an instance of ChatTemplate as 1st Graph Node
	chatTpl := prompt.FromMessages(schema.FString,
		schema.SystemMessage(`你是一个后台操作员，使用 task_manager API，对task进行操作`),
		schema.MessagesPlaceholder("message_histories", true),
		schema.UserMessage("{user_query}"),
	)

	chatModel, err := util.GetChatModel(ctx)
	if err != nil {
		fmt.Printf("NewChatModel failed, err=%v\n", err)
		return
	}

	// 3. create an instance of tool.InvokableTool for Intent recognition and execution
	taskTool, err := task.NewTaskTool(ctx, nil)
	info, err := taskTool.Info(ctx)
	if err != nil {
		fmt.Printf("Get ToolInfo failed, err=%v\n", err)
		return
	}

	// 4. bind ToolInfo to ChatModel. ToolInfo will remain in effect until the next BindTools.
	err = chatModel.BindForcedTools([]*schema.ToolInfo{info})
	if err != nil {
		fmt.Printf("BindForcedTools failed, err=%v\n", err)
		return
	}

	// 5. create an instance of ToolsNode as 3rd Graph Node
	toolsNode, err := compose.NewToolNode(ctx, &compose.ToolsNodeConfig{
		Tools: []tool.BaseTool{taskTool},
	})
	if err != nil {
		fmt.Printf("NewToolNode failed, err=%v\n", err)
		return
	}

	const (
		nodeKeyOfTemplate  = "template"
		nodeKeyOfChatModel = "chat_model"
		nodeKeyOfTools     = "tools"
	)

	// 6. create an instance of Graph
	// input type is 1st Graph Node's input type, that is ChatTemplate's input type: map[string]any
	// output type is last Graph Node's output type, that is ToolsNode's output type: []*schema.Message
	g := compose.NewGraph[map[string]any, []*schema.Message]()

	// 7. add ChatTemplate into graph
	_ = g.AddChatTemplateNode(nodeKeyOfTemplate, chatTpl)

	// 8. add ChatModel into graph
	_ = g.AddChatModelNode(nodeKeyOfChatModel, chatModel)

	// 9. add ToolsNode into graph
	_ = g.AddToolsNode(nodeKeyOfTools, toolsNode)

	// 10. add connection between nodes
	_ = g.AddEdge(compose.START, nodeKeyOfTemplate)

	_ = g.AddEdge(nodeKeyOfTemplate, nodeKeyOfChatModel)

	// chatModel 在 toolsNode 上游，这样才能调用 tools
	_ = g.AddEdge(nodeKeyOfChatModel, nodeKeyOfTools)

	_ = g.AddEdge(nodeKeyOfTools, compose.END)

	// 9. compile Graph[I, O] to Runnable[I, O]
	r, err := g.Compile(ctx)
	if err != nil {
		fmt.Printf("Compile failed, err=%v\n", err)
		return
	}

	out, err := r.Invoke(ctx, map[string]any{
		"message_histories": []*schema.Message{},
		"user_query":        `帮我在task manager系统创建一个任务：
标题：eino_agent任务标题；
内容：eino_agent任务标题；
截止时间：2025-12-25T15:15`,
	})
	if err != nil {
		fmt.Printf("Invoke failed, err=%v\n", err)
		return
	}
	fmt.Printf("Generation: %v Messages\n", len(out))
	for _, msg := range out {
		fmt.Printf("    %v\n", msg)
	}
}
