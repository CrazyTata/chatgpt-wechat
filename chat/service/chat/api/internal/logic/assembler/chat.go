package assembler

import (
	"chat/service/chat/api/internal/types"
	"chat/service/chat/model"
)

func POTODTOGetChatList(chatPos []*model.Chat, userPos []*model.WechatUser, applicationPos []*model.ApplicationConfig, customerPos []*model.CustomerConfig) (dto []types.ChatResponse) {
	if len(chatPos) <= 0 {
		return
	}
	userMap := make(map[string]string)
	customerMap := make(map[string]string)
	applicationMap := make(map[int64]string)
	if len(userPos) > 0 {
		for _, user := range userPos {
			userMap[user.User] = user.Nickname
		}
	}

	if len(applicationPos) > 0 {
		for _, application := range applicationPos {
			applicationMap[application.AgentId] = application.AgentName
		}
	}

	if len(customerPos) > 0 {
		for _, customer := range customerPos {
			customerMap[customer.KfId] = customer.KfName
		}
	}

	for _, v := range chatPos {
		u := v.User
		var a, c string

		if value1, ok1 := userMap[u]; ok1 {
			u = value1
		}
		if value2, ok2 := applicationMap[v.AgentId]; ok2 {
			a = value2
		}
		if value3, ok3 := customerMap[v.OpenKfId]; ok3 {
			c = value3
		}
		dto = append(dto, types.ChatResponse{
			Id:         v.Id,
			User:       u,
			MessageId:  v.MessageId,
			OpenKfId:   c,
			AgentId:    a,
			ReqContent: v.ReqContent,
			ResContent: v.ResContent,
			CreatedAt:  v.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt:  v.UpdatedAt.Format("2006-01-02 15:04:05"),
		})
	}
	return
}
