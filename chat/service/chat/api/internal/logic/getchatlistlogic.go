package logic

import (
	"chat/common/util"
	"chat/service/chat/api/internal/logic/assembler"
	"chat/service/chat/api/internal/repository"
	"chat/service/chat/model"
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
	wechatUser        *repository.WechatUserRepository
}

func NewGetChatListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetChatListLogic {
	return &GetChatListLogic{
		Logger:            logx.WithContext(ctx),
		ctx:               ctx,
		svcCtx:            svcCtx,
		applicationConfig: repository.NewApplicationConfigRepository(ctx, svcCtx),
		customerConfig:    repository.NewCustomerConfigRepository(ctx, svcCtx),
		chatRepository:    repository.NewChatRepository(ctx, svcCtx),
		wechatUser:        repository.NewWechatUserRepository(ctx, svcCtx),
	}
}

func (l *GetChatListLogic) GetChatList(req *types.GetChatListRequest) (resp *types.GetChatListPageResult, err error) {
	var agentId int64
	var user, openKfId string
	if req.Agent != "" {
		applicationPo, err := l.applicationConfig.GetByName(req.Agent)
		if nil != err {
			return nil, err
		}
		agentId = applicationPo.Id
	}
	if req.User != "" {
		wechatUserPo, err := l.wechatUser.GetByName(req.User)
		if nil != err {
			return nil, err
		}
		user = wechatUserPo.User
	}
	if req.Customer != "" {
		applicationPo, err := l.customerConfig.GetByName(req.Customer)
		if nil != err {
			return nil, err
		}
		openKfId = applicationPo.KfId
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
	var users, customers []string
	var applications []int64
	for _, v := range chatPos {
		if v.AgentId == 0 {
			customers = append(customers, v.OpenKfId)
			users = append(users, v.User)
		} else {
			applications = append(applications, v.AgentId)
		}
	}
	var wechatUserPos []*model.WechatUser
	var customerPos []*model.CustomerConfig
	var applicationPos []*model.ApplicationConfig
	if len(users) > 0 {

		wechatUserPos, err = l.wechatUser.GetByUsers(util.Unique(users))
		if err != nil {
			fmt.Printf("GetSystemConfig error: %v", err)
			return
		}
	}
	if len(customers) > 0 {
		customerPos, err = l.customerConfig.GetByKfIds(util.Unique(customers))
		if err != nil {
			fmt.Printf("GetSystemConfig error: %v", err)
			return
		}
	}
	if len(applications) > 0 {
		applicationPos, err = l.applicationConfig.GetByIds(util.Unique(applications))
		if err != nil {
			fmt.Printf("GetSystemConfig error: %v", err)
			return
		}
	}
	//todo 处理UpdatedAt  User OpenKfId AgentId
	chatDtos := assembler.POTODTOGetChatList(chatPos, wechatUserPos, applicationPos, customerPos)
	return &types.GetChatListPageResult{
		List:     chatDtos,
		Total:    count,
		Page:     req.Page,
		PageSize: req.PageSize,
	}, nil
}
