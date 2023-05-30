package chatAdmin

import (
	"github.com/flipped-aurora/gin-vue-admin/server/api/v1"
	"github.com/flipped-aurora/gin-vue-admin/server/middleware"
	"github.com/gin-gonic/gin"
)

type ApplicationConfigRouter struct {
}

// InitApplicationConfigRouter 初始化 ApplicationConfig 路由信息
func (s *ApplicationConfigRouter) InitApplicationConfigRouter(Router *gin.RouterGroup) {
	applicationConfigRouter := Router.Group("applicationConfig").Use(middleware.OperationRecord())
	applicationConfigRouterWithoutRecord := Router.Group("applicationConfig")
	var applicationConfigApi = v1.ApiGroupApp.ChatAdminApiGroup.ApplicationConfigApi
	{
		applicationConfigRouter.POST("createApplicationConfig", applicationConfigApi.CreateApplicationConfig)   // 新建ApplicationConfig
		applicationConfigRouter.DELETE("deleteApplicationConfig", applicationConfigApi.DeleteApplicationConfig) // 删除ApplicationConfig
		applicationConfigRouter.DELETE("deleteApplicationConfigByIds", applicationConfigApi.DeleteApplicationConfigByIds) // 批量删除ApplicationConfig
		applicationConfigRouter.PUT("updateApplicationConfig", applicationConfigApi.UpdateApplicationConfig)    // 更新ApplicationConfig
	}
	{
		applicationConfigRouterWithoutRecord.GET("findApplicationConfig", applicationConfigApi.FindApplicationConfig)        // 根据ID获取ApplicationConfig
		applicationConfigRouterWithoutRecord.GET("getApplicationConfigList", applicationConfigApi.GetApplicationConfigList)  // 获取ApplicationConfig列表
	}
}
