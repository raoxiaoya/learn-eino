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



注意：elasticsearch中的dense_vector类型，在版本**8.0 – 8.11**中，默认的最高维度是2048，在 **8.12+**之后是4096，当然，这个值越高计算越慢。数据的维度必须小于es能存储的维度，否则会报错。



配置文件`.env`

```bash
ARK_BASE_URL=https://ark.cn-beijing.volces.com/api/v3
ARK_API_KEY=xxxx
ARK_EMBEDDING_MODEL=doubao-embedding-large-text-250515
ARK_CHAT_MODEL=deepseek-v3-2-251201
```



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

**template role**： system, user, tool, assistant

system message 在这里是指令，告诉智能体大致要做什么，这个一定要写清楚。

assistant message 由agent输出，就是助手与你的交互。

tool message 由agent输出，说明调用了哪个工具，参数和结果是什么。



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



界面右上角的 Task Manager 是一个任务管理系统，任务的增删改查，同时提供了API。然后在 `util/tool/task`目录下有实现逻辑。因此你可以将 task tool 编排进 eino_agent，以此让 eino_agent 可以管理任务系统。



### 工具调用

将 task tool 添加到 flow 中

```go
// flow.go

taskTool, err := task.NewTaskTool(ctx, nil)
config.ToolsConfig.Tools = []tool.BaseTool{toolIns21, toolIns22, toolIns23, toolIns24, toolIns25, taskTool}
```

提示词

```bash
帮我在task manager系统创建一个任务：
标题：eino_agent任务标题；
内容：eino_agent任务内容；
截止时间：2025-12-25T15:15
```

日志

```json
{
    "role":"assistant",
    "content":"我来帮你创建一个任务。首先，我需要将你提供的截止时间格式转换为任务管理器需要的格式。\n\n",
    "tool_calls":[
        {
            "index":0,
            "id":"call_akjdswsyvkwv47khha8dnp1p",
            "type":"function",
            "function":{
                "name":"task_manager",
                "arguments":"{\"action\": \"add\", \"task\": {\"id\": \"task_\" + Math.random().toString(36).substr(2, 9), \"title\": \"eino_agent任务标题\", \"content\": \"eino_agent任务标题\", \"completed\": false, \"deadline\": \"2025-12-25T15:15:00Z\", \"created_at\": new Date().toISOString()}, \"list\": {\"query\": \"\", \"is_done\": false, \"limit\": 10}}"
            }
        }
]
```

但是并没有创建任务，也没有调用对应的 invoke 方法。

报错

```bash
2025/12/24 10:49:56 error publishing log: write tcp 127.0.0.1:8080->127.0.0.1:52071: wsasend: An established connection was aborted by the software in your host machine.
```



```go
type InvokableTool interface {
    BaseTool

    // InvokableRun call function with arguments in JSON format
    InvokableRun(ctx context.Context, argumentsInJSON string, opts ...Option) (string, error)
}

```

写一个测试文件

```go
func main() {
	ctx := context.Background()
	taskTool, err := task.NewTaskTool(ctx, nil)
	if err != nil {
		panic(err)
	}

	str, _ := json.Marshal(task.TaskRequest{
		Action: task.ActionAdd,
		Task: &task.Task{
			Title:    "test task",
			Content:  "test content",
			Deadline: "2025-12-30 12:12:12",
		},
	})

	tt, ok := taskTool.(tool.InvokableTool)
	if !ok {
		panic("taskTool is not InvokableTool")
	}
	res, err := tt.InvokableRun(ctx, string(str))
	if err != nil {
		panic(err)
	}
	println(res)
}
```

```bash
{"status":"success","task_list":[{"id":"905bbad6-690a-4e98-85e6-41d9232e6f9d","title":"test task","content":"test content","completed":false,"deadline":"202
5-12-30 12:12:12","is_deleted":false,"created_at":"2025-12-25T14:48:06+08:00"}],"error":""}
```

说明它是能正常写入的

```json
{"id":"f0a28183-ef37-416a-87c4-5e957fbe969d","title":"test task","content":"test content","completed":false,"deadline":"2025-12-30 12:12:12","is_deleted":false,"created_at":"2025-12-23T18:16:33+08:00"}
```



### 关于 tools call

