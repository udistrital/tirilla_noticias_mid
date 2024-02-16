// @APIVersion 1.0.0
// @Title beego Test API
// @Description beego has a very cool tools to autogenerate documents for your API
// @Contact astaxie@gmail.com
// @TermsOfServiceUrl http://beego.me/
// @License Apache 2.0
// @LicenseUrl http://www.apache.org/licenses/LICENSE-2.0.html
package routers

import (
	"github.com/udistrital/tirilla_noticias_mid/tirilla_noticias_mid/controllers"
	"github.com/astaxie/beego"
)

func init() {
    ns := beego.NewNamespace("/v1",
		beego.NSNamespace("/noticia-mid",
			beego.NSInclude(
				&controllers.Crear_noticiaController{},
			),
		),
	)
	beego.AddNamespace(ns)
}