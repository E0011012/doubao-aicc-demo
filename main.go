package main

import (
	"context"
	"fmt"
	"io"

	"github.com/volcengine/volcengine-go-sdk/service/arkruntime"
	"github.com/volcengine/volcengine-go-sdk/service/arkruntime/model"
	"github.com/volcengine/volcengine-go-sdk/volcengine"
)

func main() {
	// 设置API Key（
	apiKey := "0fc8bfc8-f908-41e1-8f33-410d5c5ab499"
	// 创建客户端
	// 设置自定义baseurl
	customBaseUrl := "http://app.cxmt.com/doubao/aicc/api/v3"
	client := arkruntime.NewClientWithApiKey(apiKey, arkruntime.WithBaseUrl(customBaseUrl))
	ctx := context.Background()
	fmt.Println("-----secure streaming request -----")
	maxCompletionTokens := 1000
	req := model.CreateChatCompletionRequest{
		Model: "ep-20251104144425-jpvbx",
		Messages: []*model.ChatCompletionMessage{
			{
				Role: model.ChatMessageRoleSystem,
				Content: &model.ChatCompletionMessageContent{
					StringValue: volcengine.String("你是豆包，是由字节跳动开发的 AI 人工智能助手"),
				},
			},
			{
				Role: model.ChatMessageRoleUser,
				Content: &model.ChatCompletionMessageContent{
					StringValue: volcengine.String("请你简单做个自我介绍"),
				},
			},
		},
		//用量统计需要
		StreamOptions: &model.StreamOptions{
			IncludeUsage: true,
		},
		MaxCompletionTokens: &maxCompletionTokens,
		Thinking: &model.Thinking{
			Type: model.ThinkingTypeDisabled,
		},
	}
	// 设置自定义请求头
	customHeaders := map[string]string{
		"x-is-encrypted":         "false",
		"x-ark-moderation-scene": "aicc-skip",
	}
	res, err := client.CreateChatCompletion(ctx, req, arkruntime.WithCustomHeaders(customHeaders))
	fmt.Println(res)

	stream, err := client.CreateChatCompletionStream(ctx, req, arkruntime.WithCustomHeaders(customHeaders))
	if err != nil {
		fmt.Printf("stream chat error: %v\n", err)
		return
	}
	defer stream.Close()
	for {
		recv, err := stream.Recv()
		if err == io.EOF {
			return
		}
		if err != nil {
			fmt.Printf("Stream chat error: %v\n", err)
			return
		}
		if len(recv.Choices) > 0 {
			fmt.Print(recv.Choices[0].Delta.Content)
		}
	}
}
