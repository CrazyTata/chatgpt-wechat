package clientStruct

type GetChatListRequest struct {
	Page           int    `json:"page"`      // 页码
	PageSize       int    `json:"page_size"` // 每页大小
	StartCreatedAt string `json:"start_created_at"`
	EndCreatedAt   string `json:"end_created_at"`
	Agent          string `json:"agent"`
	User           string `json:"user"`
	Customer       string `json:"customer"`
	ChatType       int32  `json:"chat_type"`
}

type ChatResponse struct {
	Id         int64  `json:"id"`
	User       string `json:"user"`        // weCom用户标识/customer用户标识
	MessageId  string `json:"message_id"`  // message_id customer消息唯一ID
	OpenKfId   string `json:"open_kf_id"`  // 客服标识
	AgentId    string `json:"agent_id"`    // 应用ID
	ReqContent string `json:"req_content"` // 用户发送内容
	ResContent string `json:"res_content"` // openai响应内容
	CreatedAt  string `json:"created_at"`  // 创建时间
	UpdatedAt  string `json:"updated_at"`  // 更新时间
}

type ChatPageResult struct {
	List     []ChatResponse `json:"list"`
	Total    int64          `json:"total"`
	Page     int            `json:"page"`
	PageSize int            `json:"pageSize"`
}

type ChatExportResponse struct {
	File string `json:"file"`
}
