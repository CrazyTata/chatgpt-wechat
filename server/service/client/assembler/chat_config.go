package assembler

import (
	"github.com/flipped-aurora/gin-vue-admin/server/model/chatAdmin"
	"github.com/flipped-aurora/gin-vue-admin/server/service/client/clientStruct"
	"github.com/flipped-aurora/gin-vue-admin/server/vars"
	"time"
)

func POTODTOGetApplicationList(applicationPos []clientStruct.ApplicationConfigResponse) (dto []chatAdmin.ApplicationConfig) {
	if len(applicationPos) <= 0 {
		return
	}

	for _, v := range applicationPos {
		dto = append(dto, POTODTOGetApplication(v))
	}
	return
}

func POTODTOGetApplication(applicationPo clientStruct.ApplicationConfigResponse) (dto chatAdmin.ApplicationConfig) {

	score := applicationPo.Score
	topK := int(applicationPo.TopK)
	clearContextTime := int(applicationPo.ClearContextTime)
	agentId := int(applicationPo.AgentId)
	embeddingEnable := applicationPo.EmbeddingEnable
	groupEnable := applicationPo.GroupEnable
	createdAt, _ := time.Parse(vars.TimeFormat, applicationPo.CreatedAt)
	updatedAt, _ := time.Parse(vars.TimeFormat, applicationPo.UpdatedAt)

	dto.ID = uint(applicationPo.Id)
	dto.AgentId = &agentId
	dto.AgentSecret = applicationPo.AgentSecret
	dto.AgentName = applicationPo.AgentName
	dto.Model = applicationPo.Model
	dto.PostModel = applicationPo.PostModel
	dto.BasePrompt = applicationPo.BasePrompt
	dto.Welcome = applicationPo.Welcome
	dto.GroupEnable = &groupEnable
	dto.GroupName = applicationPo.GroupName
	dto.GroupChatId = applicationPo.GroupChatId
	dto.EmbeddingEnable = &embeddingEnable
	dto.EmbeddingMode = applicationPo.EmbeddingMode
	dto.Score = &score
	dto.TopK = &topK
	dto.ClearContextTime = &clearContextTime
	dto.CreatedAt = createdAt
	dto.UpdatedAt = updatedAt
	return

}

func POTODTOGetCustomerList(customerPos []clientStruct.CustomerConfigResponse) (dto []chatAdmin.CustomerConfig) {
	if len(customerPos) <= 0 {
		return
	}

	for _, v := range customerPos {
		dto = append(dto, POTODTOGetCustomer(v))
	}
	return
}

func POTODTOGetCustomer(customerPo clientStruct.CustomerConfigResponse) (dto chatAdmin.CustomerConfig) {

	score := customerPo.Score
	topK := int(customerPo.TopK)
	clearContextTime := int(customerPo.ClearContextTime)
	embeddingEnable := customerPo.EmbeddingEnable
	createdAt, _ := time.Parse(vars.TimeFormat, customerPo.CreatedAt)
	updatedAt, _ := time.Parse(vars.TimeFormat, customerPo.UpdatedAt)
	dto.ID = uint(customerPo.Id)
	dto.KfId = customerPo.KfId
	dto.KfName = customerPo.KfName
	dto.Prompt = customerPo.Prompt
	dto.PostModel = customerPo.PostModel
	dto.EmbeddingEnable = &embeddingEnable
	dto.EmbeddingMode = customerPo.EmbeddingMode
	dto.Score = &score
	dto.TopK = &topK
	dto.ClearContextTime = &clearContextTime
	dto.CreatedAt = createdAt
	dto.UpdatedAt = updatedAt
	return
}

func DTOTOPOSaveCustomer(dto *chatAdmin.CustomerConfig) (customerPo clientStruct.CustomerConfig) {

	customerPo.Id = int64(dto.ID)
	customerPo.KfId = dto.KfId
	customerPo.KfName = dto.KfName
	customerPo.Prompt = dto.Prompt
	customerPo.PostModel = dto.PostModel
	customerPo.EmbeddingEnable = *dto.EmbeddingEnable
	customerPo.EmbeddingMode = dto.EmbeddingMode
	customerPo.Score = *dto.Score
	customerPo.TopK = int64(*dto.TopK)
	customerPo.ClearContextTime = int64(*dto.ClearContextTime)
	customerPo.CreatedAt = dto.CreatedAt.Format(vars.TimeFormat)
	customerPo.UpdatedAt = dto.UpdatedAt.Format(vars.TimeFormat)

	return
}

func DTOTOPOSaveApplication(dto *chatAdmin.ApplicationConfig) (applicationPo clientStruct.ApplicationConfig) {
	applicationPo.Id = int64(dto.ID)
	applicationPo.AgentId = *dto.AgentId
	applicationPo.AgentSecret = dto.AgentSecret
	applicationPo.AgentName = dto.AgentName
	applicationPo.Model = dto.Model
	applicationPo.PostModel = dto.PostModel
	applicationPo.BasePrompt = dto.BasePrompt
	applicationPo.Welcome = dto.Welcome
	applicationPo.GroupEnable = *dto.GroupEnable
	applicationPo.GroupName = dto.GroupName
	applicationPo.GroupChatId = dto.GroupChatId
	applicationPo.EmbeddingEnable = *dto.EmbeddingEnable
	applicationPo.EmbeddingMode = dto.EmbeddingMode
	applicationPo.Score = *dto.Score
	applicationPo.TopK = *dto.TopK
	applicationPo.ClearContextTime = *dto.ClearContextTime
	applicationPo.CreatedAt = dto.CreatedAt.Format(vars.TimeFormat)
	applicationPo.UpdatedAt = dto.UpdatedAt.Format(vars.TimeFormat)
	return
}

func DTOTOPODeleteCustomer(dto chatAdmin.CustomerConfig) clientStruct.IdRequest {
	return clientStruct.IdRequest{Id: int64(dto.ID)}
}

func DTOTOPODeleteApplication(dto chatAdmin.ApplicationConfig) clientStruct.IdRequest {
	return clientStruct.IdRequest{Id: int64(dto.ID)}

}
func DTOTOPODeleteCustomers(ids []int) (res clientStruct.IdsRequest) {
	for _, v := range ids {
		res.Ids = append(res.Ids, int64(v))
	}
	return
}

func DTOTOPODeleteApplications(ids []int) (res clientStruct.IdsRequest) {
	for _, v := range ids {
		res.Ids = append(res.Ids, int64(v))
	}
	return
}
