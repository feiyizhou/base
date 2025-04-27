package utils

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// HttpPostJson 以json格式发送post请求
func HttpPostJson(url string, headers map[string]string, body any, timeout ...time.Duration) ([]byte, error) {
	byteArr, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("parame json formatting exception, err: %s", err.Error())
	}
	request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(byteArr))
	if err != nil {
		return nil, fmt.Errorf("new http request exception, err: %s", err.Error())
	}
	for k, v := range headers {
		request.Header.Set(k, v)
	}
	return do(request, timeout...)
}

// HttpGet send get type request
func HttpGet(url string, headers map[string]string, params map[string]any, timeout ...time.Duration) ([]byte, error) {
	if len(params) > 0 {
		var paramArr []string
		for k, v := range params {
			paramArr = append(paramArr, fmt.Sprintf("%s=%v", k, v))
		}
		url = fmt.Sprintf("%s?%s", url, strings.Join(paramArr, "&"))
	}
	request, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("new http request exception, err: %s", err.Error())
	}
	for k, v := range headers {
		request.Header.Set(k, v)
	}
	return do(request, timeout...)
}

// do send request
func do(request *http.Request, timeout ...time.Duration) ([]byte, error) {
	client := &http.Client{}
	if len(timeout) == 0 {
		client.Timeout = 30 * time.Second
	} else {
		client.Timeout = timeout[0]
	}
	response, err := client.Do(request)
	if err != nil {
		return nil, fmt.Errorf("response exception, err: %s", err.Error())
	}
	if response.StatusCode != http.StatusOK {
		return nil, errors.New("request exception")
	}
	defer response.Body.Close()
	data, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("read response body exception, err: %s", err.Error())
	}
	return data, err
}
