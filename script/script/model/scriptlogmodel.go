package model

import (
	"context"
	"github.com/Masterminds/squirrel"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ ScriptLogModel = (*customScriptLogModel)(nil)

const (
	ScriptLogStatusStart   = 1
	ScriptLogStatusRunning = 2
	ScriptLogStatusSuccess = 3
	ScriptLogStatusFail    = 4
)

type (
	// ScriptLogModel is an interface to be customized, add more methods here,
	// and implement the added methods in customScriptLogModel.
	ScriptLogModel interface {
		scriptLogModel
		FindOneByQuery(ctx context.Context, rowBuilder squirrel.SelectBuilder) (*ScriptLog, error)
		RowBuilder() squirrel.SelectBuilder
		FindAll(ctx context.Context, rowBuilder squirrel.SelectBuilder) ([]*ScriptLog, error)
		FindCount(ctx context.Context, countBuilder squirrel.SelectBuilder) (int64, error)
		CountBuilder(field string) squirrel.SelectBuilder
	}

	customScriptLogModel struct {
		*defaultScriptLogModel
	}
)

// NewScriptLogModel returns a model for the database table.
func NewScriptLogModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) ScriptLogModel {
	return &customScriptLogModel{
		defaultScriptLogModel: newScriptLogModel(conn, c, opts...),
	}
}

func (m *defaultScriptLogModel) FindOneByQuery(ctx context.Context, rowBuilder squirrel.SelectBuilder) (*ScriptLog, error) {

	query, values, err := rowBuilder.ToSql()
	if err != nil {
		return nil, err
	}

	var resp ScriptLog
	err = m.QueryRowNoCacheCtx(ctx, &resp, query, values...)
	if err != nil {
		if err == sqlx.ErrNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &resp, nil
}

func (m *defaultScriptLogModel) RowBuilder() squirrel.SelectBuilder {
	return squirrel.Select(scriptLogRows).From(m.table)
}

func (m *defaultScriptLogModel) FindAll(ctx context.Context, rowBuilder squirrel.SelectBuilder) ([]*ScriptLog, error) {

	query, values, err := rowBuilder.ToSql()
	if err != nil {
		return nil, err
	}

	var resp []*ScriptLog
	err = m.QueryRowsNoCacheCtx(ctx, &resp, query, values...)
	switch err {
	case nil:
		return resp, nil
	default:
		return nil, err
	}
}

func (m *defaultScriptLogModel) FindCount(ctx context.Context, countBuilder squirrel.SelectBuilder) (int64, error) {

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

func (m *defaultScriptLogModel) CountBuilder(field string) squirrel.SelectBuilder {
	return squirrel.Select("COUNT(" + field + ")").From(m.table)
}
