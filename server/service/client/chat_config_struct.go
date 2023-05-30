package client

import "database/sql"

type GetCustomerConfigListRequest struct {
	Page           int    `json:"page"`      // 页码
	PageSize       int    `json:"page_size"` // 每页大小
	StartCreatedAt string `json:"start_created_at"`
	EndCreatedAt   string `json:"end_created_at"`
	CustomerName   string `json:"customer_name"`
	Model          string `json:"model"`
}

type FindCustomerConfigRequest struct {
	Id int64 `json:"id"`
}

type CustomerConfig struct {
	Id               int64   `json:"id,omitempty,omitempty"`
	KfId             string  `json:"kf_id,omitempty"` // 客服ID
	KfName           string  `json:"kf_name,omitempty"`
	Prompt           string  `json:"prompt,omitempty"`
	PostModel        string  `json:"post_model,omitempty"`         // 发送请求的model
	EmbeddingEnable  bool    `json:"embedding_enable,omitempty"`   // 是否启用embedding
	EmbeddingMode    string  `json:"embedding_mode,omitempty"`     // embedding的搜索模式
	Score            float64 `json:"score,omitempty"`              // 分数
	TopK             int64   `json:"top_k,omitempty"`              // topK
	ClearContextTime int64   `json:"clear_context_time,omitempty"` // 需要清理上下文的时间，按分配置，默认0不清理
	CreatedAt        string  `json:"created_at,omitempty"`         // 创建时间
	UpdatedAt        string  `json:"updated_at,omitempty"`         // 更新时间
}

type CustomerConfigResponse struct {
	Id               int64           `json:"id,omitempty,omitempty"`
	KfId             string          `json:"kf_id,omitempty"` // 客服ID
	KfName           string          `json:"kf_name,omitempty"`
	Prompt           string          `json:"prompt,omitempty"`
	PostModel        string          `json:"post_model,omitempty"`         // 发送请求的model
	EmbeddingEnable  bool            `json:"embedding_enable,omitempty"`   // 是否启用embedding
	EmbeddingMode    string          `json:"embedding_mode,omitempty"`     // embedding的搜索模式
	Score            sql.NullFloat64 `json:"score,omitempty"`              // 分数
	TopK             int64           `json:"top_k,omitempty"`              // topK
	ClearContextTime int64           `json:"clear_context_time,omitempty"` // 需要清理上下文的时间，按分配置，默认0不清理
	CreatedAt        string          `json:"created_at,omitempty"`         // 创建时间
	UpdatedAt        string          `json:"updated_at,omitempty"`         // 更新时间
}

type SyncWechatUserReq struct {
}

type SyncWechatUserReply struct {
	Message string `json:"message"`
}

type Response struct {
	Message string `json:"message"`
}

type CustomerPageResult struct {
	List     []CustomerConfigResponse `json:"list"`
	Total    int64                    `json:"total"`
	Page     int                      `json:"page"`
	PageSize int                      `json:"pageSize"`
}

type IdRequest struct {
	Id string `json:"id"`
}

type PageInfo struct {
	Page     int    `json:"page" form:"page"`         // 页码
	PageSize int    `json:"pageSize" form:"pageSize"` // 每页大小
	Keyword  string `json:"keyword" form:"keyword"`   //关键字
}

type GetApplicationConfigListRequest struct {
	Page           int    `json:"page"`      // 页码
	PageSize       int    `json:"page_size"` // 每页大小
	StartCreatedAt string `json:"start_created_at"`
	EndCreatedAt   string `json:"end_created_at"`
	AgentName      string `json:"agent_name"`
	Model          string `json:"model"`
}

type FindApplicationConfigRequest struct {
	Id int64 `json:"id"`
}

type IdsRequest struct {
	Ids []int64 `json:"id"`
}

type ApplicationConfig struct {
	Id               int64   `json:"id,omitempty"`
	AgentId          int     `json:"agent_id,omitempty"`
	AgentSecret      string  `json:"agent_secret,omitempty"`
	AgentName        string  `json:"agent_name,omitempty"`
	Model            string  `json:"model,omitempty"`
	PostModel        string  `json:"post_model,omitempty"`
	BasePrompt       string  `json:"base_prompt,omitempty"`
	Welcome          string  `json:"welcome,omitempty"`
	GroupEnable      bool    `json:"group_enable,omitempty"`
	EmbeddingEnable  bool    `json:"embedding_enable,omitempty"`
	EmbeddingMode    string  `json:"embedding_mode,omitempty"`
	Score            float64 `json:"score,omitempty"`
	TopK             int     `json:"top_k,omitempty"`
	ClearContextTime int     `json:"clear_context_time,omitempty"`
	GroupName        string  `json:"group_name,omitempty"`
	GroupChatId      string  `json:"group_chat_id,omitempty"`
}

type ApplicationConfigResponse struct {
	Id               int64           `json:"id"`
	AgentId          int64           `json:"agent_id"`
	AgentSecret      string          `json:"agent_secret"`
	AgentName        string          `json:"agent_name"`
	Model            string          `json:"model"`
	PostModel        string          `json:"post_model"`
	BasePrompt       string          `json:"base_prompt"`
	Welcome          string          `json:"welcome"`
	GroupEnable      bool            `json:"group_enable"`
	EmbeddingEnable  bool            `json:"embedding_enable"`
	EmbeddingMode    string          `json:"embedding_mode"`
	Score            sql.NullFloat64 `json:"score"`
	TopK             int64           `json:"top_k"`
	ClearContextTime int64           `json:"clear_context_time"`
	GroupName        string          `json:"group_name"`
	GroupChatId      string          `json:"group_chat_id"`
}

type ApplicationPageResult struct {
	List     []ApplicationConfigResponse `json:"list"`
	Total    int64                       `json:"total"`
	Page     int                         `json:"page"`
	PageSize int                         `json:"pageSize"`
}
