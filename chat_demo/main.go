// eino-examples/quickstart/chat/template.go
package chat_demo

import (
	"context"
	"fmt"
	"io"
	"log"

	"learn-eino/util"
	"learn-eino/util/env"

	"github.com/cloudwego/eino/components/prompt"
	"github.com/cloudwego/eino/schema"
)

func init() {
	env.MustHasEnvs("ARK_BASE_URL", "ARK_API_KEY", "ARK_CHAT_MODEL")
}

func buildMessage(ctx context.Context) ([]*schema.Message, error) {
	// 创建模板，使用 FString 格式
	template := prompt.FromMessages(
		// 模版类型
		schema.FString,

		// 系统消息模板
		schema.SystemMessage("你是一个{role}。你需要用{style}的语气回答问题。你的目标是帮助程序员保持积极乐观的心态，提供技术建议的同时也要关注他们的心理健康。"),

		// 插入需要的对话历史（新对话的话这里不填）
		schema.MessagesPlaceholder("chat_history", true),

		// 用户消息模板
		schema.UserMessage("问题: {question}"),
	)

	// 使用模板生成消息
	messages, err := template.Format(ctx, map[string]any{
		"role":     "程序员鼓励师",
		"style":    "积极、温暖且专业",
		"question": "我的代码一直报错，感觉好沮丧，该怎么办？",
		// 对话历史（这个例子里模拟两轮对话历史）
		"chat_history": []*schema.Message{
			schema.UserMessage("你好"),
			schema.AssistantMessage("嘿！我是你的程序员鼓励师！记住，每个优秀的程序员都是从 Debug 中成长起来的。有什么我可以帮你的吗？", nil),
			schema.UserMessage("我觉得自己写的代码太烂了"),
			schema.AssistantMessage("每个程序员都经历过这个阶段！重要的是你在不断学习和进步。让我们一起看看代码，我相信通过重构和优化，它会变得更好。记住，Rome wasn't built in a day，代码质量是通过持续改进来提升的。	", nil),
		},
	})

	return messages, err
}

func testChat(ctx context.Context) {
	model, err := util.GetChatModel(ctx)
	if err != nil {
		panic(err)
	}

	messages, err := buildMessage(ctx)
	result, err := model.Generate(ctx, messages)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%+v\n", result)
}

func testChatStream(ctx context.Context) {
	model, err := util.GetChatModel(ctx)
	if err != nil {
		panic(err)
	}
	messages, err := buildMessage(ctx)
	result, err := model.Stream(ctx, messages)
	if err != nil {
		panic(err)
	}
	defer result.Close()

	for {
		message, err := result.Recv()
		if err == io.EOF { // 流式输出结束
			return
		}
		if err != nil {
			log.Fatalf("recv failed: %v", err)
		}
		fmt.Print(message.Content)
	}
}

func Main() {
	ctx := context.Background()

	// testChat(ctx)
	testChatStream(ctx)
}
