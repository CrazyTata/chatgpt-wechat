package openai

import (
	"context"
	"fmt"
	copenai "github.com/sashabaranov/go-openai"
	"strings"
)

type (
	EmbeddingResponse struct {
		Object string         `json:"object"`
		Data   []Embedding    `json:"data"`
		Model  string         `json:"model"`
		Usage  EmbeddingUsage `json:"usage"`
	}

	EmbeddingUsage struct {
		PromptTokens int `json:"prompt_tokens"`
		TotalTokens  int `json:"total_tokens"`
	}

	Embedding struct {
		Object    string    `json:"object"`
		Embedding []float64 `json:"embedding"`
		Index     int       `json:"index"`
	}
)

func (c *ChatClient) CreateOpenAIEmbeddings(input string) (EmbeddingResponse, error) {
	config := c.buildConfig()

	cli := copenai.NewClientWithConfig(config)
	requestBody := copenai.EmbeddingRequest{
		Model: copenai.AdaEmbeddingV2,
		Input: []string{input},
	}
	res, err := cli.CreateEmbeddings(context.Background(), requestBody)

	if err != nil {
		if strings.Contains(err.Error(), "Incorrect API key provided") {
			//账号有问题
			loopTimes := len(c.APIKeys)
			for {
				fmt.Println("aaaaa")

				if loopTimes < 0 {
					fmt.Println("循环达到最大次数")
					return EmbeddingResponse{}, err
				}
				config = c.buildConfig()

				cli = copenai.NewClientWithConfig(config)
				res, err = cli.CreateEmbeddings(context.Background(), requestBody)
				if err != nil {
					if strings.Contains(err.Error(), "Incorrect API key provided") {
						loopTimes--
						continue
					} else {
						return EmbeddingResponse{}, err
					}
				} else {
					break
				}
			}

		} else {
			return EmbeddingResponse{}, err
		}
	}

	var arr []Embedding
	for i, v := range res.Data {
		var arr2 []float64
		for _, embedding := range v.Embedding {
			arr2 = append(arr2, float64(embedding))
		}
		arr = append(arr, Embedding{
			Index:     i,
			Object:    v.Object,
			Embedding: arr2,
		})

	}

	return EmbeddingResponse{
		Object: res.Object,
		Data:   arr,
		Model:  res.Model.String(),
		Usage: EmbeddingUsage{
			PromptTokens: res.Usage.PromptTokens,
			TotalTokens:  res.Usage.TotalTokens,
		},
	}, nil
}
