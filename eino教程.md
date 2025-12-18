cloudwego - eino 使用教程



### 参考

文档开始：https://www.cloudwego.io/zh/docs/eino/

eino：https://github.com/cloudwego/eino

eino-examples：[https://github.com/cloudwego/eino-examples](https://github.com/cloudwego/eino-examples/tree/main/quickstart/eino_assistant)



### Qwen模型

dashscope 的 key：https://dashscope.console.aliyun.com/apiKey，后来升级为阿里云百炼，但是系统是打通的，key 也是一样的。

阿里云百炼的 key：https://bailian.console.aliyun.com/?tab=model#/api-key，是阿里云的正式服务。

魔搭社区(modelscope) 的 key：https://www.modelscope.cn/my/myaccesstoken，用于调用魔搭 [API-Inference](https://www.modelscope.cn/docs/model-service/API-Inference/intro) 等其他服务。阿里云会将魔搭社区的一些模型进行部署以提供大家免费调用试玩，资源有限，非正式服务。



qwen-max详情：https://bailian.console.aliyun.com/?tab=model#/model-market/detail/qwen-max

在模型详情页可以查看免费额度。



**Ark**是字节跳动旗下的火山引擎云服务平台。



### 一、简介

**Eino[‘aino]** (近似音: i know，希望应用程序达到 “i know” 的愿景) 旨在提供基于 Golang 语言的终极大模型应用开发框架。 它从开源社区中的诸多优秀 LLM 应用开发框架，如 LangChain 和 LlamaIndex 等获取灵感，同时借鉴前沿研究成果与实际应用，提供了一个强调简洁性、可扩展性、可靠性与有效性，且更符合 Go 语言编程惯例的 LLM 应用开发框架。

```bash
go get -u github.com/cloudwego/eino@v0.3.17
go get -u github.com/cloudwego/eino-ext/components/model/qwen

注意kin-openapi要使用v0.118.0版本，否则会报错
github.com/getkin/kin-openapi v0.118.0
```



### 二、chat 实践

```go
package main

import (
	"context"
	"fmt"
	"io"
	"log"
    "os"

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
		APIKey:  os.Getenv("DASHSCOPE_API_KEY"),
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



### 三、接通 Elasticsearch

可能是本机有点老旧，导致 docker redis-stack 无法正常运行，于是改用 ES8

文档：https://www.cloudwego.io/zh/docs/eino/ecosystem_integration/indexer/indexer_es8/



使用 elasticsearch-8.10.4（[下载地址](https://artifacts.elastic.co/downloads/elasticsearch/elasticsearch-8.10.4-windows-x86_64.zip)）

```bash
D:\dev\php\magook\trunk\server\elasticsearch-8.10.4\bin

elasticsearch.bat
```

```bash
D:\dev\php\magook\trunk\server\elasticsearch-head

npm run start
```

访问：http://127.0.0.1:9200/

访问：http://localhost:9100/



使用火山引擎 https://console.volcengine.com/ark

查看[模型列表](https://www.volcengine.com/docs/82379/1330310?lang=zh)

向量化模型：模型名称`doubao-embedding-large`，对应的模型ID为`doubao-embedding-large-text-250515`，维度2048

大语言模型：`deepseek-v3.2`，对应的模型ID为`deepseek-v3-2-251201`



在API KEY管理创建一个API KEY



在开通管理需要将上面两个模型手动点击开通，提供了免费额度。

```bash
curl https://ark.cn-beijing.volces.com/api/v3/embeddings \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_API_KEY" \
  -d $'{
    "encoding_format": "float",
    "input": [
        " 天很蓝",
        "海很深"
    ],
    "model": "doubao-embedding-large-text-250515"
}'
```

如果有向量类型的字段，需要先定义 mappings。

创建一个 index：`http://127.0.0.1:9200/eino_example`

```bash
PUT http://127.0.0.1:9200/eino_example?pretty
{
    "mappings" : {
        "properties": {
            "content": { "type": "text" },
            "location": { "type": "text" },
            "content_vector": {
                "type": "dense_vector",
                "dims": 2048,
                "index": true,
                "similarity": "cosine"
            }
        }
    }
}


GET http://127.0.0.1:9200/eino_example
```



elasticsearch中的dense_vector类型，在版本**8.0 – 8.11**中，默认的最高维度是2048，在 **8.12+**之后是4096，当然，这个值越高计算越慢。数据的维度必须小于es能存储的维度，否则会报错。



配置文件 env.bat

