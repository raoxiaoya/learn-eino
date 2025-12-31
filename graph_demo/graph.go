package graph_demo

import (
	"context"
	"fmt"
	"io"
	"learn-eino/util"

	"github.com/cloudwego/eino/components/model"
	"github.com/cloudwego/eino/components/prompt"
	"github.com/cloudwego/eino/compose"
	"github.com/cloudwego/eino/schema"
)

func Main() {
	Main3()
}

func Main1() {
	ctx := context.Background()
	g := compose.NewGraph[map[string]any, *schema.Message]()

	pt := prompt.FromMessages(
		schema.FString,
		schema.UserMessage("what's the weather in {location}?"),
	)

	chatModel, err := util.GetChatModel(ctx)
	if err != nil {
		panic(err)
	}

	_ = g.AddChatTemplateNode("prompt", pt)
	_ = g.AddChatModelNode("model", chatModel, compose.WithNodeName("ChatModel"))
	_ = g.AddEdge(compose.START, "prompt")
	_ = g.AddEdge("prompt", "model")
	_ = g.AddEdge("model", compose.END)

	r, err := g.Compile(ctx)
	if err != nil {
		panic(err)
	}

	in := map[string]any{"location": "beijing"}
	s, err := r.Stream(ctx, in)
	if err != nil {
		panic(err)
	}

	defer s.Close()
	for {
		chunk, err := s.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			panic(err)
		}

		fmt.Print(chunk.Content)
	}
}

type mockChatModel struct{}

func (m *mockChatModel) Generate(ctx context.Context, input []*schema.Message, opts ...model.Option) (*schema.Message, error) {
	return schema.AssistantMessage("the weather is good", nil), nil
}

func (m *mockChatModel) Stream(ctx context.Context, input []*schema.Message, opts ...model.Option) (*schema.StreamReader[*schema.Message], error) {
	sr, sw := schema.Pipe[*schema.Message](0)
	go func() {
		defer sw.Close()
		sw.Send(schema.AssistantMessage("the weather is", nil), nil)
		sw.Send(schema.AssistantMessage("good", nil), nil)
	}()
	return sr, nil
}

func (m *mockChatModel) BindTools(tools []*schema.ToolInfo) error {
	panic("implement me")
}
