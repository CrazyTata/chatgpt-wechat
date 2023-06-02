package assembler

import (
	chatAdminReq "github.com/flipped-aurora/gin-vue-admin/server/model/chatAdmin/request"
	"github.com/flipped-aurora/gin-vue-admin/server/service/client/clientStruct"
)

func DTOTOPOGetChat(info chatAdminReq.ChatSearch) (res clientStruct.GetChatListRequest) {
	param := clientStruct.GetChatListRequest{
		Page:           info.Page,
		PageSize:       info.PageSize,
		Agent:          info.AgentId,
		User:           info.User,
		Customer:       info.OpenKfId,
		ChatType:       info.ChatType,
		StartCreatedAt: "",
		EndCreatedAt:   "",
	}
	if info.StartCreatedAt != nil {
		param.StartCreatedAt = info.StartCreatedAt.String()
	}
	if info.EndCreatedAt != nil {
		param.EndCreatedAt = info.EndCreatedAt.String()
	}
	res = param
	return
}

func Export(file string) (res clientStruct.ChatExportResponse) {
	res.File = file
	return
}
