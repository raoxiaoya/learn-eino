package einoagenttools

import (
	"context"
	"learn-eino/util"
	"learn-eino/util/tool/task"

	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/compose"
	"github.com/cloudwego/eino/flow/agent/react"
)

// reactAgentLambda component initialization function of node 'ReactAgent' in graph 'EinoAgentTools'
func reactAgentLambda(ctx context.Context) (lba *compose.Lambda, err error) {
	chatModel, err := newChatModel(ctx)
	if err != nil {
		return nil, err
	}

	taskTool, err := task.NewTaskTool(ctx, nil)
	if err != nil {
		return nil, err
	}

	ins, err := react.NewAgent(ctx, &react.AgentConfig{
		ToolCallingModel: chatModel,
		ToolsConfig: compose.ToolsNodeConfig{
			Tools: []tool.BaseTool{taskTool},
		},
		StreamToolCallChecker: util.CustomToolCallChecker,
	})
	if err != nil {
		return nil, err
	}
	lba, err = compose.AnyLambda(ins.Generate, ins.Stream, nil, nil)
	if err != nil {
		return nil, err
	}
	return lba, nil
}