tools 是由chatModel调用的，因此需要将 tools 节点放在 chatModel 的下面，并且将 tools info 注入到 chatModel 

```bash
chatModel.BindForcedTools([]*schema.ToolInfo{info})
或者
chatModel.BindTools([]*schema.ToolInfo{info})
或者
chatModel.WithTools([]*schema.ToolInfo{info})
```

可以使用 react 节点作为 tools call 调度中心。

其实 react 内部也是构建一个 graph（chatTemplate+chatModel+toolsNode）。

以下是 react 示例

```go
func Main1() {
	ctx := context.Background()

	chatModel, err := util.GetChatModel(ctx)
	if err != nil {
		fmt.Printf("Get ChatModel failed, err=%v\n", err)
		return
	}

	taskTool, err := task.NewTaskTool(ctx, nil)
	if err != nil {
		fmt.Printf("NewTaskTool failed, err=%v\n", err)
		return
	}

	agent, err := react.NewAgent(ctx, &react.AgentConfig{
		ToolCallingModel: chatModel,
		ToolsConfig: compose.ToolsNodeConfig{
			Tools: []tool.BaseTool{taskTool},
		},
		// 直接将 tool 的结果返回，而不用再回到 chatModel 进行处理
		ToolReturnDirectly: map[string]struct{}{"task_manager": {}},
		MessageModifier: func(ctx context.Context, input []*schema.Message) []*schema.Message {
			res := make([]*schema.Message, 0, len(input)+1)
			res = append(res, schema.SystemMessage("你是一个后台操作员，使用 task_manager API，对task进行操作"))
			res = append(res, input...)
			return res
		},
	})

	response, err := agent.Generate(ctx, []*schema.Message{schema.UserMessage(`帮我在task manager系统创建一个任务：
标题：eino_agent任务标题；
内容：eino_agent任务内容；
截止时间：2025-12-25T15:15；
任务ID：eac76c5a-629a-45d2-8c4e-0ac769a470f0；
任务状态：未完成；
创建时间：2026-01-07T10:03:07+08:00
`)})

	if err != nil {
		fmt.Printf("Generate failed, err=%v\n", err)
		return
	}
	fmt.Println(response.Role, response.Content)
}
```

或者将 react 包装成 graph/chain 的一个节点，然后 add lambda

```go
lba, err = compose.AnyLambda(ins.Generate, ins.Stream, nil, nil)
```

在 graph 中，也可以直接将 tools node 直接添加在 chatModel 下游。

```go
toolsNode, err := compose.NewToolNode(ctx, &compose.ToolsNodeConfig{
    Tools: []tool.BaseTool{taskTool},
})

g := compose.NewGraph[map[string]any, []*schema.Message]()
_ = g.AddChatTemplateNode(nodeKeyOfTemplate, chatTpl)
_ = g.AddChatModelNode(nodeKeyOfChatModel, chatModel)
_ = g.AddToolsNode(nodeKeyOfTools, toolsNode)
_ = g.AddEdge(compose.START, nodeKeyOfTemplate)
_ = g.AddEdge(nodeKeyOfTemplate, nodeKeyOfChatModel)
_ = g.AddEdge(nodeKeyOfChatModel, nodeKeyOfTools)
_ = g.AddEdge(nodeKeyOfTools, compose.END)
```

往往 tool node 作为最后一个节点。

chatModel 是与用户交互的节点，它来决定调用哪个工具或者什么工具都不调用，tool 的返回信息会回到 chatModel ，最后由 chatModel 加工后响应给用户。或者设置 ToolReturnDirectly，这样 agent 会直接将 tool 的结果返回，不会再调用chatModel 。不管怎么返回，其输出类型依然是`*schema.Message`，只是 role 分别是`assistant`和`tool`。

```bash
assistant 任务已成功创建！以下是创建的任务详情：

- **任务ID**: c19aa787-31c4-4b18-92c5-26a3525ece44（系统自动生成了新的ID）
- **标题**: eino_agent任务标题
- **内容**: eino_agent任务标题
- **状态**: 未完成
- **截止时间**: 2025-12-25T15:15
- **创建时间**: 2026-01-07T10:03:07+08:00（系统自动生成的当前时间）

请注意，系统自动为任务生成了新的ID（c19aa787-31c4-4b18-92c5-26a3525ece44），而不是使用您提供的eac76c5a-629a-45d2-8c4e-0ac769a470f0。创建时间也由系
统自动设置为当前时间。
```

