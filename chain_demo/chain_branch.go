package chain_demo

import (
	"context"
	"fmt"
	"io"
	"learn-eino/util"

	"github.com/cloudwego/eino/components/prompt"
	"github.com/cloudwego/eino/compose"
	"github.com/cloudwego/eino/schema"
)

func Main2() {
	ctx := context.Background()

	// 分支条件
	// 类型：func(ctx context.Context, in T) (endNode string, err error)
	//
	// branch 作为一个整体，conditon 的输入输出类型，要与各个分支 Lambda 的输入输出类型一致
	branchCondition := func(ctx context.Context, in map[string]any) (string, error) {
		switch in["language"] {
		case "golang":
			return "lambda_golang", nil
		case "python":
			return "lambda_python", nil
		default:
			return "lambda_php", nil
		}
	}
	
	branch1 := compose.InvokableLambda(func(ctx context.Context, in map[string]any) (map[string]any, error) {
		in["role"] = "Golang专家"
		return in, nil
	})
	branch2 := compose.InvokableLambda(func(ctx context.Context, in map[string]any) (map[string]any, error) {
		in["role"] = "Python专家"
		return in, nil
	})
	branch3 := compose.InvokableLambda(func(ctx context.Context, in map[string]any) (map[string]any, error) {
		in["role"] = "PHP专家"
		return in, nil
	})

	branch := compose.NewChainBranch(branchCondition).
		AddLambda("lambda_golang", branch1).
		AddLambda("lambda_python", branch2).
		AddLambda("lambda_php", branch3)

	chatModel, err := util.GetChatModel(ctx)
	if err != nil {
		panic(err)
	}

	chain := compose.NewChain[map[string]any, *schema.Message]()
	chain.AppendBranch(branch)
	chain.AppendChatTemplate(
		prompt.FromMessages(
			schema.FString,
			schema.SystemMessage("你是一个{role}。"),
			schema.MessagesPlaceholder("chat_history", true),
			schema.UserMessage("问题: {question}"),
		),
	)
	chain.AppendChatModel(chatModel)

	runable, err := chain.Compile(ctx)
	if err != nil {
		panic(err)
	}

	// param := map[string]any{
	// 	"language": "golang",
	// 	"question": "写一个冒泡算法",
	// }

	// param := map[string]any{
	// 	"language": "python",
	// 	"question": "写一个冒泡算法",
	// }

	param := map[string]any{
		"language": "php",
		"question": "写一个冒泡算法",
	}

	reader, err := runable.Stream(ctx, param)
	if err != nil {
		panic(err)
	}

	defer reader.Close()
	for {
		msg, err := reader.Recv()
		if err == io.EOF {
			return
		}
		if err != nil {
			panic(err)
		}
		fmt.Print(msg.Content)
	}
}