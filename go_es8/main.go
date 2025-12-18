package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"learn-eino/util"
	"log"
	"strings"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
)

var client *elasticsearch.Client

func init() {
	cli, err := elasticsearch.NewClient(elasticsearch.Config{
		Addresses: []string{"http://localhost:9200"},
	})
	if err != nil {
		log.Panicf("connect es8 failed, err=%v", err)
	}
	client = cli
}

func main() {
	ctx := context.Background()

	//// 1、创建 index
	// createIndex("my_index_vector")

	//// 2、索引文档
	// content := "Eino 旨在提供 Golang 语言的 AI 应用开发框架。 Eino 参考了开源社区中诸多优秀的 AI 应用开发框架，例如 LangChain、LangGraph、LlamaIndex 等，提供了更符合 Golang 编程习惯的 AI 应用开发框架。"

	// res, err := util.EmbeddingString(ctx, []string{content})
	// if err != nil {
	// 	panic(err)
	// }

	// fmt.Println("info: ", len(res), len(res[0]))

	// indexDocument("my_index_vector", "1", Document{
	// 	Title:     "Eino是什么",
	// 	Content:   content,
	// 	Embedding: res[0],
	// })

	//// 3、KNN 搜索
	res2, err := util.EmbeddingString(ctx, []string{"eino能干什么"})
	if err != nil {
		panic(err)
	}
	knnSearch("my_index_vector", res2[0], 5)
}

func createIndex(indexName string) {
	mapping := `{
		"mappings": {
			"properties": {
				"title": { "type": "text" },
				"content": { "type": "text" },
				"embedding": {
					"type": "dense_vector",
					"dims": 2048,
					"index": true,
					"similarity": "cosine"
				}
			}
		}
	}`

	req := esapi.IndicesCreateRequest{
		Index: indexName,
		Body:  strings.NewReader(mapping),
	}

	res, err := req.Do(context.Background(), client)
	if err != nil {
		log.Fatalf("Error creating index: %s", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		log.Printf("Error: %s", res.String())
	} else {
		fmt.Println("Index created successfully")
	}
}

type Document struct {
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	Embedding []float64 `json:"embedding"`
}

func indexDocument(indexName, id string, doc Document) {
	data, _ := json.Marshal(doc)

	req := esapi.IndexRequest{
		Index:      indexName,
		DocumentID: id,
		Body:       strings.NewReader(string(data)),
		Refresh:    "true",
	}

	res, err := req.Do(context.Background(), client)
	if err != nil {
		log.Fatalf("Error indexing document: %s", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		log.Printf("Error: %s", res.String())
	} else {
		fmt.Printf("Document %s indexed\n", id)
	}
}

func knnSearch(indexName string, queryVector []float64, k int) {
	query := map[string]interface{}{
		"knn": map[string]interface{}{
			"field":          "embedding",
			"query_vector":   queryVector,
			"k":              k,
			"num_candidates": k * 10, // 候选数量（建议 k*10）
		},
	}

	body, _ := json.Marshal(query)

	req := esapi.SearchRequest{
		Index: []string{indexName},
		Body:  bytes.NewReader(body),
	}

	res, err := req.Do(context.Background(), client)
	if err != nil {
		log.Fatalf("Error searching: %s", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		log.Printf("Search error: %s", res.String())
		return
	}

	var r map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		log.Fatalf("Error parsing response: %s", err)
	}

	// 打印结果
	hits := r["hits"].(map[string]interface{})["hits"].([]interface{})
	for _, hit := range hits {
		source := hit.(map[string]interface{})["_source"]
		score := hit.(map[string]interface{})["_score"]
		fmt.Printf("Score: %f, Title: %s\n", score, source.(map[string]interface{})["title"])
	}
	// Score: 0.834011, Title: Eino是什么
}

func testes() {
	client.Indices.Create("my_index")

	document := struct {
		Name string `json:"name"`
	}{
		"go-elasticsearch",
	}
	data, _ := json.Marshal(document)
	client.Index("my_index", bytes.NewReader(data))
}
