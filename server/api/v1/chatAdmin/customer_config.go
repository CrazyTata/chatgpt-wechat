package chatAdmin

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/chatAdmin"
	chatAdminReq "github.com/flipped-aurora/gin-vue-admin/server/model/chatAdmin/request"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type CustomerConfigApi struct {
}

// CreateCustomerConfig 创建CustomerConfig
// @Tags CustomerConfig
// @Summary 创建CustomerConfig
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body chatAdmin.CustomerConfig true "创建CustomerConfig"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /customerConfig/createCustomerConfig [post]
func (customerConfigApi *CustomerConfigApi) CreateCustomerConfig(c *gin.Context) {
	var customerConfig chatAdmin.CustomerConfig
	err := c.ShouldBindJSON(&customerConfig)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := ChatConfigService.CreateCustomerConfig(&customerConfig); err != nil {
		global.GVA_LOG.Error("创建失败!", zap.Error(err))
		response.FailWithMessage("创建失败", c)
	} else {
		response.OkWithMessage("创建成功", c)
	}
}

// DeleteCustomerConfig 删除CustomerConfig
// @Tags CustomerConfig
// @Summary 删除CustomerConfig
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body chatAdmin.CustomerConfig true "删除CustomerConfig"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /customerConfig/deleteCustomerConfig [delete]
func (customerConfigApi *CustomerConfigApi) DeleteCustomerConfig(c *gin.Context) {
	var customerConfig chatAdmin.CustomerConfig
	err := c.ShouldBindJSON(&customerConfig)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := ChatConfigService.DeleteCustomerConfig(customerConfig); err != nil {
		global.GVA_LOG.Error("删除失败!", zap.Error(err))
		response.FailWithMessage("删除失败", c)
	} else {
		response.OkWithMessage("删除成功", c)
	}
}

// DeleteCustomerConfigByIds 批量删除CustomerConfig
// @Tags CustomerConfig
// @Summary 批量删除CustomerConfig
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.IdsReq true "批量删除CustomerConfig"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"批量删除成功"}"
// @Router /customerConfig/deleteCustomerConfigByIds [delete]
func (customerConfigApi *CustomerConfigApi) DeleteCustomerConfigByIds(c *gin.Context) {
	var IDS request.IdsReq
	err := c.ShouldBindJSON(&IDS)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := ChatConfigService.DeleteCustomerConfigByIds(IDS); err != nil {
		global.GVA_LOG.Error("批量删除失败!", zap.Error(err))
		response.FailWithMessage("批量删除失败", c)
	} else {
		response.OkWithMessage("批量删除成功", c)
	}
}

// UpdateCustomerConfig 更新CustomerConfig
// @Tags CustomerConfig
// @Summary 更新CustomerConfig
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body chatAdmin.CustomerConfig true "更新CustomerConfig"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"更新成功"}"
// @Router /customerConfig/updateCustomerConfig [put]
func (customerConfigApi *CustomerConfigApi) UpdateCustomerConfig(c *gin.Context) {
	var customerConfig chatAdmin.CustomerConfig
	err := c.ShouldBindJSON(&customerConfig)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := ChatConfigService.UpdateCustomerConfig(customerConfig); err != nil {
		global.GVA_LOG.Error("更新失败!", zap.Error(err))
		response.FailWithMessage("更新失败", c)
	} else {
		response.OkWithMessage("更新成功", c)
	}
}

// FindCustomerConfig 用id查询CustomerConfig
// @Tags CustomerConfig
// @Summary 用id查询CustomerConfig
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query chatAdmin.CustomerConfig true "用id查询CustomerConfig"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"查询成功"}"
// @Router /customerConfig/findCustomerConfig [get]
func (customerConfigApi *CustomerConfigApi) FindCustomerConfig(c *gin.Context) {
	var customerConfig chatAdmin.CustomerConfig
	err := c.ShouldBindQuery(&customerConfig)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if recustomerConfig, err := ChatConfigService.GetCustomerConfig(customerConfig.ID); err != nil {
		global.GVA_LOG.Error("查询失败!", zap.Error(err))
		response.FailWithMessage("查询失败", c)
	} else {
		response.OkWithData(gin.H{"recustomerConfig": recustomerConfig}, c)
	}
}

// GetCustomerConfigList 分页获取CustomerConfig列表
// @Tags CustomerConfig
// @Summary 分页获取CustomerConfig列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query chatAdminReq.CustomerConfigSearch true "分页获取CustomerConfig列表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /customerConfig/getCustomerConfigList [get]
func (customerConfigApi *CustomerConfigApi) GetCustomerConfigList(c *gin.Context) {
	var pageInfo chatAdminReq.CustomerConfigSearch
	err := c.ShouldBindQuery(&pageInfo)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if list, total, err := ChatConfigService.GetCustomerConfigInfoList(pageInfo); err != nil {
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
