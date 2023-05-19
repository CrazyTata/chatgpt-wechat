package model

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ ScriptLogModel = (*customScriptLogModel)(nil)

type (
	// ScriptLogModel is an interface to be customized, add more methods here,
	// and implement the added methods in customScriptLogModel.
	ScriptLogModel interface {
		scriptLogModel
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
