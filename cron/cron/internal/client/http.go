package client

import (
	"cron/cron/internal/config"
	"cron/cron/util"
)

type HttpService struct {
}

func NewHttpService() *HttpService {
	return &HttpService{}
}

func (h *HttpService) RunScript() {
	url := config.ScriptService + config.RunScriptUri
	_, _ = util.Get(url)
}
