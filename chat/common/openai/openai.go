package openai

import (
	"chat/common/redis"
	context "context"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"

	copenai "github.com/sashabaranov/go-openai"
	"github.com/zeromicro/go-zero/core/logx"

	"golang.org/x/net/proxy"
)

const TextModel = "text-davinci-003"
const ChatModel = "gpt-3.5-turbo"
const ChatModelNew = "gpt-3.5-turbo-0301"
const ChatModel4 = "gpt-4"

const MaxToken = 2000
const Temperature = 0.8

const NeedLoopErrorMessage = "Incorrect API key provided"

type ChatModelMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ChatClient struct {
	APIKeys     []string `json:"api_keys"`
	APIKey      string   `json:"api_key"`
	HttpProxy   string   `json:"http_proxy"`
	Socks5Proxy string   `json:"socks5_proxy"`
	Model       string   `json:"model"`
	BaseHost    string   `json:"base_host"`
	MaxToken    int      `json:"max_token"`
	Temperature float32  `json:"temperature"`
}

func NewChatClient(apiKeys []string) *ChatClient {
	return &ChatClient{
		APIKeys:     apiKeys,
		MaxToken:    MaxToken,
		Temperature: Temperature,
	}
}

func (c *ChatClient) WithModel(model string) *ChatClient {
	if model != "" {
		c.Model = model
	}
	return c
}

func (c *ChatClient) WithBaseHost(baseHost string) *ChatClient {
	c.BaseHost = baseHost
	return c
}

// WithMaxToken 设置最大token数
func (c *ChatClient) WithMaxToken(maxToken int) *ChatClient {
	c.MaxToken = maxToken
	return c
}

// WithTemperature 设置响应灵活程度
func (c *ChatClient) WithTemperature(temperature float32) *ChatClient {
	c.Temperature = temperature
	return c
}

func (c *ChatClient) WithHttpProxy(proxyUrl string) *ChatClient {
	c.HttpProxy = proxyUrl
	return c
}
func (c *ChatClient) WithSocks5Proxy(proxyUrl string) *ChatClient {
	c.Socks5Proxy = proxyUrl
	return c
}

func (c *ChatClient) SpeakToTxt(voiceUrl string) (string, error) {
	config := c.buildConfig()
	cli := copenai.NewClientWithConfig(config)

	// 打印文件信息
	logx.Info("File: ", voiceUrl)
	info, err := os.Stat(voiceUrl)
	if err != nil {
		return "", err
	}

	logx.Info("FileInfo: ", info)

	req := copenai.AudioRequest{
		Model:       copenai.Whisper1,
		FilePath:    voiceUrl,
		Prompt:      "使用简体中文",
		Temperature: 0.5,
		Language:    "zh",
	}
	resp, err := cli.CreateTranscription(context.Background(), req)

	if err != nil {
		fmt.Printf("req chat stream params: %+v ,err:%+v", config, err)
		origin, err1 := c.MakeOpenAILoopRequest(&OpenAIRequest{
			FuncName: "CreateTranscription",
			Request:  req,
		})
		if err1 != nil {
			return "", err1
		}
		origin2, ok := origin.(copenai.AudioResponse)
		if !ok {
			return "", errors.New("Conversion failed")
		}
		resp = origin2
	}
	// 用完就删掉
	_ = os.Remove(voiceUrl)

	return resp.Text, nil
}

func (c *ChatClient) Completion(req string) (string, error) {
	config := c.buildConfig()
	cli := copenai.NewClientWithConfig(config)

	// 打印请求信息
	logx.Info("Completion req: ", req)
	request := copenai.CompletionRequest{
		Model:       copenai.GPT3TextDavinci003,
		Prompt:      req,
		MaxTokens:   c.MaxToken,
		Temperature: c.Temperature,
		TopP:        1,
	}
	completion, err := cli.CreateCompletion(context.Background(), request)

	if err != nil {
		fmt.Printf("req Completion stream params: %+v ,err:%+v", config, err)
		origin, err1 := c.MakeOpenAILoopRequest(&OpenAIRequest{
			FuncName: "CreateCompletion",
			Request:  request,
		})
		if err1 != nil {
			return "", err1
		}
		origin1, ok := origin.(copenai.CompletionResponse)
		if !ok {
			return "", errors.New("Conversion failed")
		}
		completion = origin1
	}
	txt := ""
	for _, choice := range completion.Choices {
		txt += choice.Text
	}
	logx.Info("Completion reps: ", txt)
	return txt, nil
}

