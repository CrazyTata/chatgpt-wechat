// 自动生成模板ApplicationConfig
package chatAdmin

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	
)

// ApplicationConfig 结构体
type ApplicationConfig struct {
      global.GVA_MODEL
      AgentId  *int `json:"agentId" form:"agentId" gorm:"column:agent_id;comment:应用ID;size:10;"`
      AgentSecret  string `json:"agentSecret" form:"agentSecret" gorm:"column:agent_secret;comment:应用secret;size:128;"`
      AgentName  string `json:"agentName" form:"agentName" gorm:"column:agent_name;comment:应用名;size:128;"`
      Model  string `json:"model" form:"model" gorm:"column:model;comment:model;size:128;"`
      PostModel  string `json:"postModel" form:"postModel" gorm:"column:post_model;comment:发送请求的model;size:128;"`
      BasePrompt  string `json:"basePrompt" form:"basePrompt" gorm:"column:base_prompt;comment:openai 基础设定（可选）;size:1000;"`
      Welcome  string `json:"welcome" form:"welcome" gorm:"column:welcome;comment:进入应用时的欢迎语;size:1000;"`
      GroupEnable  *bool `json:"groupEnable" form:"groupEnable" gorm:"column:group_enable;comment:是否启用ChatGPT应用内部交流群;"`
      EmbeddingEnable  *bool `json:"embeddingEnable" form:"embeddingEnable" gorm:"column:embedding_enable;comment:是否启用embedding;"`
      EmbeddingMode  string `json:"embeddingMode" form:"embeddingMode" gorm:"column:embedding_mode;comment:embedding的搜索模式;size:64;"`
      Score  *float64 `json:"score" form:"score" gorm:"column:score;comment:分数;size:3;"`
      TopK  *int `json:"topK" form:"topK" gorm:"column:top_k;comment:topK;"`
      ClearContextTime  *int `json:"clearContextTime" form:"clearContextTime" gorm:"column:clear_context_time;comment:需要清理上下文的时间，按分配置，默认0不清理;size:10;"`
      GroupName  string `json:"groupName" form:"groupName" gorm:"column:group_name;comment:ChatGPT群名;size:64;"`
      GroupChatId  string `json:"groupChatId" form:"groupChatId" gorm:"column:group_chat_id;comment:ChatGPT应用内部交流群chat_id;size:128;"`
}


// TableName ApplicationConfig 表名
func (ApplicationConfig) TableName() string {
  return "application_config"
}

