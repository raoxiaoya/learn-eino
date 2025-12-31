package graph_demo

import (
	"context"
	"errors"
	"fmt"
	"io"
	"runtime/debug"
	"strings"
	"unicode/utf8"

	"github.com/cloudwego/eino/compose"
	"github.com/cloudwego/eino/schema"
)

/*

Graph 可以有 graph 自身的“全局”状态，

*/

func Main3() {
	ctx := context.Background()

	const (
		nodeOfL1 = "invokable"
		nodeOfL2 = "streamable"
		nodeOfL3 = "transformable"
	)

	// 自定义 state 类型
	type testState struct {
		ms []string
	}

	gen := func(ctx context.Context) *testState {
		return &testState{}
	}

	// 在创建 Graph 时传入 WithGenLocalState Option 开启此功能
	sg := compose.NewGraph[string, string](compose.WithGenLocalState(gen))

	l1 := compose.InvokableLambda(func(ctx context.Context, in string) (out string, err error) {
		return "InvokableLambda: " + in, nil
	})

	l1StateToInput := func(ctx context.Context, in string, state *testState) (string, error) {
		state.ms = append(state.ms, in)
		return in, nil
	}

	l1StateToOutput := func(ctx context.Context, out string, state *testState) (string, error) {
		state.ms = append(state.ms, out)
		return out, nil
	}

	// Add node 时添加 Pre/Post Handler 来处理 State
	_ = sg.AddLambdaNode(nodeOfL1, l1,
		compose.WithStatePreHandler(l1StateToInput), compose.WithStatePostHandler(l1StateToOutput))

	l2 := compose.StreamableLambda(func(ctx context.Context, input string) (output *schema.StreamReader[string], err error) {
		
		// 在 Node 内部，用 ProcessState，传入一个读写 State 的 函数
		err = compose.ProcessState[*testState](ctx, func(_ context.Context, state *testState) error {
			fmt.Println("----- 在 node 内部调用 state -----")
			fmt.Println(state.ms)
			return nil
		})

		outStr := "StreamableLambda: " + input

		sr, sw := schema.Pipe[string](utf8.RuneCountInString(outStr))

		// nolint: byted_goroutine_recover
		go func() {
			for _, field := range strings.Fields(outStr) {
				sw.Send(field+" ", nil)
			}
			sw.Close()
		}()

		return sr, nil
	})

	l2StateToOutput := func(ctx context.Context, out string, state *testState) (string, error) {
		state.ms = append(state.ms, out)
		return out, nil
	}

	_ = sg.AddLambdaNode(nodeOfL2, l2, compose.WithStatePostHandler(l2StateToOutput))

	l3 := compose.TransformableLambda(func(ctx context.Context, input *schema.StreamReader[string]) (
		output *schema.StreamReader[string], err error) {

		prefix := "TransformableLambda: "
		sr, sw := schema.Pipe[string](20)

		go func() {

			defer func() {
				panicErr := recover()
				if panicErr != nil {
					fmt.Printf("panic occurs: %v\n", panicErr)
					fmt.Printf("stack: %v\n", string(debug.Stack()))
				}

			}()

			for _, field := range strings.Fields(prefix) {
				sw.Send(field+" ", nil)
			}

			for {
				chunk, err := input.Recv()
				if err != nil {
					if err == io.EOF {
						break
					}
					// TODO: how to trace this kind of error in the goroutine of processing sw
					sw.Send(chunk, err)
					break
				}

				sw.Send(chunk, nil)

			}
			sw.Close()
		}()

		return sr, nil
	})

	l3StateToOutput := func(ctx context.Context, out string, state *testState) (string, error) {
		state.ms = append(state.ms, out)
		fmt.Printf("state result: \n")
		for idx, m := range state.ms {
			fmt.Printf("    %vth: %v\n", idx, m)
		}
		return out, nil
	}

	_ = sg.AddLambdaNode(nodeOfL3, l3, compose.WithStatePostHandler(l3StateToOutput))

	_ = sg.AddEdge(compose.START, nodeOfL1)

	_ = sg.AddEdge(nodeOfL1, nodeOfL2)

	_ = sg.AddEdge(nodeOfL2, nodeOfL3)

	_ = sg.AddEdge(nodeOfL3, compose.END)

	run, err := sg.Compile(ctx)
	if err != nil {
		fmt.Printf("sg.Compile failed, err=%v", err)
		return
	}

	out, err := run.Invoke(ctx, "how are you")
	if err != nil {
		fmt.Printf("run.Invoke failed, err=%v", err)
		return
	}
	fmt.Printf("invoke result: %v\n", out)

	return

	// ----------------------------------------------------------
	stream, err := run.Stream(ctx, "how are you")
	if err != nil {
		fmt.Printf("run.Stream failed, err=%v", err)
		return
	}

	for {

		chunk, err := stream.Recv()
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			fmt.Printf("stream.Recv() failed, err=%v", err)
			break
		}

		fmt.Print(chunk)

		//    logs.Tokenf("%v", chunk)
	}
	stream.Close()

	// ----------------------------------------------------------

	sr, sw := schema.Pipe[string](1)
	sw.Send("how are you", nil)
	sw.Close()

	stream, err = run.Transform(ctx, sr)
	if err != nil {
		fmt.Printf("run.Transform failed, err=%v", err)
		return
	}

	for {

		chunk, err := stream.Recv()
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			fmt.Printf("stream.Recv() failed, err=%v", err)
			break
		}

		fmt.Printf("%v", chunk)
	}
	stream.Close()
}
