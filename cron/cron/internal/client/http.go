package client

import (
	"cron/cron/util"
)

type HttpService struct {
}

func NewHttpService() *HttpService {
	return &HttpService{}
}

func (h *HttpService) RunScript() {
	url := "http://127.0.0.1:9997/api/crontab/run-script"

	_, _ = util.Get(url)
}
