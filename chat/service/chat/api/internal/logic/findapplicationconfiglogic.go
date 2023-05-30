package logic

import (
	"chat/service/chat/api/internal/repository"
	"context"
	"fmt"

	"chat/service/chat/api/internal/svc"
	"chat/service/chat/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type FindApplicationConfigLogic struct {
	logx.Logger
	ctx                         context.Context
	svcCtx                      *svc.ServiceContext
	applicationConfigRepository *repository.ApplicationConfigRepository
}

func NewFindApplicationConfigLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FindApplicationConfigLogic {
	return &FindApplicationConfigLogic{
		Logger:                      logx.WithContext(ctx),
		ctx:                         ctx,
		svcCtx:                      svcCtx,
		applicationConfigRepository: repository.NewApplicationConfigRepository(ctx, svcCtx),
	}
}

func (l *FindApplicationConfigLogic) FindApplicationConfig(req *types.FindApplicationConfigRequest) (resp *types.ApplicationConfig, err error) {
	applicationConfigPos, err := l.applicationConfigRepository.GetById(req.Id)
	if err != nil {
		fmt.Printf("GetSystemConfig error: %v", err)
		return
	}
	var score float64
	if applicationConfigPos.Score.Valid {
		score = applicationConfigPos.Score.Float64
	}
	return &types.ApplicationConfig{
		Id:               applicationConfigPos.Id,
		AgentId:          int(applicationConfigPos.AgentId),
		AgentSecret:      applicationConfigPos.AgentSecret,
		AgentName:        applicationConfigPos.AgentName,
		Model:            applicationConfigPos.Model,
		PostModel:        applicationConfigPos.PostModel,
		BasePrompt:       applicationConfigPos.BasePrompt,
		Welcome:          applicationConfigPos.Welcome,
		GroupEnable:      applicationConfigPos.GroupEnable,
		EmbeddingEnable:  applicationConfigPos.EmbeddingEnable,
		EmbeddingMode:    applicationConfigPos.EmbeddingMode,
		Score:            score,
		TopK:             int(applicationConfigPos.TopK),
		ClearContextTime: int(applicationConfigPos.ClearContextTime),
		GroupName:        applicationConfigPos.GroupName,
		GroupChatId:      applicationConfigPos.GroupChatId,
	}, nil
}
