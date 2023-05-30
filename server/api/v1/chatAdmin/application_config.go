package chatAdmin

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/chatAdmin"
	chatAdminReq "github.com/flipped-aurora/gin-vue-admin/server/model/chatAdmin/request"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	"github.com/flipped-aurora/gin-vue-admin/server/service"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type ApplicationConfigApi struct {
}

var ChatConfigService = service.ServiceGroupApp.ClientGroup.ChatConfigService

// CreateApplicationConfig 创建ApplicationConfig
// @Tags ApplicationConfig
// @Summary 创建ApplicationConfig
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body chatAdmin.ApplicationConfig true "创建ApplicationConfig"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /applicationConfig/createApplicationConfig [post]
func (applicationConfigApi *ApplicationConfigApi) CreateApplicationConfig(c *gin.Context) {
	var applicationConfig chatAdmin.ApplicationConfig
	err := c.ShouldBindJSON(&applicationConfig)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := ChatConfigService.CreateApplicationConfig(&applicationConfig); err != nil {
		global.GVA_LOG.Error("创建失败!", zap.Error(err))
		response.FailWithMessage("创建失败", c)
	} else {
		response.OkWithMessage("创建成功", c)
	}
}

// DeleteApplicationConfig 删除ApplicationConfig
// @Tags ApplicationConfig
// @Summary 删除ApplicationConfig
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body chatAdmin.ApplicationConfig true "删除ApplicationConfig"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /applicationConfig/deleteApplicationConfig [delete]
func (applicationConfigApi *ApplicationConfigApi) DeleteApplicationConfig(c *gin.Context) {
	var applicationConfig chatAdmin.ApplicationConfig
	err := c.ShouldBindJSON(&applicationConfig)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := ChatConfigService.DeleteApplicationConfig(applicationConfig); err != nil {
		global.GVA_LOG.Error("删除失败!", zap.Error(err))
		response.FailWithMessage("删除失败", c)
	} else {
		response.OkWithMessage("删除成功", c)
	}
}

// DeleteApplicationConfigByIds 批量删除ApplicationConfig
// @Tags ApplicationConfig
// @Summary 批量删除ApplicationConfig
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.IdsReq true "批量删除ApplicationConfig"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"批量删除成功"}"
// @Router /applicationConfig/deleteApplicationConfigByIds [delete]
func (applicationConfigApi *ApplicationConfigApi) DeleteApplicationConfigByIds(c *gin.Context) {
	var IDS request.IdsReq
	err := c.ShouldBindJSON(&IDS)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := ChatConfigService.DeleteApplicationConfigByIds(IDS); err != nil {
		global.GVA_LOG.Error("批量删除失败!", zap.Error(err))
		response.FailWithMessage("批量删除失败", c)
	} else {
		response.OkWithMessage("批量删除成功", c)
	}
}

// UpdateApplicationConfig 更新ApplicationConfig
// @Tags ApplicationConfig
// @Summary 更新ApplicationConfig
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body chatAdmin.ApplicationConfig true "更新ApplicationConfig"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"更新成功"}"
// @Router /applicationConfig/updateApplicationConfig [put]
func (applicationConfigApi *ApplicationConfigApi) UpdateApplicationConfig(c *gin.Context) {
	var applicationConfig chatAdmin.ApplicationConfig
	err := c.ShouldBindJSON(&applicationConfig)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := ChatConfigService.UpdateApplicationConfig(applicationConfig); err != nil {
		global.GVA_LOG.Error("更新失败!", zap.Error(err))
		response.FailWithMessage("更新失败", c)
	} else {
		response.OkWithMessage("更新成功", c)
	}
}

// FindApplicationConfig 用id查询ApplicationConfig
// @Tags ApplicationConfig
// @Summary 用id查询ApplicationConfig
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query chatAdmin.ApplicationConfig true "用id查询ApplicationConfig"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"查询成功"}"
// @Router /applicationConfig/findApplicationConfig [get]
func (applicationConfigApi *ApplicationConfigApi) FindApplicationConfig(c *gin.Context) {
	var applicationConfig chatAdmin.ApplicationConfig
	err := c.ShouldBindQuery(&applicationConfig)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if reapplicationConfig, err := ChatConfigService.GetApplicationConfig(applicationConfig.ID); err != nil {
		global.GVA_LOG.Error("查询失败!", zap.Error(err))
		response.FailWithMessage("查询失败", c)
	} else {
		response.OkWithData(gin.H{"reapplicationConfig": reapplicationConfig}, c)
	}
}

// GetApplicationConfigList 分页获取ApplicationConfig列表
// @Tags ApplicationConfig
// @Summary 分页获取ApplicationConfig列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query chatAdminReq.ApplicationConfigSearch true "分页获取ApplicationConfig列表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /applicationConfig/getApplicationConfigList [get]
func (applicationConfigApi *ApplicationConfigApi) GetApplicationConfigList(c *gin.Context) {
	var pageInfo chatAdminReq.ApplicationConfigSearch
	err := c.ShouldBindQuery(&pageInfo)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if list, total, err := ChatConfigService.GetApplicationConfigInfoList(pageInfo); err != nil {
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
