import service from '@/utils/request'

// @Tags CustomerConfig
// @Summary 创建CustomerConfig
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.CustomerConfig true "创建CustomerConfig"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /customerConfig/createCustomerConfig [post]
export const createCustomerConfig = (data) => {
  return service({
    url: '/customerConfig/createCustomerConfig',
    method: 'post',
    data
  })
}

// @Tags CustomerConfig
// @Summary 删除CustomerConfig
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.CustomerConfig true "删除CustomerConfig"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /customerConfig/deleteCustomerConfig [delete]
export const deleteCustomerConfig = (data) => {
  return service({
    url: '/customerConfig/deleteCustomerConfig',
    method: 'delete',
    data
  })
}

// @Tags CustomerConfig
// @Summary 删除CustomerConfig
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.IdsReq true "批量删除CustomerConfig"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /customerConfig/deleteCustomerConfig [delete]
export const deleteCustomerConfigByIds = (data) => {
  return service({
    url: '/customerConfig/deleteCustomerConfigByIds',
    method: 'delete',
    data
  })
}

// @Tags CustomerConfig
// @Summary 更新CustomerConfig
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.CustomerConfig true "更新CustomerConfig"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"更新成功"}"
// @Router /customerConfig/updateCustomerConfig [put]
export const updateCustomerConfig = (data) => {
  return service({
    url: '/customerConfig/updateCustomerConfig',
    method: 'put',
    data
  })
}

// @Tags CustomerConfig
// @Summary 用id查询CustomerConfig
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query model.CustomerConfig true "用id查询CustomerConfig"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"查询成功"}"
// @Router /customerConfig/findCustomerConfig [get]
export const findCustomerConfig = (params) => {
  return service({
    url: '/customerConfig/findCustomerConfig',
    method: 'get',
    params
  })
}

// @Tags CustomerConfig
// @Summary 分页获取CustomerConfig列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query request.PageInfo true "分页获取CustomerConfig列表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /customerConfig/getCustomerConfigList [get]
export const getCustomerConfigList = (params) => {
  return service({
    url: '/customerConfig/getCustomerConfigList',
    method: 'get',
    params
  })
}
