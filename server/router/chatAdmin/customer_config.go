package chatAdmin

import (
	"github.com/flipped-aurora/gin-vue-admin/server/api/v1"
	"github.com/flipped-aurora/gin-vue-admin/server/middleware"
	"github.com/gin-gonic/gin"
)

type CustomerConfigRouter struct {
}

// InitCustomerConfigRouter 初始化 CustomerConfig 路由信息
func (s *CustomerConfigRouter) InitCustomerConfigRouter(Router *gin.RouterGroup) {
	customerConfigRouter := Router.Group("customerConfig").Use(middleware.OperationRecord())
	customerConfigRouterWithoutRecord := Router.Group("customerConfig")
	var customerConfigApi = v1.ApiGroupApp.ChatAdminApiGroup.CustomerConfigApi
	{
		customerConfigRouter.POST("createCustomerConfig", customerConfigApi.CreateCustomerConfig)   // 新建CustomerConfig
		customerConfigRouter.DELETE("deleteCustomerConfig", customerConfigApi.DeleteCustomerConfig) // 删除CustomerConfig
		customerConfigRouter.DELETE("deleteCustomerConfigByIds", customerConfigApi.DeleteCustomerConfigByIds) // 批量删除CustomerConfig
		customerConfigRouter.PUT("updateCustomerConfig", customerConfigApi.UpdateCustomerConfig)    // 更新CustomerConfig
	}
	{
		customerConfigRouterWithoutRecord.GET("findCustomerConfig", customerConfigApi.FindCustomerConfig)        // 根据ID获取CustomerConfig
		customerConfigRouterWithoutRecord.GET("getCustomerConfigList", customerConfigApi.GetCustomerConfigList)  // 获取CustomerConfig列表
	}
}
