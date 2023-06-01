package repository

import (
	"chat/service/chat/api/internal/svc"
	"chat/service/chat/model"
	"context"
	"database/sql"
	"github.com/Masterminds/squirrel"
	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/logx"
)

type CustomerConfigRepository struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCustomerConfigRepository(ctx context.Context, svcCtx *svc.ServiceContext) *CustomerConfigRepository {
	return &CustomerConfigRepository{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CustomerConfigRepository) GetAll(agentName, model, startTime, endTime, order string, page, limit uint64) (CustomerConfigPo []*model.CustomerConfig, count int64, err error) {

	countBuilder := l.svcCtx.CustomerConfigModel.CountBuilder("id")
	rowBuilder := l.svcCtx.CustomerConfigModel.RowBuilder()
	if agentName != "" {
		countBuilder = countBuilder.Where(squirrel.Eq{"kf_name": agentName})
		rowBuilder = rowBuilder.Where(squirrel.Eq{"kf_name": agentName})
	}
	if model != "" {
		countBuilder = countBuilder.Where(squirrel.Eq{"post_model": model})
		rowBuilder = rowBuilder.Where(squirrel.Eq{"post_model": model})
	}

	if startTime != "" {
		countBuilder = countBuilder.Where("created_at >= ?", startTime)
		rowBuilder = rowBuilder.Where("created_at >= ?", startTime)
	}

	if endTime != "" {
		countBuilder = countBuilder.Where("created_at < ?", endTime)
		rowBuilder = rowBuilder.Where("created_at < ?", endTime)
	}

	count, err = l.svcCtx.CustomerConfigModel.FindCount(context.Background(), countBuilder)
	if err != nil {
		return
	}
	if count <= 0 {
		return nil, 0, nil
	}

	rowBuilder = rowBuilder.OrderBy(order)
	if limit != 0 {
		offset := (page - 1) * limit
		rowBuilder = rowBuilder.Limit(limit).Offset(offset)
	}
	CustomerConfigPo, err = l.svcCtx.CustomerConfigModel.FindAll(context.Background(), rowBuilder)
	if err != nil {
		return
	}
	return
}

func (l *CustomerConfigRepository) GetById(id int64) (CustomerConfigPo *model.CustomerConfig, err error) {
	return l.svcCtx.CustomerConfigModel.FindOne(context.Background(), id)
}

func (l *CustomerConfigRepository) Insert(CustomerConfigPo *model.CustomerConfig) (sql.Result, error) {
	l.svcCtx.CustomerConfigModel.BuildFiled(nil, CustomerConfigPo)
	return l.svcCtx.CustomerConfigModel.Insert(context.Background(), CustomerConfigPo)
}

func (l *CustomerConfigRepository) Update(id int64, CustomerConfigPo *model.CustomerConfig) error {
	old, err := l.GetById(id)
	if err != nil {
		return err
	}
	if old == nil || old.Id <= 0 {
		return errors.New("record not find")
	}
	l.svcCtx.CustomerConfigModel.BuildFiled(old, CustomerConfigPo)
	return l.svcCtx.CustomerConfigModel.Update(context.Background(), CustomerConfigPo)
}

func (l *CustomerConfigRepository) GetByName(kfName string) (customerPo *model.CustomerConfig, err error) {

	customerPo, err = l.svcCtx.CustomerConfigModel.FindOneByQuery(context.Background(),
		l.svcCtx.CustomerConfigModel.RowBuilder().Where(squirrel.Eq{"kf_name": kfName}),
	)
	return
}

func (l *CustomerConfigRepository) GetByKfIds(kfId []string) (CustomerConfigPo []*model.CustomerConfig, err error) {
	return l.svcCtx.CustomerConfigModel.FindAll(context.Background(),
		l.svcCtx.CustomerConfigModel.RowBuilder().Where(squirrel.Eq{"kf_id": kfId}),
	)
}

func (l *CustomerConfigRepository) Delete(id int64) error {
	return l.svcCtx.CustomerConfigModel.Delete(context.Background(), id)
}
