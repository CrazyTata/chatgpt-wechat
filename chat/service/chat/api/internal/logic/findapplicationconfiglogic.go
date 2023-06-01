package logic

import (
	"chat/service/chat/api/internal/logic/assembler"
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
	applicationConfigPo, err := l.applicationConfigRepository.GetById(req.Id)
	if err != nil {
		fmt.Printf("GetSystemConfig error: %v", err)
		return
	}
	dto := assembler.POTODTOGetApplication(applicationConfigPo)
	return &dto, nil
}
