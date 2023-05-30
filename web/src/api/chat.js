import service from '@/utils/request'

// @Tags Chat
// @Summary 创建Chat
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.Chat true "创建Chat"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /chat/createChat [post]
export const createChat = (data) => {
  return service({
    url: '/chat/createChat',
    method: 'post',
    data
  })
}

// @Tags Chat
// @Summary 删除Chat
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.Chat true "删除Chat"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /chat/deleteChat [delete]
export const deleteChat = (data) => {
  return service({
    url: '/chat/deleteChat',
    method: 'delete',
    data
  })
}

// @Tags Chat
// @Summary 删除Chat
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.IdsReq true "批量删除Chat"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /chat/deleteChat [delete]
export const deleteChatByIds = (data) => {
  return service({
    url: '/chat/deleteChatByIds',
    method: 'delete',
    data
  })
}

// @Tags Chat
// @Summary 更新Chat
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.Chat true "更新Chat"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"更新成功"}"
// @Router /chat/updateChat [put]
export const updateChat = (data) => {
  return service({
    url: '/chat/updateChat',
    method: 'put',
    data
  })
}

// @Tags Chat
// @Summary 用id查询Chat
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query model.Chat true "用id查询Chat"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"查询成功"}"
// @Router /chat/findChat [get]
export const findChat = (params) => {
  return service({
    url: '/chat/findChat',
    method: 'get',
    params
  })
}

// @Tags Chat
// @Summary 分页获取Chat列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query request.PageInfo true "分页获取Chat列表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /chat/getChatList [get]
export const getChatList = (params) => {
  return service({
    url: '/chat/getChatList',
    method: 'get',
    params
  })
}
