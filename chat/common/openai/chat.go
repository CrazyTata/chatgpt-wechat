package openai

import (
	"errors"
	"fmt"
	"io"

	copenai "github.com/sashabaranov/go-openai"
	"github.com/zeromicro/go-zero/core/logx"
)

// ChatStream 数据流式传输
func (c *ChatClient) ChatStream(req []ChatModelMessage, channel chan string) (string, error) {
	first := 0
	for i, _ := range req {

		if i%2 == 0 {
			continue
		}
		//估算长度
		if NumTokensFromMessages(req[len(req)-i-1:], c.PostModel) < (3900 - c.MaxToken) {
			first = len(req) - i - 1
		} else {
			break
		}
	}
	var messages []copenai.ChatCompletionMessage

	if first != 0 {
		messages = c.DealMultiplePrompt(req)
	}

	for _, message := range req[first:] {
		messages = append(messages, copenai.ChatCompletionMessage{
			Role:    message.Role,
			Content: message.Content,
		})
	}
	if c.Model == "" || (c.Model != ChatModel && c.Model != ChatModelNew && c.Model != ChatModel4) {
		c.Model = ChatModel
	}
	request := copenai.ChatCompletionRequest{
		Model:       c.PostModel,
		Messages:    messages,
		MaxTokens:   c.MaxToken,
		Temperature: c.Temperature,
		TopP:        1,
	}
	var stream *copenai.ChatCompletionStream
	stream1, err1 := c.MakeOpenAILoopRequest(&OpenAIRequest{
		FuncName: "CreateChatCompletionStream",
		Request:  request,
	})
	if err1 != nil {
		fmt.Println(err1)
		return "", err1
	}
	stream2, ok := stream1.(*copenai.ChatCompletionStream)
	if !ok {
		return "", errors.New("conversion failed")
	}
	stream = stream2
	messageText := ""
	for {
		response, err := stream.Recv()
		if errors.Is(err, io.EOF) {
			logx.Info("Stream finished")
			return messageText, nil
		}

		if err != nil {
			fmt.Printf("Stream error: %v\n", err)
			close(channel)
			return messageText, err
		}

		if len(response.Choices) > 0 {
			content := response.Choices[0].Delta.Content
			channel <- content
			// 如果是对内容的进行补充
			messageText += content
			// 结算
			if response.Choices[0].FinishReason != "" {
				close(channel)
				return messageText, nil
			}
		}

		//logx.Info("Stream response:", response)
	}
}

func (c *ChatClient) Chat(req []ChatModelMessage) (string, error) {

	// 打印请求信息
	logx.Info("req: ", req)

	first := 0
	for i, _ := range req {

		if i%2 == 0 {
			continue
		}
		//估算长度
		if NumTokensFromMessages(req[len(req)-i-1:], c.PostModel) < (3900 - c.MaxToken) {
			first = len(req) - i - 1
		} else {
			break
		}
	}

	var messages []copenai.ChatCompletionMessage
	if first != 0 {
		messages = c.DealMultiplePrompt(req)
	}
	for _, message := range req[first:] {
		messages = append(messages, copenai.ChatCompletionMessage{
			Role:    message.Role,
			Content: message.Content,
		})
	}
	request := copenai.ChatCompletionRequest{
		Model:       c.PostModel,
		Messages:    messages,
		MaxTokens:   c.MaxToken,
		Temperature: c.Temperature,
		TopP:        1,
	}
	var chat copenai.ChatCompletionResponse
	chatOrigin, err1 := c.MakeOpenAILoopRequest(&OpenAIRequest{
		FuncName: "CreateChatCompletion",
		Request:  request,
	})
	if err1 != nil {
		return "", err1
	}
	chat1, ok := chatOrigin.(copenai.ChatCompletionResponse)
	if !ok {
		return "", errors.New("conversion failed")
	}
	chat = chat1

	txt := ""
	for _, choice := range chat.Choices {
		txt += choice.Message.Content
	}

	return txt, nil
}

func (c *ChatClient) DealMultiplePrompt(req []ChatModelMessage) (res []copenai.ChatCompletionMessage) {

	for _, msg := range req {
		if msg.Role == "system" {
			res = append(res, copenai.ChatCompletionMessage{
				Role:    msg.Role,
				Content: msg.Content,
			})
		}
	}

	return
}
