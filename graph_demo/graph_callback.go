package graph_demo

import (
	"context"
	"fmt"
	"log"

	"github.com/cloudwego/eino/callbacks"
	"github.com/cloudwego/eino/compose"
	"github.com/cloudwego/eino/schema"
)

func Main4() {
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

	g := compose.NewGraph[map[string]any, *schema.Message]()

	g.AddLambdaNode("lambda", compose.InvokableLambda(func(ctx context.Context, in map[string]any) (*schema.Message, error) {
		return &schema.Message{
			Role:    "assistant",
			Content: "Beijing's weather this weekend is cloud",
		}, nil
	}))
	_ = g.AddEdge(compose.START, "lambda")
	_ = g.AddEdge("lambda", compose.END)

	runable, err := g.Compile(ctx)
	if err != nil {
		panic(err)
	}

	// callback 会被加到所有的节点上
	res, err := runable.Invoke(ctx, map[string]any{"query": "Beijing's weather this weekend"}, compose.WithCallbacks(handler))
	if err != nil {
		panic(err)
	}

	fmt.Println(res.Content)
}
