package client

import (
	"encoding/json"
	"github.com/flipped-aurora/gin-vue-admin/server/model/chatAdmin"
	chatAdminReq "github.com/flipped-aurora/gin-vue-admin/server/model/chatAdmin/request"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
	"github.com/flipped-aurora/gin-vue-admin/server/service/client/assembler"
	"github.com/flipped-aurora/gin-vue-admin/server/service/client/clientStruct"
	"github.com/flipped-aurora/gin-vue-admin/server/utils"
	"github.com/flipped-aurora/gin-vue-admin/server/vars"
)

type ChatConfigService struct {
}

func (c *ChatConfigService) CreateCustomerConfig(customerConfig *chatAdmin.CustomerConfig) (err error) {
	param := assembler.DTOTOPOSaveCustomer(customerConfig)

	jsonParam, _ := json.Marshal(param)

	_, err = utils.Post(utils.GetChatServerHost()+vars.ChatSaveCustomerConfigUri, jsonParam, nil)

	if err != nil {
		return
	}
	return
}

func (c *ChatConfigService) DeleteCustomerConfig(customerConfig chatAdmin.CustomerConfig) (err error) {
	param := assembler.DTOTOPODeleteCustomer(customerConfig)

	jsonParam, _ := json.Marshal(param)

	_, err = utils.Post(utils.GetChatServerHost()+vars.ChatDeleteCustomerConfigUri, jsonParam, nil)

	if err != nil {
		return
	}
	return
}

func (c *ChatConfigService) DeleteCustomerConfigByIds(ids request.IdsReq) (err error) {
	param := assembler.DTOTOPODeleteCustomers(ids.Ids)

	jsonParam, _ := json.Marshal(param)

	_, err = utils.Post(utils.GetChatServerHost()+vars.ChatDeleteCustomerConfigByIdsUri, jsonParam, nil)

	if err != nil {
		return
	}
	return
}

func (c *ChatConfigService) UpdateCustomerConfig(customerConfig chatAdmin.CustomerConfig) (err error) {
	param := assembler.DTOTOPOSaveCustomer(&customerConfig)

	jsonParam, _ := json.Marshal(param)

	_, err = utils.Post(utils.GetChatServerHost()+vars.ChatUpdateCustomerConfigUri, jsonParam, nil)

	if err != nil {
		return
	}
	return
}

func (c *ChatConfigService) GetCustomerConfig(id uint) (customerConfig chatAdmin.CustomerConfig, err error) {
	param := clientStruct.FindCustomerConfigRequest{
		Id: int64(id),
	}

	jsonParam, _ := json.Marshal(param)

	result, err := utils.Post(utils.GetChatServerHost()+vars.ChatFindCustomerConfigUri, jsonParam, nil)

	if err != nil {
		return
	}
	var resultInfo clientStruct.CustomerConfigResponse
	err = json.Unmarshal(result, &resultInfo)
	if err != nil {
		return
	}
	if resultInfo.Id <= 0 {
		return
	}
	customerConfig = assembler.POTODTOGetCustomer(resultInfo)
	return
}

func (c *ChatConfigService) GetCustomerConfigInfoList(info chatAdminReq.CustomerConfigSearch) (list []chatAdmin.CustomerConfig, total int64, err error) {

	param := clientStruct.GetCustomerConfigListRequest{
		Page:           info.Page,
		PageSize:       info.PageSize,
		CustomerName:   info.KfName,
		Model:          info.PostModel,
		StartCreatedAt: "",
		EndCreatedAt:   "",
	}
	if info.StartCreatedAt != nil {
		param.StartCreatedAt = info.StartCreatedAt.String()
	}
	if info.EndCreatedAt != nil {
		param.EndCreatedAt = info.EndCreatedAt.String()
	}

	jsonParam, _ := json.Marshal(param)

	result, err := utils.Post(utils.GetChatServerHost()+vars.ChatGetCustomerConfigUri, jsonParam, nil)

	if err != nil {
		return
	}
	var resultInfo clientStruct.CustomerPageResult
	err = json.Unmarshal(result, &resultInfo)
	if err != nil {
		return
	}
	if resultInfo.List == nil || len(resultInfo.List) <= 0 {
		return
	}
	list = assembler.POTODTOGetCustomerList(resultInfo.List)
	total = resultInfo.Total
	return
}

