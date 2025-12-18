package einoagent

import (
	"context"
	"strings"
	"time"

	"github.com/cloudwego/eino/schema"
)

// inputToQueryLambda component initialization function of node 'InputToQuery' in graph 'EinoAgentEs8'
func inputToQueryLambda(ctx context.Context, input *UserMessage, opts ...any) (output string, err error) {
	return input.Query, nil
}

// inputToHistoryLambda component initialization function of node 'InputToHistory' in graph 'EinoAgentEs8'
func inputToHistoryLambda(ctx context.Context, input *UserMessage, opts ...any) (output map[string]any, err error) {
	return map[string]any{
		"content": input.Query,
		"history": input.History,
		"date":    time.Now().Format("2006-01-02 15:04:05"),
	}, nil
}

// 将 []*schema.Document 转换成 map[string]any
func documentListToMapLambda(ctx context.Context, input []*schema.Document, opts ...any) (output map[string]any, err error) {
	var contents []string
	for _, doc := range input {
		contents = append(contents, doc.Content)
	}
	contextText := strings.Join(contents, "\n\n")

	return map[string]interface{}{
		"retriever_result": contextText,
	}, nil
}
