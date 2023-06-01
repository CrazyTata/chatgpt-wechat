package assembler

import (
	"chat/service/chat/api/internal/types"
	"chat/service/chat/model"
)

func POTODTOGetCustomerList(customerPos []*model.CustomerConfig) (dto []types.CustomerConfig) {
	if len(customerPos) <= 0 {
		return
	}

	for _, v := range customerPos {
		dto = append(dto, POTODTOGetCustomer(v))
	}
	return
}

func POTODTOGetCustomer(customerPo *model.CustomerConfig) (dto types.CustomerConfig) {
	if customerPo == nil {
		return
	}
	var score float64
	if customerPo.Score.Valid {
		score = customerPo.Score.Float64
	}
	dto.Id = customerPo.Id
	dto.KfId = customerPo.KfId
	dto.KfName = customerPo.KfName
	dto.Prompt = customerPo.Prompt
	dto.PostModel = customerPo.PostModel
	dto.EmbeddingEnable = customerPo.EmbeddingEnable
	dto.EmbeddingMode = customerPo.EmbeddingMode
	dto.Score = score
	dto.TopK = customerPo.TopK
	dto.ClearContextTime = customerPo.ClearContextTime
	dto.CreatedAt = customerPo.CreatedAt.Format("2006-01-02 15:04:05")
	dto.UpdatedAt = customerPo.UpdatedAt.Format("2006-01-02 15:04:05")
	return
}
