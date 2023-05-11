package openai

import (
	"errors"
	copenai "github.com/sashabaranov/go-openai"
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

	requestBody := copenai.EmbeddingRequest{
		Model: copenai.AdaEmbeddingV2,
		Input: []string{input},
	}
	var res copenai.EmbeddingResponse

	origin, err1 := c.MakeOpenAILoopRequest(&OpenAIRequest{
		FuncName: "CreateEmbeddings",
		Request:  requestBody,
	})
	if err1 != nil {
		return EmbeddingResponse{}, err1
	}
	stream2, ok := origin.(copenai.EmbeddingResponse)
	if !ok {
		return EmbeddingResponse{}, errors.New("Conversion failed")
	}
	res = stream2

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
