/*
 * Copyright 2024 CloudWeGo Authors
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package react_agent

import (
	"context"
	"errors"
	"fmt"
	"io"
	"time"

	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/compose"
	"github.com/cloudwego/eino/flow/agent/react"
	"github.com/cloudwego/eino/schema"

	"learn-eino/react_agent/tools"
	"learn-eino/util"
)

func Main3() {
	ctx := context.Background()

	arkModel, err := util.GetChatModel(ctx)
	if err != nil {
		fmt.Printf("failed to create chat model: %v\n", err)
		return
	}

	restaurantTool := tools.GetRestaurantTool()
	dishTool := tools.GetDishTool()

	rAgent, err := react.NewAgent(ctx, &react.AgentConfig{
		ToolCallingModel: arkModel,
		ToolsConfig: compose.ToolsNodeConfig{
			Tools: []tool.BaseTool{restaurantTool, dishTool},
		},
		StreamToolCallChecker: util.CustomToolCallChecker,
	})
	if err != nil {
		fmt.Printf("failed to create agent: %v\n", err)
		return
	}

	sr, err := rAgent.Stream(ctx, []*schema.Message{
		{
			Role:    schema.System,
			Content: "你是一个帮助用户推荐餐厅和菜品的助手，根据用户的需要，查询餐厅信息并推荐，查询餐厅的菜品并推荐。",
		},
		{
			Role:    schema.User,
			Content: "我在北京，给我推荐一些菜，需要有口味辣一点的菜，至少推荐有 2 家餐厅",
		},
	})
	if err != nil {
		fmt.Printf("failed to stream: %v\n", err)
		return
	}

	defer sr.Close()

	fmt.Println("===== start streaming =====")

	for {
		msg, err := sr.Recv()
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			fmt.Printf("failed to recv: %v", err)
			return
		}

		// 打字机打印
		fmt.Printf("%v", msg.Content)
	}

	time.Sleep(3 * time.Second)

	fmt.Println("===== finished =====")
}
