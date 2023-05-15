package model

import (
	"context"
	"github.com/Masterminds/squirrel"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ CustomerPromptModel = (*customCustomerPromptModel)(nil)

type (
	// CustomerPromptModel is an interface to be customized, add more methods here,
	// and implement the added methods in customCustomerPromptModel.
	CustomerPromptModel interface {
		customerPromptModel
		FindOneByQuery(ctx context.Context, rowBuilder squirrel.SelectBuilder) (*CustomerPrompt, error)
		RowBuilder() squirrel.SelectBuilder
		FindAll(ctx context.Context, rowBuilder squirrel.SelectBuilder) ([]*CustomerPrompt, error)
	}

	customCustomerPromptModel struct {
		*defaultCustomerPromptModel
	}
)

// NewCustomerPromptModel returns a model for the database table.
func NewCustomerPromptModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) CustomerPromptModel {
	return &customCustomerPromptModel{
		defaultCustomerPromptModel: newCustomerPromptModel(conn, c, opts...),
	}
}

func (m *defaultCustomerPromptModel) FindOneByQuery(ctx context.Context, rowBuilder squirrel.SelectBuilder) (*CustomerPrompt, error) {

	query, values, err := rowBuilder.ToSql()
	if err != nil {
		return nil, err
	}

	var resp CustomerPrompt
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
func (m *defaultCustomerPromptModel) RowBuilder() squirrel.SelectBuilder {
	return squirrel.Select(customerPromptRows).From(m.table)
}

func (m *defaultCustomerPromptModel) FindAll(ctx context.Context, rowBuilder squirrel.SelectBuilder) ([]*CustomerPrompt, error) {

	query, values, err := rowBuilder.ToSql()
	if err != nil {
		return nil, err
	}

	var resp []*CustomerPrompt
	err = m.QueryRowsNoCacheCtx(ctx, &resp, query, values...)
	switch err {
	case nil:
		return resp, nil
	default:
		return nil, err
	}
}
