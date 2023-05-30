package request

import (
	"github.com/flipped-aurora/gin-vue-admin/server/model/chatAdmin"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
	"time"
)

type CustomerConfigSearch struct {
	chatAdmin.CustomerConfig
	StartCreatedAt *time.Time `json:"startCreatedAt" form:"startCreatedAt"`
	EndCreatedAt   *time.Time `json:"endCreatedAt" form:"endCreatedAt"`
	request.PageInfo
}
