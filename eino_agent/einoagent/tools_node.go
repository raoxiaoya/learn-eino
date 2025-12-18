package einoagent

import (
	"context"
	"time"

	"github.com/cloudwego/eino-ext/components/tool/duckduckgo/v2"
	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/schema"
)

func newTool(ctx context.Context) (bt tool.BaseTool, err error) {
	// TODO Modify component configuration here.
	config := &duckduckgo.Config{
		MaxResults: 3, // Limit to return 20 results
		Region:     duckduckgo.RegionWT,
		Timeout:    10 * time.Second,
	}
	bt, err = duckduckgo.NewTextSearchTool(ctx, config)
	if err != nil {
		return nil, err
	}
	return bt, nil
}

type Tool1Impl struct {
	config *Tool1Config
}

type Tool1Config struct {
}

func newTool1(ctx context.Context) (bt tool.BaseTool, err error) {
	// TODO Modify component configuration here.
	config := &Tool1Config{}
	bt = &Tool1Impl{config: config}
	return bt, nil
}

func (impl *Tool1Impl) Info(ctx context.Context) (*schema.ToolInfo, error) {
	return &schema.ToolInfo{
		Name: "tool1",
		Desc: "Tool1 desc",
	}, nil
}

func (impl *Tool1Impl) InvokableRun(ctx context.Context, argumentsInJSON string, opts ...tool.Option) (string, error) {
	return "tool1 result", nil
}

type Tool2Impl struct {
	config *Tool2Config
}

type Tool2Config struct {
}

func newTool2(ctx context.Context) (bt tool.BaseTool, err error) {
	// TODO Modify component configuration here.
	config := &Tool2Config{}
	bt = &Tool2Impl{config: config}
	return bt, nil
}

func (impl *Tool2Impl) Info(ctx context.Context) (*schema.ToolInfo, error) {
	return &schema.ToolInfo{
		Name: "tool2",
		Desc: "Tool2 desc",
	}, nil
}

func (impl *Tool2Impl) InvokableRun(ctx context.Context, argumentsInJSON string, opts ...tool.Option) (string, error) {
	return "tool2 result", nil
}

type Tool3Impl struct {
	config *Tool3Config
}

type Tool3Config struct {
}

func newTool3(ctx context.Context) (bt tool.BaseTool, err error) {
	// TODO Modify component configuration here.
	config := &Tool3Config{}
	bt = &Tool3Impl{config: config}
	return bt, nil
}

func (impl *Tool3Impl) Info(ctx context.Context) (*schema.ToolInfo, error) {
	return &schema.ToolInfo{
		Name: "tool3",
		Desc: "Tool3 desc",
	}, nil
}

func (impl *Tool3Impl) InvokableRun(ctx context.Context, argumentsInJSON string, opts ...tool.Option) (string, error) {
	return "tool3 result", nil
}

type Tool4Impl struct {
	config *Tool4Config
}

type Tool4Config struct {
}

func newTool4(ctx context.Context) (bt tool.BaseTool, err error) {
	// TODO Modify component configuration here.
	config := &Tool4Config{}
	bt = &Tool4Impl{config: config}
	return bt, nil
}

func (impl *Tool4Impl) Info(ctx context.Context) (*schema.ToolInfo, error) {
	return &schema.ToolInfo{
		Name: "tool4",
		Desc: "Tool4 desc",
	}, nil
}

func (impl *Tool4Impl) InvokableRun(ctx context.Context, argumentsInJSON string, opts ...tool.Option) (string, error) {
	return "tool4 result", nil
}
