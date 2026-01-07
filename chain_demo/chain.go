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

// Chain 可以视为是 Graph 的简化封装。
// chain 根据 append 的先后顺序来编排，当然，底层还是会转换成 graph 来执行。
// graph 则是先 addNode，然后 addEdge 来编排，graph 更加灵活。
// 从结果来看，chain 只能是顺序结构，而 graph 则可以通过 addEdge 来实现任意结构。
func Main() {
	Main8()
}

func Main1() {
	ctx := context.Background()

	model, err := util.GetChatModel(ctx)
	if err != nil {
		panic(err)
	}

	// 定义 chain 第一个节点的输入类型，和最后一个节点的输出类型
	// 中间各个节点的类型，要满足上游的输出类型与下游的输入类型一致
	//
	// ChatTemplate 节点，input: map[string]any, output: []*schema.Message
	// ChatModel 节点，input: []*schema.Message, output: *schema.Message
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
