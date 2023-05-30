package service

import (
	"github.com/flipped-aurora/gin-vue-admin/server/service/chatAdmin"
	"github.com/flipped-aurora/gin-vue-admin/server/service/client"
	"github.com/flipped-aurora/gin-vue-admin/server/service/example"
	"github.com/flipped-aurora/gin-vue-admin/server/service/system"
)

type ServiceGroup struct {
	SystemServiceGroup    system.ServiceGroup
	ExampleServiceGroup   example.ServiceGroup
	ChatAdminServiceGroup chatAdmin.ServiceGroup
	ClientGroup           client.ClientGroup
}

var ServiceGroupApp = new(ServiceGroup)
