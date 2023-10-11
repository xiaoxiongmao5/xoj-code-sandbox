/*
 * @Author: 小熊 627516430@qq.com
 * @Date: 2023-10-09 15:58:28
 * @LastEditors: 小熊 627516430@qq.com
 * @LastEditTime: 2023-10-09 16:14:10
 */
package utils

import (
	"encoding/json"
	"net/http"
	"testing"
)

// 响应的数据结构
type ResponseData struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    struct {
		OutputList []string `json:"OutputList"`
		Message    string   `json:"message"`
		Status     int      `json:"status"`
		JudgeInfo  struct {
			Message string `json:"message"`
			Memory  int    `json:"memory"`
			Time    int    `json:"time"`
		} `json:"judgeInfo"`
	} `json:"data"`
}

func TestSendHTTPRequest(t *testing.T) {
	method := "POST"
	targetURL := "http://127.0.0.1:8093/executeCode"
	requestBody := []byte(`
	{
		"inputList":["1 2","3 4","5 6","7 8", "9 10"],
		"code":"package main\n\nimport (\n\t\"fmt\"\n\t\"os\"\n\t\"strconv\"\n)\n\nfunc main() {\n\t// 使用 os.Args 来获取命令行参数，os.Args[0] 是程序名称，os.Args[1:] 包含所有的命令行参数\n\targs := os.Args[1:]\n\ta, _ := strconv.Atoi(args[0])\n\tb, _ := strconv.Atoi(args[1])\n\tfmt.Printf(\"%d\", a+b)\n}",
		"language":"go"
	}`)

	// 创建请求选项
	headers := map[string]string{
		"HeaderKey1": "HeaderValue1",
		"HeaderKey2": "HeaderValue2",
	}
	cookies := []*http.Cookie{{
		Name:  "CookieName",
		Value: "CookieValue",
	}}

	bodyBytes, err := SendHTTPRequest(
		method,
		targetURL,
		requestBody,
		WithHeaders(headers),
		WithCookies(cookies),
	)
	if err != nil {
		t.Error("发送 HTTP 请求失败：", err)
		return
	}

	// 解析 JSON
	var responseData ResponseData
	if err := json.Unmarshal(bodyBytes, &responseData); err != nil {
		t.Error("解析 JSON 失败：", err)
		return
	}

	if responseData.Code != 0 {
		t.Error("http响应结果code!=0: ", responseData.Code)
		return
	}

	t.Log("测试正常, 响应内容：", string(bodyBytes))
}
