package tool_demo

import (
	"context"
	"log"

	"github.com/cloudwego/eino/callbacks"
	"github.com/cloudwego/eino/compose"
	"github.com/cloudwego/eino/schema"
)

func Main() {
	ctx := context.Background()
	handler := callbacks.NewHandlerBuilder().
		OnStartFn(
			func(ctx context.Context, info *callbacks.RunInfo, input callbacks.CallbackInput) context.Context {
				log.Printf("onStart, runInfo: %v, input: %v", info, input)
				return ctx
			}).
		OnEndFn(
			func(ctx context.Context, info *callbacks.RunInfo, output callbacks.CallbackOutput) context.Context {
				log.Printf("onEnd, runInfo: %v, out: %v", info, output)
				return ctx
			}).
		Build()

	graph := compose.NewGraph[map[string]any, *schema.Message]()
	compiledGraph, err := graph.Compile(ctx)
	if err != nil {
		panic(err)
	}
	// 注入到 graph 运行中
	compiledGraph.Invoke(ctx, map[string]any{"query": "Beijing's weather this weekend"}, compose.WithCallbacks(handler))

}