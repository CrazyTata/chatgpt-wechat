package model

import (
	"context"
	"github.com/Masterminds/squirrel"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ ScriptModel = (*customScriptModel)(nil)

const (
	ScriptDeleted    = 1
	ScriptNotDeleted = 0

	ScriptEnable  = 1
	ScriptDisable = 0
)

type (
	// ScriptModel is an interface to be customized, add more methods here,
	// and implement the added methods in customScriptModel.
	ScriptModel interface {
		scriptModel
		FindOneByQuery(ctx context.Context, rowBuilder squirrel.SelectBuilder) (*Script, error)
		RowBuilder() squirrel.SelectBuilder
		FindAll(ctx context.Context, rowBuilder squirrel.SelectBuilder) ([]*Script, error)
		FindCount(ctx context.Context, countBuilder squirrel.SelectBuilder) (int64, error)
		CountBuilder(field string) squirrel.SelectBuilder
	}

	customScriptModel struct {
		*defaultScriptModel
	}
)

// NewScriptModel returns a model for the database table.
func NewScriptModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) ScriptModel {
	return &customScriptModel{
		defaultScriptModel: newScriptModel(conn, c, opts...),
	}
}

func (m *defaultScriptModel) FindOneByQuery(ctx context.Context, rowBuilder squirrel.SelectBuilder) (*Script, error) {

	query, values, err := rowBuilder.ToSql()
	if err != nil {
		return nil, err
	}

	var resp Script
	err = m.QueryRowNoCacheCtx(ctx, &resp, query, values...)
	if err != nil {
		if err == sqlx.ErrNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &resp, nil
}

func (m *defaultScriptModel) RowBuilder() squirrel.SelectBuilder {
	return squirrel.Select(scriptRows).From(m.table)
}

func (m *defaultScriptModel) FindAll(ctx context.Context, rowBuilder squirrel.SelectBuilder) ([]*Script, error) {

	query, values, err := rowBuilder.ToSql()
	if err != nil {
		return nil, err
	}

	var resp []*Script
	err = m.QueryRowsNoCacheCtx(ctx, &resp, query, values...)
	switch err {
	case nil:
		return resp, nil
	default:
		return nil, err
	}
}

func (m *defaultScriptModel) FindCount(ctx context.Context, countBuilder squirrel.SelectBuilder) (int64, error) {

	query, values, err := countBuilder.ToSql()
	if err != nil {
		return 0, err
	}

	var resp int64
	err = m.QueryRowNoCacheCtx(ctx, &resp, query, values...)
	if err != nil {
		if err == sqlx.ErrNotFound {
			return 0, nil
		}
		return 0, err
	}
	return resp, nil
}

func (m *defaultScriptModel) CountBuilder(field string) squirrel.SelectBuilder {
	return squirrel.Select("COUNT(" + field + ")").From(m.table)
}
