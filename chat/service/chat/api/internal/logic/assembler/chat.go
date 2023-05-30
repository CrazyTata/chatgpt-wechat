package assembler

import (
	"chat/service/chat/api/internal/types"
	"chat/service/chat/model"
)

func POTODTOGetChatList(po []*model.Chat) (dto []types.ChatResponse) {
	if len(po) <= 0 {
		return
	}
	var agentId string
	for _, v := range po {
		dto = append(dto, types.ChatResponse{
			Id:         v.Id,
			User:       v.User,
			MessageId:  v.MessageId,
			OpenKfId:   v.OpenKfId,
			AgentId:    agentId,
			ReqContent: v.ReqContent,
			ResContent: v.ResContent,
			CreatedAt:  v.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt:  v.UpdatedAt.Format("2006-01-02 15:04:05"),
		})
	}
	return
}
