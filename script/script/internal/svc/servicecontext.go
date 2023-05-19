package svc

import (
	"github.com/zeromicro/go-zero/core/stores/sqlx"

	"github.com/zeromicro/go-zero/rest"
	"script/script/internal/config"
	"script/script/internal/middleware"
	"script/script/model"
)

type ServiceContext struct {
	Config         config.Config
	AccessLog      rest.Middleware
	ScriptModel    model.ScriptModel
	ScriptLogModel model.ScriptLogModel
}

func NewServiceContext(c config.Config) *ServiceContext {
	conn := sqlx.NewMysql(c.Mysql.DataSource)
	return &ServiceContext{
		Config:         c,
		AccessLog:      middleware.NewAccessLogMiddleware().Handle,
		ScriptModel:    model.NewScriptModel(conn, c.RedisCache),
		ScriptLogModel: model.NewScriptLogModel(conn, c.RedisCache),
	}
}
