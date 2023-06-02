package chatAdmin

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	chatAdminReq "github.com/flipped-aurora/gin-vue-admin/server/model/chatAdmin/request"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type ChatApi struct {
}

// GetChatList 分页获取Chat列表
// @Tags Chat
// @Summary 分页获取Chat列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query chatAdminReq.ChatSearch true "分页获取Chat列表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /chat/getChatList [get]
func (chatApi *ChatApi) GetChatList(c *gin.Context) {
	var pageInfo chatAdminReq.ChatSearch
	err := c.ShouldBindQuery(&pageInfo)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if list, total, err := ChatConfigService.GetChatInfoList(pageInfo); err != nil {
		global.GVA_LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage("获取失败", c)
	} else {
		response.OkWithDetailed(response.PageResult{
			List:     list,
			Total:    total,
			Page:     pageInfo.Page,
			PageSize: pageInfo.PageSize,
		}, "获取成功", c)
	}
}

func (chatApi *ChatApi) ExportChatList(c *gin.Context) {
	var pageInfo chatAdminReq.ChatSearch
	err := c.ShouldBindQuery(&pageInfo)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if file, err := ChatConfigService.ExportChatInfoList(pageInfo); err != nil {
		global.GVA_LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage("获取失败", c)
	} else {
		response.OkWithDetailed(response.ExportResult{File: file}, "获取成功", c)
	}
}
