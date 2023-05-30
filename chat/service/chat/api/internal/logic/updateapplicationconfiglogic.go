package logic

import (
	"chat/service/chat/api/internal/repository"
	"chat/service/chat/model"
	"context"
	"database/sql"
	"errors"

	"chat/service/chat/api/internal/svc"
	"chat/service/chat/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateApplicationConfigLogic struct {
	logx.Logger
	ctx                         context.Context
	svcCtx                      *svc.ServiceContext
	applicationConfigRepository *repository.ApplicationConfigRepository
}

func NewUpdateApplicationConfigLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateApplicationConfigLogic {
	return &UpdateApplicationConfigLogic{
		Logger:                      logx.WithContext(ctx),
		ctx:                         ctx,
		svcCtx:                      svcCtx,
		applicationConfigRepository: repository.NewApplicationConfigRepository(ctx, svcCtx),
	}
}

func (l *UpdateApplicationConfigLogic) UpdateApplicationConfig(req *types.ApplicationConfig) (resp *types.Response, err error) {
	if req == nil || req.Id <= 0 {
		return nil, errors.New("缺少必传参数")
	}

	var score sql.NullFloat64
	if req.Score > 0 {
		score.Valid = true
		score.Float64 = req.Score
	}
	err = l.applicationConfigRepository.Update(req.Id, &model.ApplicationConfig{
		AgentId:          int64(req.AgentId),
		AgentSecret:      req.AgentSecret,
		AgentName:        req.AgentName,
		Model:            req.Model,
		PostModel:        req.PostModel,
		BasePrompt:       req.BasePrompt,
		Welcome:          req.Welcome,
		GroupEnable:      req.GroupEnable,
		GroupName:        req.GroupName,
		GroupChatId:      req.GroupChatId,
		EmbeddingEnable:  req.EmbeddingEnable,
		EmbeddingMode:    req.EmbeddingMode,
		Score:            score,
		TopK:             int64(req.TopK),
		ClearContextTime: int64(req.ClearContextTime),
	})
	if err != nil {
		return
	}
	return &types.Response{
		Message: "ok",
	}, nil
}
