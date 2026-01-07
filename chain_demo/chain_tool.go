package chain_demo

import (
	"context"
	"fmt"

	"learn-eino/util"
	"learn-eino/util/tool/task"

	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/compose"
	"github.com/cloudwego/eino/flow/agent/react"
	"github.com/cloudwego/eino/schema"
)

func Main8() {
	ctx := context.Background()

	chatModel, err := util.GetChatModel(ctx)
	if err != nil {
		fmt.Printf("NewChatModel failed, err=%v\n", err)
		return
	}

	taskTool, err := task.NewTaskTool(ctx, nil)
	info, err := taskTool.Info(ctx)
	if err != nil {
		fmt.Printf("Get ToolInfo failed, err=%v\n", err)
		return
	}

	err = chatModel.BindForcedTools([]*schema.ToolInfo{info})
	if err != nil {
		fmt.Printf("BindForcedTools failed, err=%v\n", err)
		return
	}

	// runableTaskTool, ok := taskTool.(tool.InvokableTool)
	// if !ok {
	// 	panic("taskTool is not InvokableTool")
	// }
	// lba, err := compose.AnyLambda(runableTaskTool.InvokableRun, nil, nil, nil)
	// if err != nil {
	// 	panic(err)
	// }

	reactAgent, err := react.NewAgent(ctx, &react.AgentConfig{
		ToolCallingModel: chatModel,
		ToolsConfig: compose.ToolsNodeConfig{
			Tools: []tool.BaseTool{taskTool},
		},
		// 直接将 tool 的结果返回，而不用再回到 chatModel 进行处理
		ToolReturnDirectly: map[string]struct{}{"task_manager": {}},
		MessageModifier: func(ctx context.Context, input []*schema.Message) []*schema.Message {
			res := make([]*schema.Message, 0, len(input)+1)
			res = append(res, schema.SystemMessage("你是一个后台操作员，使用 task_manager API，对task进行操作"))
			res = append(res, input...)
			return res
		},
	})

	lba, err := compose.AnyLambda(reactAgent.Generate, reactAgent.Stream, nil, nil)
	if err != nil {
		panic(err)
	}

	chain := compose.NewChain[[]*schema.Message, *schema.Message]().AppendLambda(lba)
	runable, err := chain.Compile(ctx)
	if err != nil {
		panic(err)
	}

	in := []*schema.Message{schema.UserMessage(`帮我在task manager系统创建一个任务：
标题：eino_agent任务标题；
内容：eino_agent任务内容；
截止时间：2026-12-25T15:15；
任务ID：eac76c5a-629a-45d2-8c4e-0ac769a470f0；
任务状态：未完成；
创建时间：2026-01-07T10:03:07+08:00
`)}

	response, err := runable.Invoke(ctx, in)
	if err != nil {
		panic(err)
	}
	fmt.Println(response.Role, response.Content)
}
