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

type FindCustomerConfigLogic struct {
	logx.Logger
	ctx                      context.Context
	svcCtx                   *svc.ServiceContext
	customerConfigRepository *repository.CustomerConfigRepository
}

func NewFindCustomerConfigLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FindCustomerConfigLogic {
	return &FindCustomerConfigLogic{
		Logger:                   logx.WithContext(ctx),
		ctx:                      ctx,
		svcCtx:                   svcCtx,
		customerConfigRepository: repository.NewCustomerConfigRepository(ctx, svcCtx),
	}
}

func (l *FindCustomerConfigLogic) FindCustomerConfig(req *types.FindCustomerConfigRequest) (resp *types.CustomerConfig, err error) {
	customerConfigPo, err := l.customerConfigRepository.GetById(req.Id)
	if err != nil {
		fmt.Printf("GetSystemConfig error: %v", err)
		return
	}

	dto := assembler.POTODTOGetCustomer(customerConfigPo)

	return &dto, nil
}
