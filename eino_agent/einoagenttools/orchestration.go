package einoagenttools

import (
	"context"

	"github.com/cloudwego/eino/compose"
	"github.com/cloudwego/eino/schema"
)

func BuildEinoAgentTools(ctx context.Context) (r compose.Runnable[*UserMessage, *schema.Message], err error) {
	const (
		ChatTemplate   = "ChatTemplate"
		ReactAgent     = "ReactAgent"
		InputToHistory = "InputToHistory"
	)
	g := compose.NewGraph[*UserMessage, *schema.Message]()
	chatTemplateKeyOfChatTemplate, err := newChatTemplate(ctx)
	if err != nil {
		return nil, err
	}
	_ = g.AddChatTemplateNode(ChatTemplate, chatTemplateKeyOfChatTemplate)
	reactAgentKeyOfLambda, err := reactAgentLambda(ctx)
	if err != nil {
		return nil, err
	}
	_ = g.AddLambdaNode(ReactAgent, reactAgentKeyOfLambda, compose.WithNodeName("ReAct Agent"))
	_ = g.AddLambdaNode(InputToHistory, compose.InvokableLambdaWithOption(inputToHistoryLambda), compose.WithNodeName("UserMessageToVariables"))
	_ = g.AddEdge(compose.START, InputToHistory)
	_ = g.AddEdge(ReactAgent, compose.END)
	_ = g.AddEdge(InputToHistory, ChatTemplate)
	_ = g.AddEdge(ChatTemplate, ReactAgent)
	r, err = g.Compile(ctx, compose.WithGraphName("EinoAgentTools"), compose.WithNodeTriggerMode(compose.AllPredecessor))
	if err != nil {
		return nil, err
	}
	return r, err
}
