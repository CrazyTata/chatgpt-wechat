package openai

import (
	"chat/common/redis"
	"chat/common/util"
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"net/url"
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
	Origin      string   `json:"origin"`
	Engine      string   `json:"engine"`
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

// WithOrigin 设置origin
func (c *ChatClient) WithOrigin(origin string) *ChatClient {
	c.Origin = origin
	return c
}

// WithEngine 设置engine
func (c *ChatClient) WithEngine(engine string) *ChatClient {
	c.Engine = engine
	return c
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

func (c *ChatClient) Completion(req string) (string, error) {
	// 打印请求信息
	logx.Info("Completion req: ", req)
	request := copenai.CompletionRequest{
		Model:       copenai.GPT3TextDavinci003,
		Prompt:      req,
		MaxTokens:   c.MaxToken,
		Temperature: c.Temperature,
		TopP:        1,
	}
	var completion copenai.CompletionResponse
	origin, err1 := c.MakeOpenAILoopRequest(&OpenAIRequest{
		FuncName: "CreateCompletion",
		Request:  request,
	})
	if err1 != nil {
		return "", err1
	}
	origin1, ok := origin.(copenai.CompletionResponse)
	if !ok {
		return "", errors.New("conversion failed")
	}
	completion = origin1

	txt := ""
	for _, choice := range completion.Choices {
		txt += choice.Text
	}
	logx.Info("Completion reps: ", txt)
	return txt, nil
}

func (c *ChatClient) buildConfig() copenai.ClientConfig {
	c.WithOpenAIKey()
	fmt.Printf("apikye:%s \n", c.APIKey)
	config := copenai.DefaultConfig(c.APIKey)
	if c.Origin == "azure" || c.Origin == "azure_ad" {
		config = copenai.DefaultAzureConfig(c.APIKey, c.BaseHost)

		// 默认只使用 一个部署的模型
		config.AzureModelMapperFunc = func(model string) string {
			//azureModelMapping := map[string]string{
			//	defaultModelType: c.Engine,
			//}
			//return azureModelMapping[model]
			return c.Engine
		}
	}
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

	if c.BaseHost != "" && c.Origin == "open_ai" {
		// trim last slash
		config.BaseURL = strings.TrimRight(c.BaseHost, "/") + "/v1"
	}
	return config
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

//
//func (c *ChatClient) MakeOpenAILoopRequestV2(req *OpenAIRequest) (interface{}, error) {
//	if strings.Contains(req.Error.Error(), NeedLoopErrorMessage) {
//		//账号有问题
//		loopTimes := len(c.APIKeys)
//		for {
//			if loopTimes < 0 {
//				fmt.Println("循环达到最大次数")
//				return "", req.Error
//			}
//			fmt.Printf("loopTimes:%d \n", loopTimes)
//			c.WithNextOpenAIKey()
//			config := c.buildConfig()
//
//			cli := copenai.NewClientWithConfig(config)
//			var result interface{}
//			var resultError error
//
//			switch req.FuncName {
//			case "CreateTranscription":
//				result, resultError = cli.CreateTranscription(context.Background(), req.Request.(copenai.AudioRequest))
//
//			case "CreateCompletion":
//				result, resultError = cli.CreateCompletion(context.Background(), req.Request.(copenai.CompletionRequest))
//
//			case "CreateChatCompletion":
//				result, resultError = cli.CreateChatCompletion(context.Background(), req.Request.(copenai.ChatCompletionRequest))
//
//			case "CreateChatCompletionStream":
//				resultStream, resultStreamError := cli.CreateChatCompletionStream(context.Background(), req.Request.(copenai.ChatCompletionRequest))
//				if nil != resultStreamError {
//					return nil, resultStreamError
//				}
//				_, resultError = resultStream.Recv()
//				if resultError == nil || !strings.Contains(resultError.Error(), NeedLoopErrorMessage) {
//					result = resultStream
//				}
//			case "CreateEmbeddings":
//				result, resultError = cli.CreateEmbeddings(context.Background(), req.Request.(copenai.EmbeddingRequest))
//
//			default:
//				fmt.Println("没有匹配到对应方法" + req.FuncName)
//				return nil, req.Error
//			}
//			if resultError != nil {
//				if strings.Contains(resultError.Error(), NeedLoopErrorMessage) {
//					loopTimes--
//					continue
//				} else {
//					return "", resultError
//				}
//			}
//
//			fmt.Println("return the result")
//			return result, nil
//		}
//
//	} else {
//		return "", req.Error
//	}
//}

func (c *ChatClient) MakeOpenAILoopRequest(req *OpenAIRequest) (interface{}, error) {
	//账号有问题
	loopTimes := len(c.APIKeys)
	var resultError error
	for {
		if loopTimes < len(c.APIKeys) {
			if loopTimes < 0 {
				fmt.Println("循环达到最大次数")
				return "", resultError
			}
			fmt.Printf("loopTimes:%d \n", loopTimes)
			c.WithNextOpenAIKey()
		}
		config := c.buildConfig()
		config.APIVersion = "2022-12-01"
		cli := copenai.NewClientWithConfig(config)
		var result interface{}
		switch req.FuncName {
		case "CreateTranscription":
			result, resultError = cli.CreateTranscription(context.Background(), req.Request.(copenai.AudioRequest))

		case "CreateCompletion":
			result, resultError = cli.CreateCompletion(context.Background(), req.Request.(copenai.CompletionRequest))

		case "CreateChatCompletion":
			result, resultError = cli.CreateChatCompletion(context.Background(), req.Request.(copenai.ChatCompletionRequest))

		case "CreateChatCompletionStream":
			resultStream, resultStreamError := cli.CreateChatCompletionStream(context.Background(), req.Request.(copenai.ChatCompletionRequest))
			if nil != resultStreamError {
				return nil, resultStreamError
			}
			_, resultError = resultStream.Recv()
			if resultError == nil || !strings.Contains(resultError.Error(), NeedLoopErrorMessage) {
				result = resultStream
			}
		case "CreateEmbeddings":
			result, resultError = cli.CreateEmbeddings(context.Background(), req.Request.(copenai.EmbeddingRequest))

		default:
			fmt.Println("没有匹配到对应方法" + req.FuncName)
			return nil, fmt.Errorf("没有匹配到对应方法")
		}
		util.Info(fmt.Sprintf("MakeOpenAILoopRequest dealMethod:%s,params: %+v ,err:%+v , response:%+v \n\n", req.FuncName, req.Request, resultError, result))
		if resultError != nil {
			if strings.Contains(resultError.Error(), NeedLoopErrorMessage) {
				fmt.Printf("MakeOpenAILoopRequest params: %+v ,err:%+v ,dealMethod:%s \n", req.Request, resultError, req.FuncName)
				loopTimes--
				continue
			} else {
				return "", resultError
			}
		}

		return result, nil
	}

}

//估算长度

func (c *ChatClient) GetNumTokens(message string) int {
	var messages []ChatModelMessage
	messages = append(messages, ChatModelMessage{
		Role:    "user",
		Content: message,
	})

	return NumTokensFromMessages(messages, ChatModel)
}