```bash
set ARK_BASE_URL=https://ark.cn-beijing.volces.com/api/v3
set ARK_API_KEY=xxxx
set ARK_EMBEDDING_MODEL=doubao-embedding-large-text-250515
set ARK_CHAT_MODEL=deepseek-v3-2-251201
set ES_USERNAME=
set ES_PASSWORD=
set ES_HTTP_CA_CERT_PATH=
```

在当前命令行运行一下bat。



#### 写入数据

在 vscode 中使用 eino dev 工具编排graph

ctrl+shift+p 打开控制面板，输入 eino

![image-20251215141002562](D:\dev\php\magook\trunk\server\md\img\image-20251215141002562.png)

生成的代码保存到 indexeres8 目录下，然后需要修改一下，主要是一些配置项

```go
// embedding.go
func newEmbedding(ctx context.Context) (eb embedding.Embedder, err error) {
	// TODO Modify component configuration here.
	config := &ark.EmbeddingConfig{
		BaseURL: os.Getenv("ARK_BASE_URL"),
		APIKey:  os.Getenv("ARK_API_KEY"),
		Model:   os.Getenv("ARK_EMBEDDING_MODEL"),
	}
	eb, err = ark.NewEmbedder(ctx, config)
	if err != nil {
		return nil, err
	}
	
	return eb, nil
}

// indexer.go
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

// transformer.go
func newDocumentTransformer(ctx context.Context) (tfr document.Transformer, err error) {
	// TODO Modify component configuration here.
	config := &markdown.HeaderConfig{
		Headers: map[string]string{
			"#": "title"}}
	tfr, err = markdown.NewHeaderSplitter(ctx, config)
	if err != nil {
		return nil, err
	}
	return tfr, nil
}
```

`main.go`

```go
package main

import (
	"context"
	"fmt"
	"io/fs"
	"learn-eino/indexer_es8/indexeres8"
	"path/filepath"
	"strings"

	"github.com/cloudwego/eino/components/document"
)

func main() {
	ctx := context.Background()

	err := indexMarkdownFiles(ctx, "./eino-docs")
	if err != nil {
		panic(err)
	}

	fmt.Println("index success")
}

func indexMarkdownFiles(ctx context.Context, dir string) error {
	runner, err := indexeres8.Buildes8Indexer(ctx)
	if err != nil {
		return fmt.Errorf("build index graph failed: %w", err)
	}

	// 遍历 dir 下的所有 markdown 文件
	err = filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return fmt.Errorf("walk dir failed: %w", err)
		}
		if d.IsDir() {
			return nil
		}

		if !strings.HasSuffix(path, ".md") {
			fmt.Printf("[skip] not a markdown file: %s\n", path)
			return nil
		}

		fmt.Printf("[start] indexing file: %s\n", path)

		ids, err := runner.Invoke(ctx, document.Source{URI: path})
		if err != nil {
			return fmt.Errorf("invoke index graph failed: %w", err)
		}

		fmt.Printf("[done] indexing file: %s, len of parts: %d\n", path, len(ids))

		return nil
	})

	return err
}
```

```bash
> go run main.go
[start] indexing file: eino-docs\_index.md
[done] indexing file: eino-docs\_index.md, len of parts: 4
[start] indexing file: eino-docs\agent_llm_with_tools.md
[done] indexing file: eino-docs\agent_llm_with_tools.md, len of parts: 1
index success
```

以一级标题为分割，得到了五个 parts

![image-20251215142842335](D:\dev\php\magook\trunk\server\md\img\image-20251215142842335.png)



#### 检索数据

```bash
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
```



出现报错：

```bash
# github.com/cloudwego/eino-ext/components/retriever/es8/search_mode
..\..\golang\path\pkg\mod\github.com\cloudwego\eino-ext\components\retriever\es8@v0.0.0-20251212100737-81e5663e756e\search_mode\dense_vector_similarity.go:81:25: cannot use &types.
Query{…} (value of type *"github.com/elastic/go-elasticsearch/v8/typedapi/types".Query) as "github.com/elastic/go-elasticsearch/v8/typedapi/types".Query value in assignment
```

打开 dense_vector_similarity.go:81行

![image-20251216092148914](D:\dev\php\magook\trunk\server\md\img\image-20251216092148914.png)

在`eino-ext\components\retriever\es8`的`go.mod`中定义的依赖是 `github.com/elastic/go-elasticsearch/v8 v8.16.0`

而我使用的是`github.com/elastic/go-elasticsearch/v8 v8.19.1`，降低版本就行了。



