package client

import (
	"encoding/json"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/chatAdmin"
	chatAdminReq "github.com/flipped-aurora/gin-vue-admin/server/model/chatAdmin/request"
	"github.com/flipped-aurora/gin-vue-admin/server/utils"
	"github.com/flipped-aurora/gin-vue-admin/server/vars"
	"time"
)

type ChatService struct {
}

func (c *ChatConfigService) GetChatInfoList(info chatAdminReq.ChatSearch) (list []chatAdmin.Chat, total int64, err error) {
	param := GetChatListRequest{
		Page:           info.Page,
		PageSize:       info.PageSize,
		AgentId:        info.AgentId,
		User:           info.User,
		OpenKfId:       info.OpenKfId,
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

	result, err := utils.Post(vars.ChatHost+vars.ChatGetChatUri, jsonParam, nil)

	if err != nil {
		utils.Info("GetChatInfoList utils.Post error " + err.Error())
		return
	}
	var resultInfo ChatPageResult
	err = json.Unmarshal(result, &resultInfo)
	if err != nil {
		utils.Info("GetChatInfoList json.Unmarshal error " + err.Error())
		return
	}
	if resultInfo.List != nil && len(resultInfo.List) > 0 {
		layout := "2006-01-02 15:04:05"
		for _, v := range resultInfo.List {
			createdAt, _ := time.Parse(layout, v.CreatedAt)
			updatedAt, _ := time.Parse(layout, v.UpdatedAt)

			list = append(list, chatAdmin.Chat{
				GVA_MODEL: global.GVA_MODEL{
					ID:        uint(v.Id),
					CreatedAt: createdAt,
					UpdatedAt: updatedAt,
				},
				User:       v.User,
				MessageId:  v.MessageId,
				OpenKfId:   v.OpenKfId,
				AgentId:    v.AgentId,
				ReqContent: v.ReqContent,
				ResContent: v.ResContent,
			})
		}
	}

	total = resultInfo.Total
	return
}
