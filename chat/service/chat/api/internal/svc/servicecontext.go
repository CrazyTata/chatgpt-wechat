package svc

import (
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"github.com/zeromicro/go-zero/rest"

	"chat/service/chat/api/internal/config"
	"chat/service/chat/api/internal/middleware"
	"chat/service/chat/model"
)

type ServiceContext struct {
	Config                 config.Config
	UserModel              model.UserModel
	ChatModel              model.ChatModel
	ChatConfigModel        model.ChatConfigModel
	CustomerPromptModel    model.CustomerPromptModel
	PromptConfigModel      model.PromptConfigModel
	ApplicationConfigModel model.ApplicationConfigModel
	CustomerConfigModel    model.CustomerConfigModel
	WechatUserModel        model.WechatUserModel
	AccessLog              rest.Middleware
}

func NewServiceContext(c config.Config) *ServiceContext {
	conn := sqlx.NewMysql(c.Mysql.DataSource)

	return &ServiceContext{
		Config:                 c,
		UserModel:              model.NewUserModel(conn, c.RedisCache),
		ChatModel:              model.NewChatModel(conn, c.RedisCache),
		ChatConfigModel:        model.NewChatConfigModel(conn, c.RedisCache),
		PromptConfigModel:      model.NewPromptConfigModel(conn, c.RedisCache),
		CustomerPromptModel:    model.NewCustomerPromptModel(conn, c.RedisCache),
		ApplicationConfigModel: model.NewApplicationConfigModel(conn, c.RedisCache),
		CustomerConfigModel:    model.NewCustomerConfigModel(conn, c.RedisCache),
		WechatUserModel:        model.NewWechatUserModel(conn, c.RedisCache),
		AccessLog:              middleware.NewAccessLogMiddleware().Handle,
	}
}
