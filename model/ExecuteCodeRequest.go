/*
 * @Author: 小熊 627516430@qq.com
 * @Date: 2023-10-02 12:24:44
 * @LastEditors: 小熊 627516430@qq.com
 * @LastEditTime: 2023-10-04 11:17:40
 */
package model

type ExecuteCodeRequest struct {
	InputList []string `json:"inputList"`
	Code      string   `json:"code"`
	Language  string   `json:"language"`
}
