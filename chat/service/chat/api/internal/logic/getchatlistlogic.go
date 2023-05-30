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

type GetChatListLogic struct {
	logx.Logger
	ctx               context.Context
	svcCtx            *svc.ServiceContext
	chatRepository    *repository.ChatRepository
	applicationConfig *repository.ApplicationConfigRepository
	customerConfig    *repository.CustomerConfigRepository
}

func NewGetChatListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetChatListLogic {
	return &GetChatListLogic{
		Logger:            logx.WithContext(ctx),
		ctx:               ctx,
		svcCtx:            svcCtx,
		applicationConfig: repository.NewApplicationConfigRepository(ctx, svcCtx),
		customerConfig:    repository.NewCustomerConfigRepository(ctx, svcCtx),
		chatRepository:    repository.NewChatRepository(ctx, svcCtx),
	}
}

func (l *GetChatListLogic) GetChatList(req *types.GetChatListRequest) (resp *types.GetChatListPageResult, err error) {
	var agentId, user, openKfId string
	if req.AgentId != "" {
		//todo
	}
	if req.User != "" {
		//todo
	}
	if req.OpenKfId != "" {
		//todo
	}
	chatPos, count, err := l.chatRepository.GetAll(agentId, openKfId, user, req.StartCreatedAt, req.EndCreatedAt, "id asc", uint64(req.Page), uint64(req.PageSize))
	if err != nil {
		fmt.Printf("GetSystemConfig error: %v", err)
		return
	}

	if count <= 0 || len(chatPos) <= 0 {
		return &types.GetChatListPageResult{
			List:     nil,
			Total:    0,
			Page:     0,
			PageSize: 0,
		}, nil
	}
	//todo 处理UpdatedAt  User OpenKfId AgentId
	chatDtos := assembler.POTODTOGetChatList(chatPos)
	return &types.GetChatListPageResult{
		List:     chatDtos,
		Total:    count,
		Page:     req.Page,
		PageSize: req.PageSize,
	}, nil
}
