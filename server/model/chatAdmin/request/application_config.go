package request

import (
	"github.com/flipped-aurora/gin-vue-admin/server/model/chatAdmin"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
	"time"
)

type ApplicationConfigSearch struct{
    chatAdmin.ApplicationConfig
    StartCreatedAt *time.Time `json:"startCreatedAt" form:"startCreatedAt"`
    EndCreatedAt   *time.Time `json:"endCreatedAt" form:"endCreatedAt"`
    request.PageInfo
}
