// 自动生成模板CustomerConfig
package chatAdmin

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
)

// CustomerConfig 结构体
type CustomerConfig struct {
	global.GVA_MODEL
	KfId             string   `json:"kfId" form:"kfId" gorm:"column:kf_id;comment:客服ID;size:128;"`
	KfName           string   `json:"kfName" form:"kfName" gorm:"column:kf_name;comment:;size:50;"`
	Prompt           string   `json:"prompt" form:"prompt" gorm:"column:prompt;comment:;size:1000;"`
	PostModel        string   `json:"postModel" form:"postModel" gorm:"column:post_model;comment:发送请求的model;size:128;"`
	EmbeddingEnable  *bool    `json:"embeddingEnable" form:"embeddingEnable" gorm:"column:embedding_enable;comment:是否启用embedding;"`
	EmbeddingMode    string   `json:"embeddingMode" form:"embeddingMode" gorm:"column:embedding_mode;comment:embedding的搜索模式;size:64;"`
	Score            *float64 `json:"score" form:"score" gorm:"column:score;comment:分数;size:3;"`
	TopK             *int     `json:"topK" form:"topK" gorm:"column:top_k;comment:topK;"`
	ClearContextTime *int     `json:"clearContextTime" form:"clearContextTime" gorm:"column:clear_context_time;comment:需要清理上下文的时间，按分配置，默认0不清理;size:10;"`
}

// TableName CustomerConfig 表名
func (CustomerConfig) TableName() string {
	return "customer_config"
}
