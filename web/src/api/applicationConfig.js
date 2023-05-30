import service from '@/utils/request'

// @Tags ApplicationConfig
// @Summary 创建ApplicationConfig
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.ApplicationConfig true "创建ApplicationConfig"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /applicationConfig/createApplicationConfig [post]
export const createApplicationConfig = (data) => {
  return service({
    url: '/applicationConfig/createApplicationConfig',
    method: 'post',
    data
  })
}

// @Tags ApplicationConfig
// @Summary 删除ApplicationConfig
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.ApplicationConfig true "删除ApplicationConfig"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /applicationConfig/deleteApplicationConfig [delete]
export const deleteApplicationConfig = (data) => {
  return service({
    url: '/applicationConfig/deleteApplicationConfig',
    method: 'delete',
    data
  })
}

// @Tags ApplicationConfig
// @Summary 删除ApplicationConfig
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.IdsReq true "批量删除ApplicationConfig"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /applicationConfig/deleteApplicationConfig [delete]
export const deleteApplicationConfigByIds = (data) => {
  return service({
    url: '/applicationConfig/deleteApplicationConfigByIds',
    method: 'delete',
    data
  })
}

// @Tags ApplicationConfig
// @Summary 更新ApplicationConfig
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.ApplicationConfig true "更新ApplicationConfig"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"更新成功"}"
// @Router /applicationConfig/updateApplicationConfig [put]
export const updateApplicationConfig = (data) => {
  return service({
    url: '/applicationConfig/updateApplicationConfig',
    method: 'put',
    data
  })
}

// @Tags ApplicationConfig
// @Summary 用id查询ApplicationConfig
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query model.ApplicationConfig true "用id查询ApplicationConfig"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"查询成功"}"
// @Router /applicationConfig/findApplicationConfig [get]
export const findApplicationConfig = (params) => {
  return service({
    url: '/applicationConfig/findApplicationConfig',
    method: 'get',
    params
  })
}

// @Tags ApplicationConfig
// @Summary 分页获取ApplicationConfig列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query request.PageInfo true "分页获取ApplicationConfig列表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /applicationConfig/getApplicationConfigList [get]
export const getApplicationConfigList = (params) => {
  return service({
    url: '/applicationConfig/getApplicationConfigList',
    method: 'get',
    params
  })
}
