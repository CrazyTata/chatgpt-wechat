// Code generated by goctl. DO NOT EDIT.

package model

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/zeromicro/go-zero/core/stores/builder"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlc"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"github.com/zeromicro/go-zero/core/stringx"
)

var (
	customerConfigFieldNames          = builder.RawFieldNames(&CustomerConfig{})
	customerConfigRows                = strings.Join(customerConfigFieldNames, ",")
	customerConfigRowsExpectAutoSet   = strings.Join(stringx.Remove(customerConfigFieldNames, "`id`", "`create_at`", "`create_time`", "`created_at`", "`update_at`", "`update_time`", "`updated_at`"), ",")
	customerConfigRowsWithPlaceHolder = strings.Join(stringx.Remove(customerConfigFieldNames, "`id`", "`create_at`", "`create_time`", "`created_at`", "`update_at`", "`update_time`", "`updated_at`"), "=?,") + "=?"

	cacheCustomerConfigIdPrefix = "cache:customerConfig:id:"
)

type (
	customerConfigModel interface {
		Insert(ctx context.Context, data *CustomerConfig) (sql.Result, error)
		FindOne(ctx context.Context, id int64) (*CustomerConfig, error)
		Update(ctx context.Context, data *CustomerConfig) error
		Delete(ctx context.Context, id int64) error
	}

	defaultCustomerConfigModel struct {
		sqlc.CachedConn
		table string
	}

	CustomerConfig struct {
		Id               int64           `db:"id"`
		KfId             string          `db:"kf_id"` // 客服ID
		KfName           string          `db:"kf_name"`
		Prompt           string          `db:"prompt"`
		PostModel        string          `db:"post_model"`         // 发送请求的model
		EmbeddingEnable  bool            `db:"embedding_enable"`   // 是否启用embedding
		EmbeddingMode    string          `db:"embedding_mode"`     // embedding的搜索模式
		Score            sql.NullFloat64 `db:"score"`              // 分数
		TopK             int64           `db:"top_k"`              // topK
		ClearContextTime int64           `db:"clear_context_time"` // 需要清理上下文的时间，按分配置，默认0不清理
		CreatedAt        time.Time       `db:"created_at"`         // 创建时间
		UpdatedAt        time.Time       `db:"updated_at"`         // 更新时间
		IsDeleted        int64           `db:"is_deleted"`         // 是否删除
	}
)

func newCustomerConfigModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) *defaultCustomerConfigModel {
	return &defaultCustomerConfigModel{
		CachedConn: sqlc.NewConn(conn, c, opts...),
		table:      "`customer_config`",
	}
}

func (m *defaultCustomerConfigModel) Delete(ctx context.Context, id int64) error {
	customerConfigIdKey := fmt.Sprintf("%s%v", cacheCustomerConfigIdPrefix, id)
	_, err := m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("update %s set is_deleted=1 where `id` = ?", m.table)
		return conn.ExecCtx(ctx, query, id)
	}, customerConfigIdKey)
	return err
}

func (m *defaultCustomerConfigModel) FindOne(ctx context.Context, id int64) (*CustomerConfig, error) {
	customerConfigIdKey := fmt.Sprintf("%s%v", cacheCustomerConfigIdPrefix, id)
	var resp CustomerConfig
	err := m.QueryRowCtx(ctx, &resp, customerConfigIdKey, func(ctx context.Context, conn sqlx.SqlConn, v any) error {
		query := fmt.Sprintf("select %s from %s where `id` = ? limit 1", customerConfigRows, m.table)
		return conn.QueryRowCtx(ctx, v, query, id)
	})
	switch err {
	case nil:
		return &resp, nil
	case sqlc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *defaultCustomerConfigModel) Insert(ctx context.Context, data *CustomerConfig) (sql.Result, error) {
	customerConfigIdKey := fmt.Sprintf("%s%v", cacheCustomerConfigIdPrefix, data.Id)
	ret, err := m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("insert into %s (%s) values (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)", m.table, customerConfigRowsExpectAutoSet)
		return conn.ExecCtx(ctx, query, data.KfId, data.KfName, data.Prompt, data.PostModel, data.EmbeddingEnable, data.EmbeddingMode, data.Score, data.TopK, data.ClearContextTime, data.IsDeleted)
	}, customerConfigIdKey)
	return ret, err
}

func (m *defaultCustomerConfigModel) Update(ctx context.Context, data *CustomerConfig) error {
	customerConfigIdKey := fmt.Sprintf("%s%v", cacheCustomerConfigIdPrefix, data.Id)
	_, err := m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("update %s set %s where `id` = ?", m.table, customerConfigRowsWithPlaceHolder)
		return conn.ExecCtx(ctx, query, data.KfId, data.KfName, data.Prompt, data.PostModel, data.EmbeddingEnable, data.EmbeddingMode, data.Score, data.TopK, data.ClearContextTime, data.IsDeleted, data.Id)
	}, customerConfigIdKey)
	return err
}

func (m *defaultCustomerConfigModel) formatPrimary(primary any) string {
	return fmt.Sprintf("%s%v", cacheCustomerConfigIdPrefix, primary)
}

func (m *defaultCustomerConfigModel) queryPrimary(ctx context.Context, conn sqlx.SqlConn, v, primary any) error {
	query := fmt.Sprintf("select %s from %s where `id` = ? limit 1", customerConfigRows, m.table)
	return conn.QueryRowCtx(ctx, v, query, primary)
}

func (m *defaultCustomerConfigModel) tableName() string {
	return m.table
}
