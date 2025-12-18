package einoagent

import (
	"context"

	"github.com/cloudwego/eino/compose"
	"github.com/cloudwego/eino/schema"
)

func BuildEinoAgentEs8(ctx context.Context) (r compose.Runnable[*UserMessage, *schema.Message], err error) {
	const (
		InputToQuery            = "InputToQuery"
		ChatTemplate            = "ChatTemplate"
		ReactAgent              = "ReactAgent"
		InputToHistory          = "InputToHistory"
		Retriever1              = "Retriever1"
		DocumentListToMapLambda = "DocumentListToMapLambda"
	)
	g := compose.NewGraph[*UserMessage, *schema.Message]()
	_ = g.AddLambdaNode(InputToQuery, compose.InvokableLambdaWithOption(inputToQueryLambda), compose.WithNodeName("UserMessageToQuery"))
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
	_ = g.AddLambdaNode(DocumentListToMapLambda, compose.InvokableLambdaWithOption(documentListToMapLambda), compose.WithNodeName("DocumentListToMapLambda"))
	retriever1KeyOfRetriever, err := newRetriever(ctx)
	if err != nil {
		return nil, err
	}
	_ = g.AddRetrieverNode(Retriever1, retriever1KeyOfRetriever)
	_ = g.AddEdge(compose.START, InputToQuery)
	_ = g.AddEdge(compose.START, InputToHistory)
	_ = g.AddEdge(ReactAgent, compose.END)
	_ = g.AddEdge(InputToQuery, Retriever1)
	_ = g.AddEdge(InputToHistory, ChatTemplate)
	_ = g.AddEdge(Retriever1, DocumentListToMapLambda)
	_ = g.AddEdge(DocumentListToMapLambda, ChatTemplate)
	_ = g.AddEdge(ChatTemplate, ReactAgent)
	r, err = g.Compile(ctx, compose.WithGraphName("EinoAgentEs8"), compose.WithNodeTriggerMode(compose.AllPredecessor))
	if err != nil {
		return nil, err
	}
	return r, err
}
