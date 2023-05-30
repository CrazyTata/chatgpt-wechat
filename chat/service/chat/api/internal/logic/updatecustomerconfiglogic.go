package logic

import (
	"chat/service/chat/api/internal/repository"
	"chat/service/chat/api/internal/svc"
	"chat/service/chat/api/internal/types"
	"chat/service/chat/model"
	"context"
	"database/sql"
	"github.com/cockroachdb/errors"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateCustomerConfigLogic struct {
	logx.Logger
	ctx                      context.Context
	svcCtx                   *svc.ServiceContext
	customerConfigRepository *repository.CustomerConfigRepository
}

func NewUpdateCustomerConfigLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateCustomerConfigLogic {
	return &UpdateCustomerConfigLogic{
		Logger:                   logx.WithContext(ctx),
		ctx:                      ctx,
		svcCtx:                   svcCtx,
		customerConfigRepository: repository.NewCustomerConfigRepository(ctx, svcCtx),
	}
}

func (l *UpdateCustomerConfigLogic) UpdateCustomerConfig(req *types.CustomerConfig) (resp *types.Response, err error) {
	if req == nil || req.Id <= 0 {
		return nil, errors.New("缺少必传参数")
	}

	var score sql.NullFloat64
	if req.Score > 0 {
		score.Valid = true
		score.Float64 = req.Score
	}
	err = l.customerConfigRepository.Update(req.Id, &model.CustomerConfig{
		Id:               req.Id,
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
