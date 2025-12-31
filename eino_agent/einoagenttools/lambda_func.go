package einoagenttools

import (
	"context"
	"time"
)

// inputToHistoryLambda component initialization function of node 'InputToHistory' in graph 'EinoAgentTools'
func inputToHistoryLambda(ctx context.Context, input *UserMessage, opts ...any) (output map[string]any, err error) {
	return map[string]any{
		"content": input.Query,
		"history": input.History,
		"date":    time.Now().Format("2006-01-02 15:04:05"),
	}, nil
}
