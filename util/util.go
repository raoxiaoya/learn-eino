package util

import (
	"context"
	"errors"
	"io"
	"learn-eino/util/env"
	"os"

	"github.com/cloudwego/eino/schema"

	embeddingArk "github.com/cloudwego/eino-ext/components/embedding/ark"
	"github.com/cloudwego/eino-ext/components/model/ark"
)

func init() {
	env.MustHasEnvs("ARK_BASE_URL", "ARK_API_KEY", "ARK_CHAT_MODEL", "ARK_EMBEDDING_MODEL")
}

func EmbeddingString(ctx context.Context, content []string) ([][]float64, error) {
	ebs, err := embeddingArk.NewEmbedder(ctx, &embeddingArk.EmbeddingConfig{
		BaseURL: os.Getenv("ARK_BASE_URL"),
		APIKey:  os.Getenv("ARK_API_KEY"),
		Model:   os.Getenv("ARK_EMBEDDING_MODEL"),
	})
	if err != nil {
		return nil, err
	}

	res, err := ebs.EmbedStrings(ctx, content)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func GetChatModel(ctx context.Context) (*ark.ChatModel, error) {
	chatModel, err := ark.NewChatModel(ctx, &ark.ChatModelConfig{
		BaseURL: os.Getenv("ARK_API_URL"),
		APIKey:  os.Getenv("ARK_API_KEY"),
		Model:   os.Getenv("ARK_CHAT_MODEL"),
	})
	return chatModel, err
}

func CustomToolCallChecker(ctx context.Context, sr *schema.StreamReader[*schema.Message]) (bool, error) {
	defer sr.Close()
	for {
		msg, err := sr.Recv()
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			return false, err
		}
		if len(msg.ToolCalls) > 0 {
			return true, nil
		}
	}
	return false, nil
}
