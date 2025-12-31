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

// Chain 可以视为是 Graph 的简化封装

func Main() {
	Main4()
}

func Main1() {
	ctx := context.Background()

	model, err := util.GetChatModel(ctx)
	if err != nil {
		panic(err)
	}

	// 定义 chain 第一个节点的输入类型，和最后一个节点的输出类型
	// 中间各个节点的类型，要满足上游的输出类型与下游的输入类型一致
	chain := compose.NewChain[map[string]any, *schema.Message]().
		AppendChatTemplate(
			prompt.FromMessages(schema.FString,
				schema.SystemMessage("你是一个{role}。"),
				schema.MessagesPlaceholder("chat_history", true),
				schema.UserMessage("问题: {question}"),
			),
		).
		AppendChatModel(model)

	runable, err := chain.Compile(ctx)
	if err != nil {
		panic(err)
	}

	in := map[string]any{
		"role":     "Golang专家",
		"question": "用go语言开发AI智能体有优势吗？",
	}

	reader, err := runable.Stream(ctx, in)
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
