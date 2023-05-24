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
	url := "http://script:9997/api/crontab/run-script"
	_, _ = util.Get(url)
}