func (c *ChatConfigService) CreateApplicationConfig(applicationConfig *chatAdmin.ApplicationConfig) (err error) {
	param := assembler.DTOTOPOSaveApplication(applicationConfig)

	jsonParam, _ := json.Marshal(param)

	_, err = utils.Post(utils.GetChatServerHost()+vars.ChatSaveApplicationConfigUri, jsonParam, nil)

	if err != nil {
		return
	}
	return
}

func (c *ChatConfigService) DeleteApplicationConfig(applicationConfig chatAdmin.ApplicationConfig) (err error) {
	param := assembler.DTOTOPODeleteApplication(applicationConfig)

	jsonParam, _ := json.Marshal(param)

	_, err = utils.Post(utils.GetChatServerHost()+vars.ChatDeleteApplicationConfigUri, jsonParam, nil)

	if err != nil {
		return
	}
	return
}

func (c *ChatConfigService) DeleteApplicationConfigByIds(ids request.IdsReq) (err error) {
	param := assembler.DTOTOPODeleteApplications(ids.Ids)

	jsonParam, _ := json.Marshal(param)

	_, err = utils.Post(utils.GetChatServerHost()+vars.ChatDeleteApplicationConfigByIdsUri, jsonParam, nil)

	if err != nil {
		return
	}
	return
}

func (c *ChatConfigService) UpdateApplicationConfig(applicationConfig chatAdmin.ApplicationConfig) (err error) {
	param := assembler.DTOTOPOSaveApplication(&applicationConfig)

	jsonParam, _ := json.Marshal(param)

	_, err = utils.Post(utils.GetChatServerHost()+vars.ChatUpdateApplicationConfigUri, jsonParam, nil)

	if err != nil {
		return
	}
	return
}

func (c *ChatConfigService) GetApplicationConfig(id uint) (applicationConfig chatAdmin.ApplicationConfig, err error) {
	param := clientStruct.FindApplicationConfigRequest{
		Id: int64(id),
	}

	jsonParam, _ := json.Marshal(param)

	result, err := utils.Post(utils.GetChatServerHost()+vars.ChatFindApplicationConfigUri, jsonParam, nil)

	if err != nil {
		return
	}
	var resultInfo clientStruct.ApplicationConfigResponse
	err = json.Unmarshal(result, &resultInfo)
	if err != nil {
		return
	}
	if resultInfo.Id <= 0 {
		return
	}
	applicationConfig = assembler.POTODTOGetApplication(resultInfo)
	return
}

func (c *ChatConfigService) GetApplicationConfigInfoList(info chatAdminReq.ApplicationConfigSearch) (list []chatAdmin.ApplicationConfig, total int64, err error) {
	param := clientStruct.GetApplicationConfigListRequest{
		Page:           info.Page,
		PageSize:       info.PageSize,
		AgentName:      info.AgentName,
		Model:          info.PostModel,
		StartCreatedAt: "",
		EndCreatedAt:   "",
	}
	if info.StartCreatedAt != nil {
		param.StartCreatedAt = info.StartCreatedAt.String()
	}
	if info.EndCreatedAt != nil {
		param.EndCreatedAt = info.EndCreatedAt.String()
	}

	jsonParam, _ := json.Marshal(param)

	result, err := utils.Post(utils.GetChatServerHost()+vars.ChatGetApplicationConfigUri, jsonParam, nil)

	if err != nil {
		return
	}
	var resultInfo clientStruct.ApplicationPageResult
	err = json.Unmarshal(result, &resultInfo)
	if err != nil {
		return
	}
	if resultInfo.List == nil || len(resultInfo.List) <= 0 {
		return
	}
	list = assembler.POTODTOGetApplicationList(resultInfo.List)
	total = resultInfo.Total
	return
}
