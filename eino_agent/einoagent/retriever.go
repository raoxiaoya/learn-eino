package einoagent

import (
	"context"
	"encoding/json"
	"log"

	"github.com/cloudwego/eino-ext/components/retriever/es8"
	"github.com/cloudwego/eino-ext/components/retriever/es8/search_mode"
	"github.com/cloudwego/eino/components/retriever"
	"github.com/cloudwego/eino/schema"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types"
)

// newRetriever component initialization function of node 'Retriever1' in graph 'EinoAgentEs8'
func newRetriever(ctx context.Context) (rtr retriever.Retriever, err error) {
	// TODO Modify component configuration here.
	client, err := elasticsearch.NewClient(elasticsearch.Config{
		Addresses: []string{"http://localhost:9200"},
	})
	if err != nil {
		log.Panicf("connect es8 failed, err=%v", err)
	}
	k_value := 10
	config := &es8.RetrieverConfig{
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
	}
	embeddingIns11, err := newEmbedding(ctx)
	if err != nil {
		return nil, err
	}
	config.Embedding = embeddingIns11
	rtr, err = es8.NewRetriever(ctx, config)
	if err != nil {
		return nil, err
	}
	return rtr, nil
}
