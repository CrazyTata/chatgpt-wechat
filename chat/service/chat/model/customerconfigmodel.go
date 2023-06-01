package model

import (
	"chat/common/util"
	"context"
	"github.com/Masterminds/squirrel"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ CustomerConfigModel = (*customCustomerConfigModel)(nil)

type (
	// CustomerConfigModel is an interface to be customized, add more methods here,
	// and implement the added methods in customCustomerConfigModel.
	CustomerConfigModel interface {
		customerConfigModel
		FindOneByQuery(ctx context.Context, rowBuilder squirrel.SelectBuilder) (*CustomerConfig, error)
		RowBuilder() squirrel.SelectBuilder
		FindAll(ctx context.Context, rowBuilder squirrel.SelectBuilder) ([]*CustomerConfig, error)
		CountBuilder(field string) squirrel.SelectBuilder
		FindCount(ctx context.Context, countBuilder squirrel.SelectBuilder) (int64, error)
		BuildFiled(old, new *CustomerConfig)
	}

	customCustomerConfigModel struct {
		*defaultCustomerConfigModel
	}
)

// NewCustomerConfigModel returns a model for the database table.
func NewCustomerConfigModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) CustomerConfigModel {
	return &customCustomerConfigModel{
		defaultCustomerConfigModel: newCustomerConfigModel(conn, c, opts...),
	}
}

func (m *defaultCustomerConfigModel) FindOneByQuery(ctx context.Context, rowBuilder squirrel.SelectBuilder) (*CustomerConfig, error) {

	query, values, err := rowBuilder.ToSql()
	if err != nil {
		return nil, err
	}

	var resp CustomerConfig
	err = m.QueryRowNoCacheCtx(ctx, &resp, query, values...)
	if err != nil {
		if err == sqlx.ErrNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &resp, nil
}

// export logic
func (m *defaultCustomerConfigModel) RowBuilder() squirrel.SelectBuilder {
	return squirrel.Select(customerConfigRows).From(m.table).Where(squirrel.Eq{"is_deleted": 0})
}

func (m *defaultCustomerConfigModel) FindAll(ctx context.Context, rowBuilder squirrel.SelectBuilder) ([]*CustomerConfig, error) {

	query, values, err := rowBuilder.ToSql()
	if err != nil {
		return nil, err
	}

	var resp []*CustomerConfig
	err = m.QueryRowsNoCacheCtx(ctx, &resp, query, values...)
	switch err {
	case nil:
		return resp, nil
	default:
		return nil, err
	}
}

func (m *defaultCustomerConfigModel) CountBuilder(field string) squirrel.SelectBuilder {
	return squirrel.Select("COUNT(" + field + ")").From(m.table).Where(squirrel.Eq{"is_deleted": 0})
}

func (m *defaultCustomerConfigModel) FindCount(ctx context.Context, countBuilder squirrel.SelectBuilder) (int64, error) {

	query, values, err := countBuilder.ToSql()
	if err != nil {
		return 0, err
	}

	var resp int64
	err = m.QueryRowNoCacheCtx(ctx, &resp, query, values...)
	switch err {
	case nil:
		return resp, nil
	default:
		return 0, err
	}
}

func (m *defaultCustomerConfigModel) BuildFiled(old, new *CustomerConfig) {

	if new == nil {
		return
	}
	if old == nil && new.Id == 0 {
		new.Id = util.GenerateSnowflakeInt64()
		return
	}
	new.Id = old.Id
	//
	//if new.ClearContextTime == 0 {
	//	new.ClearContextTime = old.ClearContextTime
	//}
	//
	//if new.TopK == 0 {
	//	new.TopK = old.TopK
	//}
	//
	//if new.KfId == "" {
	//	new.KfId = old.KfId
	//}
	//
	//if new.KfName == "" {
	//	new.KfName = old.KfName
	//}
	//
	//if new.Prompt == "" {
	//	new.Prompt = old.Prompt
	//}
	//
	//if new.PostModel == "" {
	//	new.PostModel = old.PostModel
	//}
	//
	//if new.EmbeddingMode == "" {
	//	new.EmbeddingMode = old.EmbeddingMode
	//}
	//
	//if !new.Score.Valid {
	//	new.Score = old.Score
	//}
	//
	//
	return
}
