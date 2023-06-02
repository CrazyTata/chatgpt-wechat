package chatAdmin

import (
	"github.com/flipped-aurora/gin-vue-admin/server/api/v1"
	"github.com/gin-gonic/gin"
)

type ChatRouter struct {
}

// InitChatRouter 初始化 Chat 路由信息
func (s *ChatRouter) InitChatRouter(Router *gin.RouterGroup) {
	chatRouterWithoutRecord := Router.Group("chat")
	var chatApi = v1.ApiGroupApp.ChatAdminApiGroup.ChatApi
	{
		chatRouterWithoutRecord.GET("getChatList", chatApi.GetChatList)       // 获取Chat列表
		chatRouterWithoutRecord.GET("exportChatList", chatApi.ExportChatList) // 导出Chat列表
	}
}
