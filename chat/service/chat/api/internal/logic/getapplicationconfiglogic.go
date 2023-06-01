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

type GetApplicationConfigLogic struct {
	logx.Logger
	ctx                         context.Context
	svcCtx                      *svc.ServiceContext
	applicationConfigRepository *repository.ApplicationConfigRepository
}

func NewGetApplicationConfigLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetApplicationConfigLogic {
	return &GetApplicationConfigLogic{
		Logger:                      logx.WithContext(ctx),
		ctx:                         ctx,
		svcCtx:                      svcCtx,
		applicationConfigRepository: repository.NewApplicationConfigRepository(ctx, svcCtx),
	}
}

func (l *GetApplicationConfigLogic) GetApplicationConfig(req *types.GetApplicationConfigListRequest) (resp *types.PageResult, err error) {
	applicationConfigPos, count, err := l.applicationConfigRepository.GetAll(req.AgentName, req.Model, req.StartCreatedAt, req.EndCreatedAt, "id asc", uint64(req.Page), uint64(req.PageSize))
	if err != nil {
		fmt.Printf("GetSystemConfig error: %v", err)
		return
	}
	if count <= 0 || len(applicationConfigPos) <= 0 {
		return &types.PageResult{
			List:     nil,
			Total:    0,
			Page:     0,
			PageSize: 0,
		}, nil
	}
	applicationConfigDto := assembler.POTODTOGetApplicationList(applicationConfigPos)
	return &types.PageResult{
		List:     applicationConfigDto,
		Total:    count,
		Page:     req.Page,
		PageSize: req.PageSize,
	}, nil
}
