package clientStruct

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
	Id               int64   `json:"id"`
	KfId             string  `json:"kf_id"` // 客服ID
	KfName           string  `json:"kf_name"`
	Prompt           string  `json:"prompt"`
	PostModel        string  `json:"post_model"`         // 发送请求的model
	EmbeddingEnable  bool    `json:"embedding_enable"`   // 是否启用embedding
	EmbeddingMode    string  `json:"embedding_mode"`     // embedding的搜索模式
	Score            float64 `json:"score"`              // 分数
	TopK             int64   `json:"top_k"`              // topK
	ClearContextTime int64   `json:"clear_context_time"` // 需要清理上下文的时间，按分配置，默认0不清理
	CreatedAt        string  `json:"created_at"`         // 创建时间
	UpdatedAt        string  `json:"updated_at"`         // 更新时间
}

type CustomerConfigResponse struct {
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
	Id int64 `json:"id"`
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
	Id               int64   `json:"id"`
	AgentId          int     `json:"agent_id"`
	AgentSecret      string  `json:"agent_secret"`
	AgentName        string  `json:"agent_name"`
	Model            string  `json:"model"`
	PostModel        string  `json:"post_model"`
	BasePrompt       string  `json:"base_prompt"`
	Welcome          string  `json:"welcome"`
	GroupEnable      bool    `json:"group_enable"`
	EmbeddingEnable  bool    `json:"embedding_enable"`
	EmbeddingMode    string  `json:"embedding_mode"`
	Score            float64 `json:"score"`
	TopK             int     `json:"top_k"`
	ClearContextTime int     `json:"clear_context_time"`
	GroupName        string  `json:"group_name"`
	GroupChatId      string  `json:"group_chat_id"`
	CreatedAt        string  `json:"created_at"` // 创建时间
	UpdatedAt        string  `json:"updated_at"` // 更新时间
}

type ApplicationConfigResponse struct {
	Id               int64   `json:"id"`
	AgentId          int64   `json:"agent_id"`
	AgentSecret      string  `json:"agent_secret"`
	AgentName        string  `json:"agent_name"`
	Model            string  `json:"model"`
	PostModel        string  `json:"post_model"`
	BasePrompt       string  `json:"base_prompt"`
	Welcome          string  `json:"welcome"`
	GroupEnable      bool    `json:"group_enable"`
	EmbeddingEnable  bool    `json:"embedding_enable"`
	EmbeddingMode    string  `json:"embedding_mode"`
	Score            float64 `json:"score"`
	TopK             int64   `json:"top_k"`
	ClearContextTime int64   `json:"clear_context_time"`
	GroupName        string  `json:"group_name"`
	GroupChatId      string  `json:"group_chat_id"`
	CreatedAt        string  `json:"created_at,omitempty"` // 创建时间
	UpdatedAt        string  `json:"updated_at,omitempty"` // 更新时间
}

type ApplicationPageResult struct {
	List     []ApplicationConfigResponse `json:"list"`
	Total    int64                       `json:"total"`
	Page     int                         `json:"page"`
	PageSize int                         `json:"pageSize"`
}
