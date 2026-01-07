package chain_demo

import (
	"context"
	"fmt"
	"learn-eino/util"
	"log"
	"math/rand"

	"github.com/cloudwego/eino/components/prompt"
	"github.com/cloudwego/eino/compose"
	"github.com/cloudwego/eino/schema"
)

func Main4() {
	ctx := context.Background()
	// build branch func
	branchCond := func(ctx context.Context, input map[string]any) (string, error) { // nolint: byted_all_nil_return
		if rand.Intn(2) == 1 {
			return "b1", nil
		}

		return "b2", nil
	}

	b1 := compose.InvokableLambda(func(ctx context.Context, kvs map[string]any) (map[string]any, error) {
		fmt.Println("hello in branch lambda b1")
		if kvs == nil {
			return nil, fmt.Errorf("nil map")
		}

		kvs["role"] = "cat"
		return kvs, nil
	})

	b2 := compose.InvokableLambda(func(ctx context.Context, kvs map[string]any) (map[string]any, error) {
		fmt.Println("hello in branch lambda b2")
		if kvs == nil {
			return nil, fmt.Errorf("nil map")
		}

		kvs["role"] = "dog"
		return kvs, nil
	})

	parallel := compose.NewParallel()
	parallel.
		AddLambda("role", compose.InvokableLambda(func(ctx context.Context, kvs map[string]any) (string, error) {
			role, ok := kvs["role"].(string)
			if !ok || role == "" {
				role = "bird"
			}

			return role, nil
		})).
		AddLambda("input", compose.InvokableLambda(func(ctx context.Context, kvs map[string]any) (string, error) {
			return "你的叫声是怎样的？", nil
		}))

   // parallel 会收集多个 lambda 的输出，合并到一个 map[string]any 中，作为下一个 lambda 的输入

	cm, err := util.GetChatModel(ctx)
	if err != nil {
		log.Panic(err)
		return
	}

	rolePlayerChain := compose.NewChain[map[string]any, *schema.Message]()
	rolePlayerChain.
		AppendLambda(compose.InvokableLambda(func(ctx context.Context, kvs map[string]any) (map[string]any, error) {
			// 查看 parallel 输出类型
         for k, v := range kvs {
				fmt.Printf("%v: %v\n", k, v)
			}
			return kvs, nil
		})).
		AppendChatTemplate(
			prompt.FromMessages(
				schema.FString,
				schema.SystemMessage(`You are a {role}.`),
				schema.UserMessage(`{input}`),
			),
		).
		AppendChatModel(cm)

	// =========== build chain ===========
	chain := compose.NewChain[map[string]any, string]()
	chain.
		AppendLambda(compose.InvokableLambda(func(ctx context.Context, kvs map[string]any) (map[string]any, error) {
			return kvs, nil
		})).
		AppendBranch(compose.NewChainBranch(branchCond).AddLambda("b1", b1).AddLambda("b2", b2)).
		AppendPassthrough().
		AppendParallel(parallel).
		AppendGraph(rolePlayerChain).
		AppendLambda(compose.InvokableLambda(func(ctx context.Context, m *schema.Message) (string, error) {
			return m.Content, nil
		}))

	r, err := chain.Compile(ctx)
	if err != nil {
		log.Panic(err)
		return
	}

	output, err := r.Invoke(ctx, map[string]any{})
	if err != nil {
		log.Panic(err)
		return
	}

	fmt.Printf("%v\n", output)
}
