package util

import (
	"context"
	"os"

	"github.com/cloudwego/eino-ext/components/embedding/ark"
)

func EmbeddingString(ctx context.Context, content []string) ([][]float64, error) {
	ebs, err := ark.NewEmbedder(ctx, &ark.EmbeddingConfig{
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