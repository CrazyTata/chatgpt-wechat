package model

import (
	"github.com/Masterminds/squirrel"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"golang.org/x/net/context"
)

var _ WechatUserModel = (*customWechatUserModel)(nil)

type (
	// WechatUserModel is an interface to be customized, add more methods here,
	// and implement the added methods in customWechatUserModel.
	WechatUserModel interface {
		wechatUserModel
		FindOneByQuery(ctx context.Context, rowBuilder squirrel.SelectBuilder) (*WechatUser, error)
		RowBuilder() squirrel.SelectBuilder
		FindAll(ctx context.Context, rowBuilder squirrel.SelectBuilder) ([]*WechatUser, error)
		FindCount(ctx context.Context, countBuilder squirrel.SelectBuilder) (int64, error)
	}

	customWechatUserModel struct {
		*defaultWechatUserModel
	}
)

// NewWechatUserModel returns a model for the database table.
func NewWechatUserModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) WechatUserModel {
	return &customWechatUserModel{
		defaultWechatUserModel: newWechatUserModel(conn, c, opts...),
	}
}

func (m *defaultWechatUserModel) FindOneByQuery(ctx context.Context, rowBuilder squirrel.SelectBuilder) (*WechatUser, error) {

	query, values, err := rowBuilder.ToSql()
	if err != nil {
		return nil, err
	}

	var resp WechatUser
	err = m.QueryRowNoCacheCtx(ctx, &resp, query, values...)
	switch err {
	case nil:
		return &resp, nil
	default:
		return nil, err
	}
}

// export logic
func (m *defaultWechatUserModel) RowBuilder() squirrel.SelectBuilder {
	return squirrel.Select(wechatUserRows).From(m.table)
}

func (m *defaultWechatUserModel) FindAll(ctx context.Context, rowBuilder squirrel.SelectBuilder) ([]*WechatUser, error) {

	query, values, err := rowBuilder.ToSql()
	if err != nil {
		return nil, err
	}

	var resp []*WechatUser
	err = m.QueryRowsNoCacheCtx(ctx, &resp, query, values...)
	switch err {
	case nil:
		return resp, nil
	default:
		return nil, err
	}
}

func (m *defaultWechatUserModel) FindCount(ctx context.Context, countBuilder squirrel.SelectBuilder) (int64, error) {

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
