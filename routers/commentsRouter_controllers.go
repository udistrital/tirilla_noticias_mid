package routers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context/param"
)

func init() {

    beego.GlobalControllerRouter["github.com/udistrital/tirilla_noticias_mid/controllers:Crear_noticiaController"] = append(beego.GlobalControllerRouter["github.com/udistrital/tirilla_noticias_mid/controllers:Crear_noticiaController"],
        beego.ControllerComments{
            Method: "PostNoticia",
            Router: "/",
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/tirilla_noticias_mid/controllers:Crear_noticiaController"] = append(beego.GlobalControllerRouter["github.com/udistrital/tirilla_noticias_mid/controllers:Crear_noticiaController"],
        beego.ControllerComments{
            Method: "GetAllNoticias",
            Router: "/",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/tirilla_noticias_mid/controllers:Crear_noticiaController"] = append(beego.GlobalControllerRouter["github.com/udistrital/tirilla_noticias_mid/controllers:Crear_noticiaController"],
        beego.ControllerComments{
            Method: "GetOne",
            Router: "/:id",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/tirilla_noticias_mid/controllers:Crear_noticiaController"] = append(beego.GlobalControllerRouter["github.com/udistrital/tirilla_noticias_mid/controllers:Crear_noticiaController"],
        beego.ControllerComments{
            Method: "PutNoticia",
            Router: "/:id",
            AllowHTTPMethods: []string{"put"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/tirilla_noticias_mid/controllers:Crear_noticiaController"] = append(beego.GlobalControllerRouter["github.com/udistrital/tirilla_noticias_mid/controllers:Crear_noticiaController"],
        beego.ControllerComments{
            Method: "Delete",
            Router: "/:id",
            AllowHTTPMethods: []string{"delete"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

}