func (c *ChatClient) Chat(req []ChatModelMessage) (string, error) {

	config := c.buildConfig()
	cli := copenai.NewClientWithConfig(config)

	// 打印请求信息
	logx.Info("req: ", req)

	var messages []copenai.ChatCompletionMessage

	for _, message := range req {
		messages = append(messages, copenai.ChatCompletionMessage{
			Role:    message.Role,
			Content: message.Content,
		})
	}
	request := copenai.ChatCompletionRequest{
		Model:       ChatModel,
		Messages:    messages,
		MaxTokens:   c.MaxToken,
		Temperature: c.Temperature,
		TopP:        1,
	}
	chat, err := cli.CreateChatCompletion(context.Background(), request)
	if err != nil {
		fmt.Printf("req chat params: %+v ,err:%+v", config, err)
		chatOrigin, err1 := c.MakeOpenAILoopRequest(&OpenAIRequest{
			FuncName: "CreateChatCompletion",
			Request:  request,
		})
		if err1 != nil {
			return "", err1
		}
		chat1, ok := chatOrigin.(copenai.ChatCompletionResponse)
		if !ok {
			return "", errors.New("Conversion failed")
		}
		chat = chat1
	}
	txt := ""
	for _, choice := range chat.Choices {
		txt += choice.Message.Content
	}

	return txt, nil
}

func (c *ChatClient) buildConfig() copenai.ClientConfig {
	c.WithOpenAIKey()
	config := copenai.DefaultConfig(c.APIKey)
	if c.HttpProxy != "" {
		proxyUrl, _ := url.Parse(c.HttpProxy)
		transport := &http.Transport{
			Proxy: http.ProxyURL(proxyUrl),
		}
		config.HTTPClient = &http.Client{
			Transport: transport,
		}
	} else if c.Socks5Proxy != "" {
		socks5Transport := &http.Transport{}
		dialer, _ := proxy.SOCKS5("tcp", c.Socks5Proxy, nil, proxy.Direct)
		socks5Transport.DialContext = func(ctx context.Context, network, addr string) (net.Conn, error) {
			return dialer.Dial(network, addr)
		}
		config.HTTPClient = &http.Client{
			Transport: socks5Transport,
		}
	}

	if c.BaseHost != "" {
		// trim last slash
		config.BaseURL = strings.TrimRight(c.BaseHost, "/") + "/v1"
	}
	return config
}

