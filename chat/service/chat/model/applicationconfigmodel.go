package model

import (
	"context"
	"github.com/Masterminds/squirrel"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ ApplicationConfigModel = (*customApplicationConfigModel)(nil)

type (
	// ApplicationConfigModel is an interface to be customized, add more methods here,
	// and implement the added methods in customApplicationConfigModel.
	ApplicationConfigModel interface {
		applicationConfigModel
		FindOneByQuery(ctx context.Context, rowBuilder squirrel.SelectBuilder) (*ApplicationConfig, error)
		RowBuilder() squirrel.SelectBuilder
		FindAll(ctx context.Context, rowBuilder squirrel.SelectBuilder) ([]*ApplicationConfig, error)
	}

	customApplicationConfigModel struct {
		*defaultApplicationConfigModel
	}
)

// NewApplicationConfigModel returns a model for the database table.
func NewApplicationConfigModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) ApplicationConfigModel {
	return &customApplicationConfigModel{
		defaultApplicationConfigModel: newApplicationConfigModel(conn, c, opts...),
	}
}

func (m *defaultApplicationConfigModel) FindOneByQuery(ctx context.Context, rowBuilder squirrel.SelectBuilder) (*ApplicationConfig, error) {

	query, values, err := rowBuilder.ToSql()
	if err != nil {
		return nil, err
	}

	var resp ApplicationConfig
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
func (m *defaultApplicationConfigModel) RowBuilder() squirrel.SelectBuilder {
	return squirrel.Select(applicationConfigRows).From(m.table)
}

func (m *defaultApplicationConfigModel) FindAll(ctx context.Context, rowBuilder squirrel.SelectBuilder) ([]*ApplicationConfig, error) {

	query, values, err := rowBuilder.ToSql()
	if err != nil {
		return nil, err
	}

	var resp []*ApplicationConfig
	err = m.QueryRowsNoCacheCtx(ctx, &resp, query, values...)
	switch err {
	case nil:
		return resp, nil
	default:
		return nil, err
	}
}
