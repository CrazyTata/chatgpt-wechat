package model

import (
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
	return squirrel.Select(customerConfigRows).From(m.table)
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
