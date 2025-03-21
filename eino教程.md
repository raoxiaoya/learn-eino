cloudwego - eino 使用教程



文档开始：https://www.cloudwego.io/zh/docs/eino/

github：https://github.com/cloudwego/eino



### 一、简介

**Eino[‘aino]** (近似音: i know，希望应用程序达到 “i know” 的愿景) 旨在提供基于 Golang 语言的终极大模型应用开发框架。 它从开源社区中的诸多优秀 LLM 应用开发框架，如 LangChain 和 LlamaIndex 等获取灵感，同时借鉴前沿研究成果与实际应用，提供了一个强调简洁性、可扩展性、可靠性与有效性，且更符合 Go 语言编程惯例的 LLM 应用开发框架。

```bash
go get -u github.com/cloudwego/eino@v0.3.17
go get -u github.com/cloudwego/eino-ext/components/model/qwen

注意kin-openapi要使用v0.118.0版本，否则会报错
github.com/getkin/kin-openapi v0.118.0
```

示例参考：[eino-example](github.com/cloudwego/eino-example)

### 二、实践

```go
package main

import (
	"context"
	"fmt"
	"io"
	"log"

	"github.com/cloudwego/eino-ext/components/model/qwen"
	"github.com/cloudwego/eino/components/prompt"
	"github.com/cloudwego/eino/schema"
)

func buildMessage() ([]*schema.Message, error) {
	// 创建模板，使用 FString 格式
	template := prompt.FromMessages(schema.FString,
		// 系统消息模板
		schema.SystemMessage("你是一个{role}。你需要用{style}的语气回答问题。你的目标是帮助程序员保持积极乐观的心态，提供技术建议的同时也要关注他们的心理健康。"),

		// 插入需要的对话历史（新对话的话这里不填）
		schema.MessagesPlaceholder("chat_history", true),

		// 用户消息模板
		schema.UserMessage("问题: {question}"),
	)

	// 使用模板生成消息
	messages, err := template.Format(context.Background(), map[string]any{
		"role":     "程序员鼓励师",
		"style":    "积极、温暖且专业",
		"question": "我的代码一直报错，感觉好沮丧，该怎么办？",
		// 对话历史（这个例子里模拟两轮对话历史）
		"chat_history": []*schema.Message{
			schema.UserMessage("你好"),
			schema.AssistantMessage("嘿！我是你的程序员鼓励师！记住，每个优秀的程序员都是从 Debug 中成长起来的。有什么我可以帮你的吗？", nil),
			schema.UserMessage("我觉得自己写的代码太烂了"),
			schema.AssistantMessage("每个程序员都经历过这个阶段！重要的是你在不断学习和进步。让我们一起看看代码，我相信通过重构和优化，它会变得更好。记住，Rome wasn't built in a day，代码质量是通过持续改进来提升的。", nil),
		},
	})

	return messages, err
}

func getQwenChatModel() (*qwen.ChatModel, error) {
	chatModel, err := qwen.NewChatModel(context.Background(), &qwen.ChatModelConfig{
		BaseURL: "https://dashscope.aliyuncs.com/compatible-mode/v1",
		Model:   "qwen-max",
		APIKey:  "xxxxxxxxxxxxxxxxxxx",
	})
	return chatModel, err
}

func testQwen() {
	model, err := getQwenChatModel()
	if err != nil {
		panic(err)
	}

	messages, err := buildMessage()
	result, err := model.Generate(context.Background(), messages)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%+v\n", result)
}

func testQwenStream() {
	model, err := getQwenChatModel()
	if err != nil {
		panic(err)
	}
	messages, err := buildMessage()
	result, err := model.Stream(context.Background(), messages)
	if err != nil {
		panic(err)
	}
	reportStream(result)
}

func reportStream(sr *schema.StreamReader[*schema.Message]) {
	defer sr.Close()

	i := 0
	for {
		message, err := sr.Recv()
		if err == io.EOF {
			return
		}
		if err != nil {
			log.Fatalf("recv failed: %v", err)
		}
		log.Printf("message[%d]: %+v\n", i, message)
		i++
	}
}

func main() {
	testQwenStream()
}

```



