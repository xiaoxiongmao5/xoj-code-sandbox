/*
 * @Author: 小熊 627516430@qq.com
 * @Date: 2023-10-03 19:11:56
 * @LastEditors: 小熊 627516430@qq.com
 * @LastEditTime: 2023-10-04 18:04:10
 * @FilePath: /xoj-code-sandbox/routers/router.go
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
// @APIVersion 1.0.0
// @Title beego Test API
// @Description beego has a very cool tools to autogenerate documents for your API
// @Contact astaxie@gmail.com
// @TermsOfServiceUrl http://beego.me/
// @License Apache 2.0
// @LicenseUrl http://www.apache.org/licenses/LICENSE-2.0.html
package routers

import (
	beego "github.com/beego/beego/v2/server/web"
	"github.com/xiaoxiongmao5/xoj/xoj-code-sandbox/controllers"
)

func init() {

	beego.CtrlGet("health", controllers.CheckController.Health)
	beego.CtrlGet("do", controllers.CheckController.Do)
	beego.CtrlGet("dockerdo", controllers.DockerController.Do)
	beego.CtrlGet("check", controllers.CheckController.CheckExec)
}
