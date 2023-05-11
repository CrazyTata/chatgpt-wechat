package util

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"time"
)

const TIMEOUT = 4

func Get(url string, headerData map[string]string) (string, error) {
	client := &http.Client{Timeout: time.Duration(TIMEOUT * time.Second)}
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}

	if len(headerData) > 0 {
		for k, v := range headerData {
			request.Header.Set(k, v)
		}
	}

	response, err := client.Do(request)
	if err != nil {
		return "", err
	}
	if response == nil {
		return "", nil
	}

	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}

func Post(url string, params []byte, headerData map[string]string) ([]byte, error) {
	client := &http.Client{Timeout: time.Duration(TIMEOUT * time.Second)} //4秒超时
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(params))
	if err != nil {
		return nil, err
	}

	if len(headerData) > 0 {
		for k, v := range headerData {
			req.Header.Set(k, v)
		}
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	if resp == nil {
		return nil, nil
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}
