package react_agent

import (
	"context"
	"errors"
	"fmt"
	"io"
	"learn-eino/util"
	"learn-eino/util/tool/task"

	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/compose"
	"github.com/cloudwego/eino/flow/agent/react"
	"github.com/cloudwego/eino/schema"
)

func Main() {
	Main3()
}

/*

react 作为一个节点，在其内部，chatModel 负责调用 tools。

*/

func Main1() {
	ctx := context.Background()

	chatModel, err := util.GetChatModel(ctx)
	if err != nil {
		fmt.Printf("Get ChatModel failed, err=%v\n", err)
		return
	}

	taskTool, err := task.NewTaskTool(ctx, nil)
	if err != nil {
		fmt.Printf("NewTaskTool failed, err=%v\n", err)
		return
	}

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
		StreamToolCallChecker: util.CustomToolCallChecker,
	})

	response, err := reactAgent.Stream(ctx, []*schema.Message{schema.UserMessage(`帮我在task manager系统创建一个任务：
标题：eino_agent任务标题；
内容：eino_agent任务内容；
截止时间：2025-12-25T15:15；
任务ID：eac76c5a-629a-45d2-8c4e-0ac769a470f0；
任务状态：未完成；
创建时间：2026-01-07T10:03:07+08:00
`)})
	if err != nil {
		fmt.Printf("Generate failed, err=%v\n", err)
		return
	}
	defer response.Close()
	for {
		msg, err := response.Recv()
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			fmt.Printf("failed to recv: %v", err)
			return
		}
		fmt.Printf("%v", msg.Content)
	}
}
