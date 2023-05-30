package logic

import (
	"chat/service/chat/api/internal/repository"
	"chat/service/chat/api/internal/svc"
	"chat/service/chat/api/internal/types"
	"chat/service/chat/model"
	"context"
	"database/sql"
	"github.com/zeromicro/go-zero/core/logx"
)

type CreateApplicationConfigLogic struct {
	logx.Logger
	ctx                         context.Context
	svcCtx                      *svc.ServiceContext
	applicationConfigRepository *repository.ApplicationConfigRepository
}

func NewCreateApplicationConfigLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateApplicationConfigLogic {
	return &CreateApplicationConfigLogic{
		Logger:                      logx.WithContext(ctx),
		ctx:                         ctx,
		svcCtx:                      svcCtx,
		applicationConfigRepository: repository.NewApplicationConfigRepository(ctx, svcCtx),
	}
}

func (l *CreateApplicationConfigLogic) CreateApplicationConfig(req *types.ApplicationConfig) (resp *types.Response, err error) {
	var score sql.NullFloat64
	if req.Score > 0 {
		score.Valid = true
		score.Float64 = req.Score
	}
	_, err = l.applicationConfigRepository.Insert(&model.ApplicationConfig{
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
