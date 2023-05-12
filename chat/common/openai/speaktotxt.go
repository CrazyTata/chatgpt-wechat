package openai

import (
	"errors"
	"fmt"
	"os"

	copenai "github.com/sashabaranov/go-openai"
	"github.com/zeromicro/go-zero/core/logx"
)

type Speaker interface {
	SpeakToTxt(voiceUrl string) (string, error)
}

func (c *ChatClient) SpeakToTxt(voiceUrl string) (string, error) {
	// 打印文件信息
	logx.Info("File: ", voiceUrl)
	info, err := os.Stat(voiceUrl)
	if err != nil {
		fmt.Println(err)
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
	var resp copenai.AudioResponse
	origin, err1 := c.MakeOpenAILoopRequest(&OpenAIRequest{
		FuncName: "CreateTranscription",
		Request:  req,
	})
	if err1 != nil {
		return "", err1
	}
	origin2, ok := origin.(copenai.AudioResponse)
	if !ok {
		return "", errors.New("conversion failed")
	}
	resp = origin2

	// 用完就删掉
	_ = os.Remove(voiceUrl)

	return resp.Text, nil
}
