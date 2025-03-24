package utils

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// HttpPostJson 以json格式发送post请求
func HttpPostJson(ctx context.Context, url string, headers map[string]string, body interface{}, timeout ...time.Duration) ([]byte, error) {
	if body == nil {
		return nil, fmt.Errorf("the body of post request must not be null ")
	}
	byteArr, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("参数json格式化异常，异常原因：%v，context：%v", err, ctx)
	}
	request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(byteArr))
	if err != nil {
		return nil, fmt.Errorf("新建请求体异常，异常原因：%v，context：%v", err, ctx)
	}
	for k, v := range headers {
		request.Header.Set(k, v)
	}
	return do(ctx, request, timeout...)
}

// HttpGet 发送get请求
func HttpGet(ctx context.Context, url string, headers map[string]string, params map[string]interface{}, timeout ...time.Duration) ([]byte, error) {
	var paramArr []string
	for k, v := range params {
		paramArr = append(paramArr, fmt.Sprintf("%s=%v", k, v))
	}
	paramStr := strings.Join(paramArr, "&")
	if len(paramStr) > 0 {
		url = fmt.Sprintf("%s?%s", url, paramStr)
	}
	request, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("新建请求体异常，异常原因：%v，context：%v", err, ctx)
	}
	for k, v := range headers {
		request.Header.Set(k, v)
	}
	return do(ctx, request, timeout...)
}

// do 发送请求
func do(ctx context.Context, request *http.Request, timeout ...time.Duration) ([]byte, error) {
	client := &http.Client{}
	if len(timeout) == 0 {
		client.Timeout = 30 * time.Second
	} else {
		client.Timeout = timeout[0]
	}
	response, err := client.Do(request)
	if err != nil {
		return nil, fmt.Errorf("响应异常，异常原因：%v，context：%v", err, ctx)
	}
	if response.StatusCode != http.StatusOK {
		return nil, errors.New("请求异常")
	}
	defer response.Body.Close()
	data, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应数据，异常原因：%v，context：%v", err, ctx)
	}
	return data, err
}
