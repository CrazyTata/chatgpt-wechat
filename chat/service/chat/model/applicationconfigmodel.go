package model

import (
	"chat/common/util"
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
		CountBuilder(field string) squirrel.SelectBuilder
		FindCount(ctx context.Context, countBuilder squirrel.SelectBuilder) (int64, error)
		BuildFiled(old, new *ApplicationConfig)
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
func (m *defaultApplicationConfigModel) CountBuilder(field string) squirrel.SelectBuilder {
	return squirrel.Select("COUNT(" + field + ")").From(m.table)
}

func (m *defaultApplicationConfigModel) FindCount(ctx context.Context, countBuilder squirrel.SelectBuilder) (int64, error) {

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

func (m *defaultApplicationConfigModel) BuildFiled(old, new *ApplicationConfig) {
	if new == nil {
		return
	}
	if old == nil && new.Id == 0 {
		new.Id = util.GenerateSnowflakeInt64()
		return
	}
	new.Id = old.Id
	if new.AgentId == 0 {
		new.AgentId = old.AgentId
	}

	if new.ClearContextTime == 0 {
		new.ClearContextTime = old.ClearContextTime
	}

	if new.TopK == 0 {
		new.TopK = old.TopK
	}

	new.GroupEnable = old.GroupEnable
	new.EmbeddingEnable = old.EmbeddingEnable

	if new.AgentSecret == "" {
		new.AgentSecret = old.AgentSecret
	}

	if new.GroupChatId == "" {
		new.GroupChatId = old.GroupChatId
	}

	if new.GroupName == "" {
		new.GroupName = old.GroupName
	}

	if new.Welcome == "" {
		new.Welcome = old.Welcome
	}

	if new.BasePrompt == "" {
		new.BasePrompt = old.BasePrompt
	}

	if new.PostModel == "" {
		new.PostModel = old.PostModel
	}

	if new.Model == "" {
		new.Model = old.Model
	}

	if new.AgentName == "" {
		new.AgentName = old.AgentName
	}

	if new.EmbeddingMode == "" {
		new.EmbeddingMode = old.EmbeddingMode
	}

	if !new.Score.Valid {
		new.Score = old.Score
	}
	return
}
