package main

import (
	"context"
	"encoding/json"
	"log"
	"os"

	"github.com/cloudwego/eino/schema"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types"

	"github.com/cloudwego/eino-ext/components/embedding/ark"
	"github.com/cloudwego/eino-ext/components/retriever/es8"
	"github.com/cloudwego/eino-ext/components/retriever/es8/search_mode"
)

func main() {
	ctx := context.Background()

	// es 支持多种连接方式
	// username := os.Getenv("ES_USERNAME")
	// password := os.Getenv("ES_PASSWORD")
	// httpCACertPath := os.Getenv("ES_HTTP_CA_CERT_PATH")

	// cert, err := os.ReadFile(httpCACertPath)
	// if err != nil {
	//         log.Fatalf("read file failed, err=%v", err)
	// }

	client, err := elasticsearch.NewClient(elasticsearch.Config{
		Addresses: []string{"http://localhost:9200"},
		// Username:  username,
		// Password:  password,
		// CACert:    cert,
	})
	if err != nil {
		log.Panicf("connect es8 failed, err=%v", err)
	}

	emb, err := ark.NewEmbedder(ctx, &ark.EmbeddingConfig{
		BaseURL: os.Getenv("ARK_BASE_URL"),
		APIKey:  os.Getenv("ARK_API_KEY"),
		Model:   os.Getenv("ARK_EMBEDDING_MODEL"),
	})
	if err != nil {
		panic(err)
	}

	// 创建检索器组件
	k_value := 10
	retriever, err := es8.NewRetriever(ctx, &es8.RetrieverConfig{
		Client: client,
		Index:  "eino_example",
		SearchMode: search_mode.SearchModeApproximate(&search_mode.ApproximateConfig{
			QueryFieldName:  "content",
			VectorFieldName: "content_vector",
			Hybrid:          true,
			// RRF 仅在特定许可证下可用
			// 参见：https://www.elastic.co/subscriptions
			RRF:             false,
			RRFRankConstant: nil,
			RRFWindowSize:   nil,
			K:               &k_value,
			NumCandidates:   &k_value,
		}),
		ResultParser: func(ctx context.Context, hit types.Hit) (doc *schema.Document, err error) {
			doc = &schema.Document{
				ID:       *hit.Id_,
				Content:  "",
				MetaData: map[string]any{},
			}

			var src map[string]any
			if err = json.Unmarshal(hit.Source_, &src); err != nil {
				return nil, err
			}

			for field, val := range src {
				switch field {
				case "content":
					doc.Content = val.(string)
				case "content_vector":
					var v []float64
					for _, item := range val.([]interface{}) {
						v = append(v, item.(float64))
					}
					doc.WithDenseVector(v)
				case "location":
					if loc, ok := val.(string); ok {
						doc.MetaData["location"] = loc
					} else {
						doc.MetaData["location"] = ""
					}
				}
			}

			if hit.Score_ != nil {
				doc.WithScore(float64(*hit.Score_))
			}

			return doc, nil
		},
		Embedding: emb, // 你的 embedding 组件
	})
	if err != nil {
		log.Panicf("create retriever failed, err=%v", err)
	}

	// 无过滤条件搜索
	docs, err := retriever.Retrieve(ctx, "tourist attraction")
	if err != nil {
		log.Panicf("retrieve docs failed1, err=%v", err)
	}
	for _, doc := range docs {
		log.Printf("doc1=%v\n", doc.String())
	}

	// 带过滤条件搜索
	trueof := true
	docs, err = retriever.Retrieve(ctx, "tourist attraction",
		es8.WithFilters([]types.Query{{
			Term: map[string]types.TermQuery{
				"location": {
					CaseInsensitive: &trueof,
					Value:           "China",
				},
			},
		}}),
	)
	if err != nil {
		log.Panicf("retrieve docs failed2, err=%v", err)
	}
	for _, doc := range docs {
		log.Printf("doc2=%v\n", doc.String())
	}
}