// ChatStream 数据流式传输
func (c *ChatClient) ChatStream(req []ChatModelMessage, channel chan string) (string, error) {

	config := c.buildConfig()

	cli := copenai.NewClientWithConfig(config)

	// 打印请求信息
	logx.Info("req: ", req)
	first := 0
	var system ChatModelMessage
	for i, msg := range req {
		if msg.Role == "system" {
			system = msg
		}
		if i%2 == 0 {
			continue
		}
		//估算长度
		if NumTokensFromMessages(req[len(req)-i-1:], ChatModel) < (3900 - c.MaxToken) {
			first = len(req) - i - 1
		} else {
			break
		}
	}

	var messages []copenai.ChatCompletionMessage

	if first != 0 {
		messages = append(messages, copenai.ChatCompletionMessage{
			Role:    system.Role,
			Content: system.Content,
		})
	}

	for _, message := range req[first:] {
		messages = append(messages, copenai.ChatCompletionMessage{
			Role:    message.Role,
			Content: message.Content,
		})
	}
	request := copenai.ChatCompletionRequest{
		Model:       ChatModel,
		Messages:    messages,
		MaxTokens:   c.MaxToken,
		Temperature: c.Temperature,
		TopP:        1,
	}
	var stream *copenai.ChatCompletionStream
	stream, err := cli.CreateChatCompletionStream(context.Background(), request)

	if err != nil {
		fmt.Printf("req chat stream params: %+v ,err:%+v", config, err)
		stream1, err1 := c.MakeOpenAILoopRequest(&OpenAIRequest{
			FuncName: "CreateChatCompletionStream",
			Request:  request,
		})
		if err1 != nil {
			return "", err1
		}
		stream2, ok := stream1.(copenai.ChatCompletionStream)
		if !ok {
			return "", errors.New("Conversion failed")
		}
		stream = &stream2
	}
	defer stream.Close()

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

func (c *ChatClient) WithKey(key string) *ChatClient {
	if key != "" {
		c.APIKey = key
	}
	return c
}

func (c *ChatClient) WithOpenAIKey() *ChatClient {
	keys := c.APIKeys
	openAiKey, err := redis.Rdb.Get(context.Background(), redis.OpenAIUsedKey).Result()
	var currentKey int
	if err == nil {
		currentKey, _ = strconv.Atoi(openAiKey)
	}
	c.APIKey = keys[currentKey]
	return c
}

func (c *ChatClient) WithNextOpenAIKey() *ChatClient {
	ctx := context.Background()
	keys := c.APIKeys
	openAiKey, err := redis.Rdb.Get(ctx, redis.OpenAIUsedKey).Result()
	var currentKey int
	if err == nil {
		currentKey, _ = strconv.Atoi(openAiKey)
	}
	if len(keys) > currentKey+1 {
		redis.Rdb.Incr(ctx, redis.OpenAIUsedKey)
		c.APIKey = keys[currentKey+1]
	} else {
		redis.Rdb.Del(ctx, redis.OpenAIUsedKey)
		c.APIKey = keys[0]
	}
	return c
}

type OpenAIRequest struct {
	Error    error
	FuncName string
	Request  interface{}
}

func (c *ChatClient) MakeOpenAILoopRequest(req *OpenAIRequest) (interface{}, error) {
	if strings.Contains(req.Error.Error(), NeedLoopErrorMessage) {
		//账号有问题
		loopTimes := len(c.APIKeys)
		for {
			fmt.Printf("loopTimes:%d \n", loopTimes)
			if loopTimes < 0 {
				fmt.Println("循环达到最大次数")
				return "", req.Error
			}
			c.WithNextOpenAIKey()
			config := c.buildConfig()

			cli := copenai.NewClientWithConfig(config)
			var result interface{}
			var resultError error

			switch req.FuncName {
			case "CreateTranscription":
				result, resultError = cli.CreateTranscription(context.Background(), req.Request.(copenai.AudioRequest))

			case "CreateCompletion":
				result, resultError = cli.CreateCompletion(context.Background(), req.Request.(copenai.CompletionRequest))

			case "Chat":
				result, resultError = cli.CreateChatCompletion(context.Background(), req.Request.(copenai.ChatCompletionRequest))

			case "ChatStream":
				result, resultError = cli.CreateChatCompletionStream(context.Background(), req.Request.(copenai.ChatCompletionRequest))

			case "CreateEmbeddings":
				result, resultError = cli.CreateEmbeddings(context.Background(), req.Request.(copenai.EmbeddingRequest))

			default:
				fmt.Println("没有匹配到对应方法")
				return nil, req.Error
			}
			if resultError != nil {
				if strings.Contains(resultError.Error(), NeedLoopErrorMessage) {
					loopTimes--
					continue
				} else {
					return "", resultError
				}
			}
			if result != nil {
				fmt.Println("1111111111111")
				return result, nil
			}
		}

	} else {
		return "", req.Error
	}
}
