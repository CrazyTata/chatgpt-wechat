package logic

import (
	"chat/service/chat/api/internal/repository"
	"chat/service/chat/model"
	"context"
	"database/sql"

	"chat/service/chat/api/internal/svc"
	"chat/service/chat/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateCustomerConfigLogic struct {
	logx.Logger
	ctx                      context.Context
	svcCtx                   *svc.ServiceContext
	customerConfigRepository *repository.CustomerConfigRepository
}

func NewCreateCustomerConfigLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateCustomerConfigLogic {
	return &CreateCustomerConfigLogic{
		Logger:                   logx.WithContext(ctx),
		ctx:                      ctx,
		svcCtx:                   svcCtx,
		customerConfigRepository: repository.NewCustomerConfigRepository(ctx, svcCtx),
	}
}

func (l *CreateCustomerConfigLogic) CreateCustomerConfig(req *types.CustomerConfig) (resp *types.Response, err error) {
	var score sql.NullFloat64
	if req.Score > 0 {
		score.Valid = true
		score.Float64 = req.Score
	}
	_, err = l.customerConfigRepository.Insert(&model.CustomerConfig{
		KfId:             req.KfId,
		KfName:           req.KfName,
		Prompt:           req.Prompt,
		PostModel:        req.PostModel,
		EmbeddingEnable:  req.EmbeddingEnable,
		EmbeddingMode:    req.EmbeddingMode,
		Score:            score,
		TopK:             req.TopK,
		ClearContextTime: req.ClearContextTime,
	})
	if err != nil {
		return
	}
	return &types.Response{
		Message: "ok",
	}, nil
}
