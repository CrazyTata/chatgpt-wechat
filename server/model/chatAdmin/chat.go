// 自动生成模板Chat
package chatAdmin

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
)

// Chat 结构体
type Chat struct {
	global.GVA_MODEL
	User       string `json:"user" form:"user" gorm:"column:user;comment:weCom用户标识/customer用户标识;size:191;"`
	MessageId  string `json:"message_id" form:"message_id" gorm:"column:message_id;comment:message_id customer消息唯一ID;size:191;"`
	OpenKfId   string `json:"open_kf_id" form:"open_kf_id" gorm:"column:open_kf_id;comment:客服标识;size:191;"`
	AgentId    string `json:"agent_id" form:"agent_id" gorm:"column:agent_id;comment:应用ID;size:20;"`
	ReqContent string `json:"req_content" form:"req_content" gorm:"column:req_content;comment:用户发送内容;size:500;"`
	ResContent string `json:"res_content" form:"res_content" gorm:"column:res_content;comment:openai响应内容;"`
}

// TableName Chat 表名
func (Chat) TableName() string {
	return "chat"
}
