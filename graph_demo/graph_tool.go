package graph_demo

import (
	"context"
	"fmt"
	"learn-eino/util"

	"github.com/cloudwego/eino/callbacks"
	"github.com/cloudwego/eino/components/prompt"
	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/components/tool/utils"
	"github.com/cloudwego/eino/compose"
	"github.com/cloudwego/eino/schema"
)

func Main2() {
	ctx := context.Background()

	callbacks.AppendGlobalHandlers(&loggerCallbacks{})

	// 1. create an instance of ChatTemplate as 1st Graph Node
	chatTpl := prompt.FromMessages(schema.FString,
		schema.SystemMessage(`你是一名房产经纪人，结合用户的薪酬和工作，使用 user_info API，为其提供相关的房产信息。邮箱是必须的`),
		schema.MessagesPlaceholder("message_histories", true),
		schema.UserMessage("{user_query}"),
	)

	chatModel, err := util.GetChatModel(ctx)
	if err != nil {
		fmt.Printf("NewChatModel failed, err=%v\n", err)
		return
	}

	// 3. create an instance of tool.InvokableTool for Intent recognition and execution
	userInfoTool := utils.NewTool(
		&schema.ToolInfo{
			Name: "user_info",
			Desc: "根据用户的姓名和邮箱，查询用户的公司、职位、薪酬信息",
			ParamsOneOf: schema.NewParamsOneOfByParams(map[string]*schema.ParameterInfo{
				"name": {
					Type: "string",
					Desc: "用户的姓名",
				},
				"email": {
					Type: "string",
					Desc: "用户的邮箱",
				},
			}),
		},
		func(ctx context.Context, input *userInfoRequest) (output *userInfoResponse, err error) {
			// LLM 将参数以 map 的 json 格式返回，agent 将 json string 反序列化为 userInfoRequest struct
			fmt.Println("-------------调用tool-------------")
			return &userInfoResponse{
				Name:     input.Name,
				Email:    input.Email,
				Company:  "Bytedance",
				Position: "CEO",
				Salary:   "9999",
			}, nil
		})

	info, err := userInfoTool.Info(ctx)
	if err != nil {
		fmt.Printf("Get ToolInfo failed, err=%v\n", err)
		return
	}

	// 4. bind ToolInfo to ChatModel. ToolInfo will remain in effect until the next BindTools.
	err = chatModel.BindForcedTools([]*schema.ToolInfo{info})
	if err != nil {
		fmt.Printf("BindForcedTools failed, err=%v\n", err)
		return
	}

	// 5. create an instance of ToolsNode as 3rd Graph Node
	toolsNode, err := compose.NewToolNode(ctx, &compose.ToolsNodeConfig{
		Tools: []tool.BaseTool{userInfoTool},
	})
	if err != nil {
		fmt.Printf("NewToolNode failed, err=%v\n", err)
		return
	}

	const (
		nodeKeyOfTemplate  = "template"
		nodeKeyOfChatModel = "chat_model"
		nodeKeyOfTools     = "tools"
	)

	// 6. create an instance of Graph
	// input type is 1st Graph Node's input type, that is ChatTemplate's input type: map[string]any
	// output type is last Graph Node's output type, that is ToolsNode's output type: []*schema.Message
	g := compose.NewGraph[map[string]any, []*schema.Message]()

	// 7. add ChatTemplate into graph
	_ = g.AddChatTemplateNode(nodeKeyOfTemplate, chatTpl)

	// 8. add ChatModel into graph
	_ = g.AddChatModelNode(nodeKeyOfChatModel, chatModel)

	// 9. add ToolsNode into graph
	_ = g.AddToolsNode(nodeKeyOfTools, toolsNode)

	// 10. add connection between nodes
	_ = g.AddEdge(compose.START, nodeKeyOfTemplate)

	_ = g.AddEdge(nodeKeyOfTemplate, nodeKeyOfChatModel)

	_ = g.AddEdge(nodeKeyOfChatModel, nodeKeyOfTools)

	_ = g.AddEdge(nodeKeyOfTools, compose.END)

	// 9. compile Graph[I, O] to Runnable[I, O]
	r, err := g.Compile(ctx)
	if err != nil {
		fmt.Printf("Compile failed, err=%v\n", err)
		return
	}

	out, err := r.Invoke(ctx, map[string]any{
		"message_histories": []*schema.Message{},
		"user_query":        "我叫 zhangsan, 邮箱是 zhangsan@bytedance.com, 帮我推荐一处房产",
	})
	if err != nil {
		fmt.Printf("Invoke failed, err=%v\n", err)
		return
	}
	fmt.Printf("Generation: %v Messages\n", len(out))
	for _, msg := range out {
		fmt.Printf("    %v\n", msg)
	}
}

type userInfoRequest struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

type userInfoResponse struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Company  string `json:"company"`
	Position string `json:"position"`
	Salary   string `json:"salary"`
}

type loggerCallbacks struct{}

func (l *loggerCallbacks) OnStart(ctx context.Context, info *callbacks.RunInfo, input callbacks.CallbackInput) context.Context {
	fmt.Printf("name: %v, type: %v, component: %v, input: %v\n", info.Name, info.Type, info.Component, input)
	return ctx
}

func (l *loggerCallbacks) OnEnd(ctx context.Context, info *callbacks.RunInfo, output callbacks.CallbackOutput) context.Context {
	fmt.Printf("name: %v, type: %v, component: %v, output: %v\n", info.Name, info.Type, info.Component, output)
	return ctx
}

func (l *loggerCallbacks) OnError(ctx context.Context, info *callbacks.RunInfo, err error) context.Context {
	fmt.Printf("name: %v, type: %v, component: %v, error: %v\n", info.Name, info.Type, info.Component, err)
	return ctx
}

func (l *loggerCallbacks) OnStartWithStreamInput(ctx context.Context, info *callbacks.RunInfo, input *schema.StreamReader[callbacks.CallbackInput]) context.Context {
	return ctx
}

func (l *loggerCallbacks) OnEndWithStreamOutput(ctx context.Context, info *callbacks.RunInfo, output *schema.StreamReader[callbacks.CallbackOutput]) context.Context {
	return ctx
}
