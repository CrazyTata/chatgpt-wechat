package chatAdmin

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/chatAdmin"
	chatAdminReq "github.com/flipped-aurora/gin-vue-admin/server/model/chatAdmin/request"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
)

type ChatService struct {
}

// CreateChat 创建Chat记录
// Author [piexlmax](https://github.com/piexlmax)
func (chatService *ChatService) CreateChat(chat *chatAdmin.Chat) (err error) {
	err = global.GVA_DB.Create(chat).Error
	return err
}

// DeleteChat 删除Chat记录
// Author [piexlmax](https://github.com/piexlmax)
func (chatService *ChatService) DeleteChat(chat chatAdmin.Chat) (err error) {
	err = global.GVA_DB.Delete(&chat).Error
	return err
}

// DeleteChatByIds 批量删除Chat记录
// Author [piexlmax](https://github.com/piexlmax)
func (chatService *ChatService) DeleteChatByIds(ids request.IdsReq) (err error) {
	err = global.GVA_DB.Delete(&[]chatAdmin.Chat{}, "id in ?", ids.Ids).Error
	return err
}

// UpdateChat 更新Chat记录
// Author [piexlmax](https://github.com/piexlmax)
func (chatService *ChatService) UpdateChat(chat chatAdmin.Chat) (err error) {
	err = global.GVA_DB.Save(&chat).Error
	return err
}

// GetChat 根据id获取Chat记录
// Author [piexlmax](https://github.com/piexlmax)
func (chatService *ChatService) GetChat(id uint) (chat chatAdmin.Chat, err error) {
	err = global.GVA_DB.Where("id = ?", id).First(&chat).Error
	return
}

// GetChatInfoList 分页获取Chat记录
// Author [piexlmax](https://github.com/piexlmax)
func (chatService *ChatService) GetChatInfoList(info chatAdminReq.ChatSearch) (list []chatAdmin.Chat, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	// 创建db
	db := global.GVA_DB.Model(&chatAdmin.Chat{})
	var chats []chatAdmin.Chat
	// 如果有条件搜索 下方会自动创建搜索语句
	if info.StartCreatedAt != nil && info.EndCreatedAt != nil {
		db = db.Where("created_at BETWEEN ? AND ?", info.StartCreatedAt, info.EndCreatedAt)
	}
	if info.User != "" {
		db = db.Where("user = ?", info.User)
	}
	if info.OpenKfId != "" {
		db = db.Where("open_kf_id = ?", info.OpenKfId)
	}
	if info.AgentId != "" {
		db = db.Where("agent_id = ?", info.AgentId)
	}
	err = db.Count(&total).Error
	if err != nil {
		return
	}

	err = db.Limit(limit).Offset(offset).Find(&chats).Error
	return chats, total, err
}
