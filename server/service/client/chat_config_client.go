package client

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/chatAdmin"
	chatAdminReq "github.com/flipped-aurora/gin-vue-admin/server/model/chatAdmin/request"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
	"github.com/flipped-aurora/gin-vue-admin/server/utils"
	"github.com/flipped-aurora/gin-vue-admin/server/vars"
)

type ChatConfigService struct {
}

func (c *ChatConfigService) CreateCustomerConfig(customerConfig *chatAdmin.CustomerConfig) (err error) {

	return errors.New("此功能暂未开放")
}

func (c *ChatConfigService) DeleteCustomerConfig(customerConfig chatAdmin.CustomerConfig) (err error) {
	return errors.New("此功能暂未开放")

}

func (c *ChatConfigService) DeleteCustomerConfigByIds(ids request.IdsReq) (err error) {
	return errors.New("此功能暂未开放")

}

func (c *ChatConfigService) UpdateCustomerConfig(customerConfig chatAdmin.CustomerConfig) (err error) {
	err = global.GVA_DB.Save(&customerConfig).Error
	return err
}

func (c *ChatConfigService) GetCustomerConfig(id uint) (customerConfig chatAdmin.CustomerConfig, err error) {
	return chatAdmin.CustomerConfig{}, errors.New("此功能暂未开放")

}

func (c *ChatConfigService) GetCustomerConfigInfoList(info chatAdminReq.CustomerConfigSearch) (list []chatAdmin.CustomerConfig, total int64, err error) {

	param := GetCustomerConfigListRequest{
		Page:           info.Page,
		PageSize:       info.PageSize,
		CustomerName:   info.KfName,
		Model:          info.PostModel,
		StartCreatedAt: "",
		EndCreatedAt:   "",
	}
	if info.StartCreatedAt != nil {
		param.StartCreatedAt = info.StartCreatedAt.String()
	}
	if info.EndCreatedAt != nil {
		param.EndCreatedAt = info.EndCreatedAt.String()
	}

	jsonParam, _ := json.Marshal(param)

	result, err := utils.Post(vars.ChatHost+vars.ChatGetCustomerConfigUri, jsonParam, nil)

	if err != nil {
		fmt.Print(err.Error())
		return
	}
	var resultInfo CustomerPageResult
	err = json.Unmarshal(result, &resultInfo)
	if err != nil {
		fmt.Print(err.Error())
		return
	}
	if resultInfo.List != nil && len(resultInfo.List) > 0 {
		for _, v := range resultInfo.List {
			var clearContextTime = int(v.ClearContextTime)
			var topK = int(v.TopK)
			var score float64
			if v.Score.Valid {
				score = v.Score.Float64
			}
			list = append(list, chatAdmin.CustomerConfig{
				KfId:             v.KfId,
				KfName:           v.KfName,
				Prompt:           v.Prompt,
				PostModel:        v.PostModel,
				EmbeddingEnable:  &v.EmbeddingEnable,
				EmbeddingMode:    v.EmbeddingMode,
				Score:            &score,
				TopK:             &topK,
				ClearContextTime: &clearContextTime,
			})
		}
	}

	total = resultInfo.Total
	return
}

func (c *ChatConfigService) CreateApplicationConfig(applicationConfig *chatAdmin.ApplicationConfig) (err error) {
	return errors.New("此功能暂未开放")
}

func (c *ChatConfigService) DeleteApplicationConfig(applicationConfig chatAdmin.ApplicationConfig) (err error) {
	return errors.New("此功能暂未开放")
}

func (c *ChatConfigService) DeleteApplicationConfigByIds(ids request.IdsReq) (err error) {
	return errors.New("此功能暂未开放")
}

func (c *ChatConfigService) UpdateApplicationConfig(applicationConfig chatAdmin.ApplicationConfig) (err error) {
	return errors.New("此功能暂未开放")
}

func (c *ChatConfigService) GetApplicationConfig(id uint) (applicationConfig chatAdmin.ApplicationConfig, err error) {
	return chatAdmin.ApplicationConfig{}, errors.New("此功能暂未开放")
}

func (c *ChatConfigService) GetApplicationConfigInfoList(info chatAdminReq.ApplicationConfigSearch) (list []chatAdmin.ApplicationConfig, total int64, err error) {
	param := GetApplicationConfigListRequest{
		Page:           info.Page,
		PageSize:       info.PageSize,
		AgentName:      info.AgentName,
		Model:          info.PostModel,
		StartCreatedAt: "",
		EndCreatedAt:   "",
	}
	if info.StartCreatedAt != nil {
		param.StartCreatedAt = info.StartCreatedAt.String()
	}
	if info.EndCreatedAt != nil {
		param.EndCreatedAt = info.EndCreatedAt.String()
	}

	jsonParam, _ := json.Marshal(param)

	result, err := utils.Post(vars.ChatHost+vars.ChatGetApplicationConfigUri, jsonParam, nil)

	if err != nil {
		fmt.Print(err.Error())
		return
	}
	fmt.Println(string(result))
	var resultInfo ApplicationPageResult
	err = json.Unmarshal(result, &resultInfo)
	if err != nil {
		fmt.Print(err.Error())
		return
	}
	if resultInfo.List != nil && len(resultInfo.List) > 0 {
		for _, v := range resultInfo.List {
			var clearContextTime = int(v.ClearContextTime)
			var topK = int(v.TopK)
			var agentId = int(v.AgentId)
			var score float64
			if v.Score.Valid {
				score = v.Score.Float64
			}
			list = append(list, chatAdmin.ApplicationConfig{
				AgentId:          &agentId,
				AgentSecret:      v.AgentSecret,
				AgentName:        v.AgentName,
				Model:            v.Model,
				PostModel:        v.PostModel,
				BasePrompt:       v.BasePrompt,
				Welcome:          v.Welcome,
				GroupEnable:      &v.GroupEnable,
				EmbeddingEnable:  &v.EmbeddingEnable,
				EmbeddingMode:    v.EmbeddingMode,
				Score:            &score,
				TopK:             &topK,
				ClearContextTime: &clearContextTime,
				GroupName:        v.GroupName,
				GroupChatId:      v.GroupChatId,
			})
		}
	}

	total = resultInfo.Total
	return
}