```bash
tool {"status":"success","task_list":[{"id":"04dca8bd-552d-4559-8acf-faebeff5d898","title":"eino_agent任务标题","content":"eino_agent任务标题","co
mpleted":false,"deadline":"2025-12-25T15:15","is_deleted":false,"created_at":"2026-01-07T10:13:41+08:00"}],"error":""}
```



但是这里面也是有不确定性，比如 chatModel 能否正确的理解你的意思，chatModel 能正确的构建参数等等。



### react  与 stream

在使用 react 的时候发现一个奇怪的现象，使用 Generate 运行 react 是好的，能够正常调用 tool 并响应信息。在使用 Stream 运行 react 的时候，在 chatModel 第一次输出之后，由于 EOF 导致程序退出，于是无法调用 tool。按理说，chatModel 应该会多次输出，第一次告诉你，它将要做什么，然后去调用tool，然后合并结果再次输出，如果涉及到多个 tool，那么 chatModel 就会有更多的输出。也就是说，这个stream是不应该中途中断的。在 github 上找到了对应的 [issue](https://github.com/cloudwego/eino/issues/613)

这个问题跟模型有关系，解决方案为 [StreamToolCallChecker](https://www.cloudwego.io/zh/docs/eino/core_modules/flow_integration_components/react_agent_manual/#streamtoolcallchecker)

agent 是否去调用 tool ，是根据 chatModel 的返回信息，如果是在 stream 上，则需要每个 chunk 都要检查，但这样做效率低下，因此 eino 默认的 toolCallChecker 为 firstChunkStreamToolCallChecker，也就是只检查第一个 chunk。

```go
func firstChunkStreamToolCallChecker(_ context.Context, sr *schema.StreamReader[*schema.Message]) (bool, error) {
	defer sr.Close()

	for {
		msg, err := sr.Recv()
		if err == io.EOF {
			return false, nil
		}
		if err != nil {
			return false, err
		}

		if len(msg.ToolCalls) > 0 {
			return true, nil
		}

		if len(msg.Content) == 0 { // skip empty chunks at the front
			continue
		}

		return false, nil
	}
}
```

此方法的问题在于，如果在输出 tool call 之前先输出了别的内容，此方法会返回 false，意味着 agent 不会，然而这并不准确。我用的 deepseek 模型会先告诉你它要干什么，后面才会返回 tool 信息，因此需要改进一下

```go
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
```

这样就不会输出chatModel在中途返回的信息，只要最终的结果。

![image-20260108165858614](D:\dev\php\magook\trunk\server\md\img\image-20260108165858614.png)



### 性能观测

##### APMPlus

如果在运行时，在 .env 文件中指定了 `APMPLUS_APP_KEY`，便可在 [火山引擎 APMPlus](https://console.volcengine.com/apmplus-server) 平台中，登录对应的账号，查看 Trace 以及 Metrics 详情。

点服务端监控，开通APMPLUS服务。

点击[接入中心](https://console.volcengine.com/observe/apmplus-server/region:apmplus-server+cn-beijing/server/access/service)，复制 APP Key

Eino接入教程：https://www.volcengine.com/docs/6431/1471121?lang=zh



```go
// 创建apmplus handler
cbh, showdown, err := apmplus.NewApmplusHandler(&apmplus.Config{
    Host: "apmplus-cn-beijing.volces.com:4317",
    AppKey:      os.Getenv("APMPLUS_APP_KEY"),
    ServiceName: "eino-app",
    Release:     "release/v0.0.1",
})
if err != nil {
    log.Fatal(err)
}

// 设置apmplus为全局callback
callbacks.AppendGlobalHandlers(cbh)

// 运行 agent 

// 最后，等待所有trace和metrics上报完成后退出
showdown(ctx)
```

以官方的查询菜品的DEMO为例：[eino-examples/flow/agent/react](https://github.com/cloudwego/eino-examples/tree/main/flow/agent/react)

![image-20260109104705394](D:\dev\php\magook\trunk\server\md\img\image-20260109104705394.png)

![image-20260109105349574](D:\dev\php\magook\trunk\server\md\img\image-20260109105349574.png)



##### Langfuse

如果在运行时，在 .env 文件中指定了 `LANGFUSE_PUBLIC_KEY` 和 `LANGFUSE_SECRET_KEY`，便可在 Langfuse 平台中，登录对应的账号，查看请求的 Trace 详情。



##### React流程分析

![image-20260112093558744](D:\dev\php\magook\trunk\server\md\img\image-20260112093558744.png)

React 

```bash
input

system：你是一个帮助用户推荐餐厅和菜品的助手，根据用户的需要，查询餐厅信息并推荐，查询餐厅的菜品并推荐。
user：我在北京，给我推荐一些菜，需要有口味辣一点的菜，至少推荐有 2 家餐厅
```

```bash
output

assistant：
根据查询结果，我为您推荐以下北京餐厅的辣味菜品：

## 推荐餐厅及辣味菜品：

### 1. 云边小馆（评分：3分）
**推荐辣味菜品：**
- **韩式辣白菜**（20元）- 这可是开过光的辣白菜，好吃得很，评分9分
- **酸辣土豆丝**（10元）- 酸酸辣辣的土豆丝，评分9分  
- **酸辣粉**（5元）- 酸酸辣辣的粉，价格实惠

### 2. 聚福轩食府（评分：5分）
**推荐辣味菜品：**
- **火辣辣的吻**（60元）- 凉拌猪嘴，口味辣而不腻，评分9分
- **辣椒拌皮蛋**（15元）- 擂椒皮蛋，下饭的神器，评分8分
- **大刀回锅肉**（40元）- 经典的回锅肉，肉很大，评分8分（微辣）

## 推荐总结：
两家餐厅都有不错的辣味菜品选择：
- **云边小馆**的辣味菜品更偏向韩式和酸辣口味，价格相对实惠
- **聚福轩食府**的辣味菜品更传统，特别是"火辣辣的吻"这道菜评分很高，辣而不腻

如果您喜欢韩式辣味和酸辣口味，推荐云边小馆；如果您喜欢传统中式辣味和凉拌菜，推荐聚福轩食府。两家餐厅都有多种辣味菜品可供选择。
```

ChatModel 

```bash
input

system：你是一个帮助用户推荐餐厅和菜品的助手，根据用户的需要，查询餐厅信息并推荐，查询餐厅的菜品并推荐。
user：我在北京，给我推荐一些菜，需要有口味辣一点的菜，至少推荐有 2 家餐厅
```

```bash
output

assistant：
我来帮您在北京推荐一些口味辣一点的菜品和餐厅。首先让我查询一下北京的餐厅信息。
```

tools

```bash
input：
{"role":"assistant","content":"我来帮您在北京推荐一些口味辣一点的菜品和餐厅。首先让我查询一下北京的餐厅信息。\n\n","tool_calls":[{"index":0,"id":"call_9zao24vvrqpwm6qn2c575ei9","type":"function","function":{"name":"query_restaurants","arguments":"{\"location\": \"北京\", \"topn\": 10}"}}],"response_meta":{"finish_reason":"tool_calls","usage":{"prompt_tokens":464,"prompt_token_details":{"cached_tokens":0},"completion_tokens":82,"total_tokens":546,"completion_token_details":{}}},"extra":{"ark-model-name":"deepseek-v3-2-251201","ark-service-tier":"default","ark-request-id":"02176792534812543cc1cf6ced7e437334ae832707c2826e6512f"}}

output：
[[{"role":"tool","content":"[{\"id\":\"1001\",\"name\":\"云边小馆\",\"place\":\"北京\",\"desc\":\"\",\"score\":3},{\"id\":\"1002\",\"name\":\"聚福轩食府\",\"place\":\"北京\",\"desc\":\"\",\"score\":5},{\"id\":\"1003\",\"name\":\"花影食舍\",\"place\":\"上海\",\"desc\":\"\",\"score\":10}]","tool_call_id":"call_9zao24vvrqpwm6qn2c575ei9","tool_name":"query_restaurants"}]]
```

query_restaurants 

```bash
input：
"{\"location\": \"北京\", \"topn\": 10}"

output：
"[{\"id\":\"1001\",\"name\":\"云边小馆\",\"place\":\"北京\",\"desc\":\"\",\"score\":3},{\"id\":\"1002\",\"name\":\"聚福轩食府\",\"place\":\"北京\",\"desc\":\"\",\"score\":5},{\"id\":\"1003\",\"name\":\"花影食舍\",\"place\":\"上海\",\"desc\":\"\",\"score\":10}]"
```

ChatModel 

```bash
input

system：你是一个帮助用户推荐餐厅和菜品的助手，根据用户的需要，查询餐厅信息并推荐，查询餐厅的菜品并推荐。
user：我在北京，给我推荐一些菜，需要有口味辣一点的菜，至少推荐有 2 家餐厅
assistant：我来帮您在北京推荐一些口味辣一点的菜品和餐厅。首先让我查询一下北京的餐厅信息。
tool：[{"id":"1001","name":"云边小馆","place":"北京","desc":"","score":3},{"id":"1002","name":"聚福轩食府","place":"北京","desc":"","score":5},{"id":"1003","name":"花影食舍","place":"上海","desc":"","score":10}]
```

```bash
output

assistant：我看到北京有两家餐厅。让我分别查询这两家餐厅的菜品，看看哪些是辣味的：
```

tools

```bash
input：
{"role":"assistant","content":"我看到北京有两家餐厅。让我分别查询这两家餐厅的菜品，看看哪些是辣味的：\n\n\n\n","tool_calls":[{"index":0,"id":"call_akfts8ntgtxq62h30i3kpce6","type":"function","function":{"name":"query_dishes","arguments":"{\"restaurant_id\": \"1001\", \"topn\": 10}"}}],"response_meta":{"finish_reason":"tool_calls","usage":{"prompt_tokens":635,"prompt_token_details":{"cached_tokens":0},"completion_tokens":85,"total_tokens":720,"completion_token_details":{}}},"extra":{"ark-service-tier":"default","ark-request-id":"02176792535212643cc1cf6ced7e437334ae832707c28266cb682","ark-model-name":"deepseek-v3-2-251201"}}

output：
[[{"role":"tool","content":"[{\"name\":\"红烧肉\",\"desc\":\"一块红烧肉\",\"price\":20,\"score\":8},{\"name\":\"清泉牛肉\",\"desc\":\"很多的水煮牛肉\",\"price\":50,\"score\":8},{\"name\":\"清炒小南瓜\",\"desc\":\"炒的糊糊的南瓜\",\"price\":5,\"score\":5},{\"name\":\"韩式辣白菜\",\"desc\":\"这可是开过光的辣白菜，好吃得很\",\"price\":20,\"score\":9},{\"name\":\"酸辣土豆丝\",\"desc\":\"酸酸辣辣的土豆丝\",\"price\":10,\"score\":9},{\"name\":\"酸辣粉\",\"desc\":\"酸酸辣辣的粉\",\"price\":5,\"score\":0}]","tool_call_id":"call_akfts8ntgtxq62h30i3kpce6","tool_name":"query_dishes"}]]
```

query_dishes 

```bash
input：
"{\"restaurant_id\": \"1001\", \"topn\": 10}"

output：
"[{\"name\":\"红烧肉\",\"desc\":\"一块红烧肉\",\"price\":20,\"score\":8},{\"name\":\"清泉牛肉\",\"desc\":\"很多的水煮牛肉\",\"price\":50,\"score\":8},{\"name\":\"清炒小南瓜\",\"desc\":\"炒的糊糊的南瓜\",\"price\":5,\"score\":5},{\"name\":\"韩式辣白菜\",\"desc\":\"这可是开过光的辣白菜，好吃得很\",\"price\":20,\"score\":9},{\"name\":\"酸辣土豆丝\",\"desc\":\"酸酸辣辣的土豆丝\",\"price\":10,\"score\":9},{\"name\":\"酸辣粉\",\"desc\":\"酸酸辣辣的粉\",\"price\":5,\"score\":0}]"
```

ChatModel 

```bash
input

system：你是一个帮助用户推荐餐厅和菜品的助手，根据用户的需要，查询餐厅信息并推荐，查询餐厅的菜品并推荐。
user：我在北京，给我推荐一些菜，需要有口味辣一点的菜，至少推荐有 2 家餐厅
assistant：我来帮您在北京推荐一些口味辣一点的菜品和餐厅。首先让我查询一下北京的餐厅信息。
tool：[{"id":"1001","name":"云边小馆","place":"北京","desc":"","score":3},{"id":"1002","name":"聚福轩食府","place":"北京","desc":"","score":5},{"id":"1003","name":"花影食舍","place":"上海","desc":"","score":10}]
assistant：我看到北京有两家餐厅。让我分别查询这两家餐厅的菜品，看看哪些是辣味的：
tool：[{"name":"红烧肉","desc":"一块红烧肉","price":20,"score":8},{"name":"清泉牛肉","desc":"很多的水煮牛肉","price":50,"score":8},{"name":"清炒小南瓜","desc":"炒的糊糊的南瓜","price":5,"score":5},{"name":"韩式辣白菜","desc":"这可是开过光的辣白菜，好吃得很","price":20,"score":9},{"name":"酸辣土豆丝","desc":"酸酸辣辣的土豆丝","price":10,"score":9},{"name":"酸辣粉","desc":"酸酸辣辣的粉","price":5,"score":0}]
```

```bash
output

assistant：现在查询第二家餐厅的菜品：
```

tools

```bash
input

{"role":"assistant","content":"现在查询第二家餐厅的菜品：\n\n","tool_calls":[{"index":0,"id":"call_l7k9djvc5ek8ztegmlwafz1q","type":"function","function":{"name":"query_dishes","arguments":"{\"restaurant_id\": \"1002\", \"topn\": 10}"}}],"response_meta":{"finish_reason":"tool_calls","usage":{"prompt_tokens":875,"prompt_token_details":{"cached_tokens":0},"completion_tokens":71,"total_tokens":946,"completion_token_details":{}}},"extra":{"ark-request-id":"02176792535621643cc1cf6ced7e437334ae832707c28268bf259","ark-model-name":"deepseek-v3-2-251201","ark-service-tier":"default"}}
```

```bash
output

[[{"role":"tool","content":"[{\"name\":\"红烧排骨\",\"desc\":\"一块一块的排骨\",\"price\":43,\"score\":7},{\"name\":\"大刀回锅肉\",\"desc\":\"经典的回锅肉, 肉很大\",\"price\":40,\"score\":8},{\"name\":\"火辣辣的吻\",\"desc\":\"凉拌猪嘴，口味辣而不腻\",\"price\":60,\"score\":9},{\"name\":\"辣椒拌皮蛋\",\"desc\":\"擂椒皮蛋，下饭的神器\",\"price\":15,\"score\":8}]","tool_call_id":"call_l7k9djvc5ek8ztegmlwafz1q","tool_name":"query_dishes"}]]
```

query_dishes 

```bash
input

"{\"restaurant_id\": \"1002\", \"topn\": 10}"
```

```bash
output

"[{\"name\":\"红烧排骨\",\"desc\":\"一块一块的排骨\",\"price\":43,\"score\":7},{\"name\":\"大刀回锅肉\",\"desc\":\"经典的回锅肉, 肉很大\",\"price\":40,\"score\":8},{\"name\":\"火辣辣的吻\",\"desc\":\"凉拌猪嘴，口味辣而不腻\",\"price\":60,\"score\":9},{\"name\":\"辣椒拌皮蛋\",\"desc\":\"擂椒皮蛋，下饭的神器\",\"price\":15,\"score\":8}]"
```

ChatModel 

```bash
input

system：你是一个帮助用户推荐餐厅和菜品的助手，根据用户的需要，查询餐厅信息并推荐，查询餐厅的菜品并推荐。
user：我在北京，给我推荐一些菜，需要有口味辣一点的菜，至少推荐有 2 家餐厅
assistant：我来帮您在北京推荐一些口味辣一点的菜品和餐厅。首先让我查询一下北京的餐厅信息。
tool：[{"id":"1001","name":"云边小馆","place":"北京","desc":"","score":3},{"id":"1002","name":"聚福轩食府","place":"北京","desc":"","score":5},{"id":"1003","name":"花影食舍","place":"上海","desc":"","score":10}]
assistant：我看到北京有两家餐厅。让我分别查询这两家餐厅的菜品，看看哪些是辣味的：
tool：[{"name":"红烧肉","desc":"一块红烧肉","price":20,"score":8},{"name":"清泉牛肉","desc":"很多的水煮牛肉","price":50,"score":8},{"name":"清炒小南瓜","desc":"炒的糊糊的南瓜","price":5,"score":5},{"name":"韩式辣白菜","desc":"这可是开过光的辣白菜，好吃得很","price":20,"score":9},{"name":"酸辣土豆丝","desc":"酸酸辣辣的土豆丝","price":10,"score":9},{"name":"酸辣粉","desc":"酸酸辣辣的粉","price":5,"score":0}]
assistant：现在查询第二家餐厅的菜品
tool：[{"name":"红烧排骨","desc":"一块一块的排骨","price":43,"score":7},{"name":"大刀回锅肉","desc":"经典的回锅肉, 肉很大","price":40,"score":8},{"name":"火辣辣的吻","desc":"凉拌猪嘴，口味辣而不腻","price":60,"score":9},{"name":"辣椒拌皮蛋","desc":"擂椒皮蛋，下饭的神器","price":15,"score":8}]
```

```bash
output

根据查询结果，我为您推荐以下北京餐厅的辣味菜品：

## 推荐餐厅及辣味菜品：

### 1. 云边小馆（评分：3分）
**推荐辣味菜品：**
- **韩式辣白菜**（20元）- 这可是开过光的辣白菜，好吃得很，评分9分
- **酸辣土豆丝**（10元）- 酸酸辣辣的土豆丝，评分9分  
- **酸辣粉**（5元）- 酸酸辣辣的粉，价格实惠

### 2. 聚福轩食府（评分：5分）
**推荐辣味菜品：**
- **火辣辣的吻**（60元）- 凉拌猪嘴，口味辣而不腻，评分9分
- **辣椒拌皮蛋**（15元）- 擂椒皮蛋，下饭的神器，评分8分
- **大刀回锅肉**（40元）- 经典的回锅肉，肉很大，评分8分（微辣）

## 推荐总结：
两家餐厅都有不错的辣味菜品选择：
- **云边小馆**的辣味菜品更偏向韩式和酸辣口味，价格相对实惠
- **聚福轩食府**的辣味菜品更传统，特别是"火辣辣的吻"这道菜评分很高，辣而不腻

如果您喜欢韩式辣味和酸辣口味，推荐云边小馆；如果您喜欢传统中式辣味和凉拌菜，推荐聚福轩食府。两家餐厅都有多种辣味菜品可供选择。
```





### Callback调试

`RunInfo`表示自己是谁，也就是当前节点信息

`callbacks.Handler`

```go
type Handler interface {
    OnStart(ctx context.Context, info *RunInfo, input CallbackInput) context.Context
    
    OnEnd(ctx context.Context, info *RunInfo, output CallbackOutput) context.Context

    OnError(ctx context.Context, info *RunInfo, err error) context.Context

    OnStartWithStreamInput(ctx context.Context, info *RunInfo,
        input *schema.StreamReader[CallbackInput]) context.Context
    
    OnEndWithStreamOutput(ctx context.Context, info *RunInfo,
        output *schema.StreamReader[CallbackOutput]) context.Context
}
```



如果一个 Handler，不想关注所有的 5 个触发时机，只想关注一部分，比如只关注 OnStart，建议使用 `NewHandlerBuilder().OnStartFn(...).Build()`。如果不想关注所有的组件类型，只想关注特定组件，比如 ChatModel，建议使用 `NewHandlerHelper().ChatModel(...).Handler()`，可以只接收 ChatModel 的回调并拿到一个具体类型的 CallbackInput/CallbackOutput。

```go
handler := callbacks.NewHandlerBuilder().
		OnStartFn(
			func(ctx context.Context, info *callbacks.RunInfo, input callbacks.CallbackInput) context.Context {
				log.Printf("onStart, runInfo: %v, input: %v", info, input)
				return ctx
			}).
		OnEndFn(
			func(ctx context.Context, info *callbacks.RunInfo, output callbacks.CallbackOutput) context.Context {
				log.Printf("onEnd, runInfo: %v, out: %v", info, output)
				return ctx
			}).
		Build()
```



全局注入

```go
callbacks.AppendGlobalHandlers(handler)
```







