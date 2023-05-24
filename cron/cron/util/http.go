package util

import (
	"bytes"
	"io"
	"net/http"
	"strings"
	"time"
)

const HTTP_TIME_OUT = 4

func Post(url string, params string, headerData map[string]string) (string, error) {

	client := &http.Client{Timeout: time.Duration(HTTP_TIME_OUT * time.Second)} //4秒超时
	req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(params)))
	if err != nil {
		return "", err
	}

	if len(headerData) > 0 {
		for k, v := range headerData {
			req.Header.Set(k, v)
		}
	}

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	if resp == nil {
		return "", nil
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(body), nil
}

func PostV2(url string, params string) (string, error) {
	resp, err := http.Post(url,
		"application/x-www-form-urlencoded",
		strings.NewReader(params))

	if err != nil {
		return "", err
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}

func Get(url string) (string, error) {
	client := &http.Client{}
	//提交请求
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		Info(" url:" + url + " http.NewRequest err:" + err.Error())
		return "", err
	}
	response, err := client.Do(request)
	if err != nil {
		Info(" url:" + url + " client.Do  err:" + err.Error())
		return "", err
	}
	defer response.Body.Close()
	rel, _ := io.ReadAll(response.Body)
	Info(" url:" + url + " response:" + string(rel))
	return string(rel), nil
}