### 四、Eino 智能助手

文档：[Eino 智能助手](https://www.cloudwego.io/zh/docs/eino/overview/bytedance_eino_practice/)



eino dev 可以将编排好的 graph 导出为 json schema，也可以导入 json schema 来创建graph。

于是，可以将 `eino-examples/quickstart/eino_assistant/eino/eino_agent.json`导入进来，但是这里面还是有问题，需要自行完善，修改完后导出到`eino_agent_es8.json`

![image-20251218112802680](D:\dev\php\magook\trunk\server\md\img\image-20251218112802680.png)

创建一个`.env`文件，就是上面的配置内容

```bash
> cd learn-eino\eino_agent\cmd\einoagent
> go build .
> einoagent.exe
```

访问：http://127.0.0.1:8080/agent/



会出来一条会话历史记录，这是来自 `data/memory`

然后提问`eino能做什么`

![image-20251218114424322](D:\dev\php\magook\trunk\server\md\img\image-20251218114424322.png)



右下角是日志，同时在 log 目录下也能看到。



**一些说明**

template role： system, user, tool, assistant

system message 在这里是指令，告诉智能体大致要做什么

assistant message 由agent输出，说明调用哪些tool

tool message 由agent输出，说明调用了哪个工具，参数和结果是什么



`schema.MessagesPlaceholder(key string, optional bool)`，可用于把一个 `[]*schema.Message` 插入到 message 列表中，常用于插入历史对话。optional 为 true 表示如果在 Variables 中没有这个字段，会填充空数组`[]`，如果为 false ，而 Variables 中没有这个字段，就会报错。

在 eino 中，使用 Variables 中的字段的时候，如果不存在这个key，那么是会报错的。



`FormatType: schema.FString`是使用`{variable}` 语法进行变量替换，简单直观，适合基础文本替换场景。示例：`"你是一个{role}，请帮我{task}。"`



`Variables`是由用户来维护的，在源码中就是随处可见的`vs map[string]any`，比如在 `lambda`中，在 `chatTemplate`中定义的 key

```json
{"Variables":
 {
     "content":"eino是什么",
     "date":"2025-12-17 16:04:50",
     "history":[
         {"role":"user","content":"Eino 好不好"},
         {"role":"assistant","content":"Eino是一个有其自身优势的框架。xxxx。","response_meta":{"finish_reason":"stop","usage":{"prompt_tokens":923,"prompt_token_details":{"cached_tokens":0},"completion_tokens":745,"total_tokens":1668,"completion_token_details":{}}}}],
     "retriever_result":"# Eino  是什么\n\n\u003e 💡\n\u003e Go AI 集成组件的研发框架。"
 },
 "Templates":null,
 "Extra":null
}
```



`DocumentListToMapLambda`节点的作用是数据格式转换，因为 es8 的输出是`[]*schema.Document`，而 chatTemplate 的输入是 `map[string]any`，这会导致报错。

```bash
Error running agent: failed to build agent graph: graph edge[Retriever1]-[ChatTemplate]: start node's output type[[]*schema.Document
] and end node's input type[map[string]interface {}] mismatch
```

我们将其转换到`Variables`中的`retriever_result`。这样就可以将查询到的内容放到 prompt 中去，使用`{retriever_result}`

```go
// 将 []*schema.Document 转换成 map[string]any
func documentListToMapLambda(ctx context.Context, input []*schema.Document, opts ...any) (output map[string]any, err error) {
	var contents []string
	for _, doc := range input {
		contents = append(contents, doc.Content)
	}
	contextText := strings.Join(contents, "\n\n")

	return map[string]interface{}{
		"retriever_result": contextText,
	}, nil
}
```



关于 tools_node，这里并没有调用到，因为 es8 已经查询到内容了，关于 tool 的使用说明：https://www.cloudwego.io/zh/docs/eino/core_modules/components/tools_node_guide/how_to_create_a_tool/

tool 的返回值是字符串，但最好定义为 json string，语义清晰，方便大模型理解

```go
return `{"status": "success", "result": "tool1 result"}`, nil
```



`React -- Tools -- ChatModel` 的执行流程 https://www.cloudwego.io/zh/docs/eino/core_modules/flow_integration_components/react_agent_manual/

![image-20251218144359711](D:\dev\php\magook\trunk\server\md\img\image-20251218144359711.png)

`React`调用 Tool ，得到结果后，将其追加到 ChatModel 的 input message，其 role 为 tool，这样 ChatModel 就会综合这些 messages 来做应答。



























