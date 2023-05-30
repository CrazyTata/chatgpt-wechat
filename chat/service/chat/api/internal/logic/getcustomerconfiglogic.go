package logic

import (
	"chat/service/chat/api/internal/repository"
	"context"
	"fmt"

	"chat/service/chat/api/internal/svc"
	"chat/service/chat/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetCustomerConfigLogic struct {
	logx.Logger
	ctx                      context.Context
	svcCtx                   *svc.ServiceContext
	customerConfigRepository *repository.CustomerConfigRepository
}

func NewGetCustomerConfigLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetCustomerConfigLogic {
	return &GetCustomerConfigLogic{
		Logger:                   logx.WithContext(ctx),
		ctx:                      ctx,
		svcCtx:                   svcCtx,
		customerConfigRepository: repository.NewCustomerConfigRepository(ctx, svcCtx),
	}
}

func (l *GetCustomerConfigLogic) GetCustomerConfig(req *types.GetCustomerConfigListRequest) (resp *types.PageResult, err error) {
	customerConfigPos, count, err := l.customerConfigRepository.GetAll(req.CustomerName, req.Model, req.StartCreatedAt, req.EndCreatedAt, "id asc", uint64(req.Page), uint64(req.PageSize))
	if err != nil {
		fmt.Printf("GetSystemConfig error: %v", err)
		return
	}
	if count <= 0 || len(customerConfigPos) <= 0 {
		return &types.PageResult{
			List:     nil,
			Total:    0,
			Page:     0,
			PageSize: 0,
		}, nil
	}
	return &types.PageResult{
		List:     customerConfigPos,
		Total:    count,
		Page:     req.Page,
		PageSize: req.PageSize,
	}, nil
}
