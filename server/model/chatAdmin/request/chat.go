package request

import (
	"github.com/flipped-aurora/gin-vue-admin/server/model/chatAdmin"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
	"time"
)

type ChatSearch struct {
	chatAdmin.Chat
	StartCreatedAt *time.Time `json:"startCreatedAt" form:"startCreatedAt"`
	EndCreatedAt   *time.Time `json:"endCreatedAt" form:"endCreatedAt"`
	ChatType       int32      `json:"chat_type" form:"chat_type"`
	request.PageInfo
}
