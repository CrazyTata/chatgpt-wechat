package client

import (
	"encoding/json"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/chatAdmin"
	chatAdminReq "github.com/flipped-aurora/gin-vue-admin/server/model/chatAdmin/request"
	"github.com/flipped-aurora/gin-vue-admin/server/service/client/assembler"
	"github.com/flipped-aurora/gin-vue-admin/server/service/client/clientStruct"
	"github.com/flipped-aurora/gin-vue-admin/server/utils"
	"github.com/flipped-aurora/gin-vue-admin/server/vars"
	"time"
)

type ChatService struct {
}

func (c *ChatConfigService) GetChatInfoList(info chatAdminReq.ChatSearch) (list []chatAdmin.Chat, total int64, err error) {
	param := assembler.DTOTOPOGetChat(info)

	jsonParam, _ := json.Marshal(param)

	result, err := utils.Post(utils.GetChatServerHost()+vars.ChatGetChatUri, jsonParam, nil)

	if err != nil {
		utils.Info("GetChatInfoList utils.Post error " + err.Error())
		return
	}
	var resultInfo clientStruct.ChatPageResult
	err = json.Unmarshal(result, &resultInfo)
	if err != nil {
		utils.Info("GetChatInfoList json.Unmarshal error " + err.Error())
		return
	}
	if resultInfo.List != nil && len(resultInfo.List) > 0 {

		for _, v := range resultInfo.List {
			createdAt, _ := time.Parse(vars.TimeFormat, v.CreatedAt)
			updatedAt, _ := time.Parse(vars.TimeFormat, v.UpdatedAt)

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

func (c *ChatConfigService) ExportChatInfoList(info chatAdminReq.ChatSearch) (file string, err error) {
	param := assembler.DTOTOPOGetChat(info)

	jsonParam, _ := json.Marshal(param)

	result, err := utils.Post(utils.GetChatServerHost()+vars.ChatExportChatUri, jsonParam, nil)

	if err != nil {
		utils.Info("ExportChatInfoList utils.Post error " + err.Error())
		return
	}
	var resultInfo clientStruct.ChatExportResponse
	err = json.Unmarshal(result, &resultInfo)
	if err != nil {
		utils.Info("ExportChatInfoList json.Unmarshal error " + err.Error())
		return
	}
	file = resultInfo.File
	return
}
