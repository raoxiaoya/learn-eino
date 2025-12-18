package indexeres8

import (
	"context"
	"log"

	"github.com/cloudwego/eino-ext/components/indexer/es8"
	"github.com/cloudwego/eino/components/indexer"
	"github.com/cloudwego/eino/schema"
	"github.com/elastic/go-elasticsearch/v8"
)

// newIndexer component initialization function of node 'Indexer1' in graph 'es8Indexer'
func newIndexer(ctx context.Context) (idr indexer.Indexer, err error) {
	// TODO Modify component configuration here.
	client, err := elasticsearch.NewClient(elasticsearch.Config{
		Addresses: []string{"http://localhost:9200"},
		// Username:  username,
		// Password:  password,
		// CACert:    cert,
	})
	if err != nil {
		log.Panicf("connect es8 failed, err=%v", err)
	}

	config := &es8.IndexerConfig{
		Index:     "eino_example",
		BatchSize: 10,
		Client:    client,
		DocumentToFields: func(ctx context.Context, doc *schema.Document) (field2Value map[string]es8.FieldValue, err error) {
			return map[string]es8.FieldValue{
				"content": {
					Value:    doc.Content,
					EmbedKey: "content_vector", // 对文档内容进行向量化并保存向量到 "content_vector" 字段
				},
				"location": {
					Value: doc.MetaData["location"],
				},
			}, nil
		},
	}
	embeddingIns11, err := newEmbedding(ctx)
	if err != nil {
		return nil, err
	}
	config.Embedding = embeddingIns11
	idr, err = es8.NewIndexer(ctx, config)
	if err != nil {
		return nil, err
	}
	return idr, nil
}
