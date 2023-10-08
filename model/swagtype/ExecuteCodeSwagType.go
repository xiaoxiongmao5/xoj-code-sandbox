/*
 * @Author: 小熊 627516430@qq.com
 * @Date: 2023-10-08 14:15:20
 * @LastEditors: 小熊 627516430@qq.com
 * @LastEditTime: 2023-10-08 14:29:38
 */
package swagtype

import "github.com/xiaoxiongmao5/xoj/xoj-code-sandbox/model"

type ExecuteCodeResponse struct {
	Code    int                       `json:"code"`
	Message string                    `json:"message"`
	Data    model.ExecuteCodeResponse `json:"data"`
}
