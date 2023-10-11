/*
 * @Author: 小熊 627516430@qq.com
 * @Date: 2023-10-09 15:40:33
 * @LastEditors: 小熊 627516430@qq.com
 * @LastEditTime: 2023-10-09 16:11:42
 */
package utils

import (
	"bytes"
	"io"
	"net/http"
)

// 发送 HTTP 请求函数
func SendHTTPRequest(method, targetURL string, requestBody []byte, opts ...func(*http.Request)) (bodyBytes []byte, err error) {
	// 创建 HTTP 请求
	req, err := http.NewRequest(method, targetURL, bytes.NewBuffer(requestBody))
	if err != nil {
		return bodyBytes, err
	}

	// 应用选项函数
	for _, opt := range opts {
		opt(req)
	}

	// 发送 HTTP 请求
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return bodyBytes, err
	}
	defer resp.Body.Close()

	// 读取响应 Body 内容
	bodyBytes, err = io.ReadAll(resp.Body)
	if err != nil {
		return bodyBytes, err
	}

	return bodyBytes, nil
}

// 设置请求 Header 选项函数
func WithHeaders(headers map[string]string) func(*http.Request) {
	return func(req *http.Request) {
		for key, value := range headers {
			req.Header.Set(key, value)
		}
	}
}

// 设置请求 Cookie 选项函数
func WithCookies(cookies []*http.Cookie) func(*http.Request) {
	return func(req *http.Request) {
		for _, cookie := range cookies {
			req.AddCookie(cookie)
		}
	}
}
