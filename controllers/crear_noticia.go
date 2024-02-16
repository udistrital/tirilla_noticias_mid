package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/udistrital/tirilla_noticias_mid/tirilla_noticias_mid/helpers"
	"github.com/udistrital/tirilla_noticias_mid/tirilla_noticias_mid/models"
	"github.com/udistrital/utils_oas/errorhandler"
)

// Crear_noticiaController operations for Crear_noticia
type Crear_noticiaController struct {
	beego.Controller
}

// URLMapping ...
func (c *Crear_noticiaController) URLMapping() {
	c.Mapping("Post", c.Post)
	c.Mapping("GetOne", c.GetOne)
	c.Mapping("GetAll", c.GetAll)
	c.Mapping("Put", c.Put)
	c.Mapping("Delete", c.Delete)
}

// Post ...
// @Title Create
// @Description create Crear_noticia
// @Param	body		body 	models.Crear_noticia	true		"body for Crear_noticia content"
// @Success 201 {object} models.Crear_noticia
// @Failure 403 body is empty
// @router / [post]
func (c *Crear_noticiaController) Post() {

	apiNoticiaURL := beego.AppConfig.String("router.noticia")
	apiEtiquetaURL := beego.AppConfig.String("router.etiqueta")
	apiContenidoURL := beego.AppConfig.String("router.contenido")
	apiModuloURL := beego.AppConfig.String("router.modulo")

	var reqBody models.NoticiaRequest
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &reqBody); err != nil {
		logs.Error("Error al decodificar el cuerpo de la solicitud:", err)
		c.CustomAbort(http.StatusBadRequest, "Error al decodificar el cuerpo de la solicitud")
		return
	}

	noticia := reqBody.Noticia
	etiqueta := reqBody.Etiqueta
	contenido := reqBody.Contenido
	moduloPublicaion := reqBody.ModuloPublicacion

	apiResp := helpers.SendRequestToCRUDAPI(apiNoticiaURL, noticia)
	if apiResp.Err != nil {
		logs.Error("Error al enviar la solicitud a la API CRUD para Noticia:", apiResp.Err)
		c.CustomAbort(http.StatusInternalServerError, "Error al enviar la solicitud a la API CRUD para Noticia")
		return
	}

	var noticiaResp models.NoticiaResponse
	if err := json.Unmarshal(apiResp.Body, &noticiaResp); err != nil {
		logs.Error("Error al decodificar la respuesta de la API CRUD para Noticia:", err)
		c.CustomAbort(http.StatusInternalServerError, "Error al decodificar la respuesta de la API CRUD para Noticia")
		return
	}
	noticiaID := noticiaResp.Data.ID

	// Hacer solicitudes POST a la API CRUD para Etiqueta
	for _, etiquetaID := range etiqueta.IdTipoEtiqueta {
		etiquetaData := models.EtiquetaData{
			Activo: true,
			IdNoticia: struct {
				Id int `json:"Id"`
			}{
				Id: noticiaID,
			},
			IdTipoEtiqueta: struct {
				Id int `json:"Id"`
			}{
				Id: etiquetaID,
			},
		}

		// Enviar la solicitud POST a la API CRUD para la etiqueta
		apiResp := helpers.SendRequestToCRUDAPI(apiEtiquetaURL, etiquetaData)
		if apiResp.Err != nil {
			logs.Error("Error al enviar la solicitud a la API CRUD para la etiqueta:", apiResp.Err)
			c.CustomAbort(http.StatusInternalServerError, "Error al enviar la solicitud a la API CRUD para la etiqueta")
			return
		}
	}

	// Hacer solicitudes POST a la API CRUD para Contenido
	for i, _ := range contenido.Id {
		datoJSON := fmt.Sprintf(`{"dato": "%s"}`, contenido.Dato[i])
		contenidoData := models.ContenidoData{
			Activo: true,
			Dato:   datoJSON,
			IdNoticia: struct {
				Id int `json:"Id"`
			}{
				Id: noticiaID,
			},
			IdTipoContenido: struct {
				Id int `json:"Id"`
			}{
				Id: contenido.Id[i],
			},
		}

		// Enviar la solicitud POST a la API CRUD para el contenido
		apiResp := helpers.SendRequestToCRUDAPI(apiContenidoURL, contenidoData)
		if apiResp.Err != nil {
			logs.Error("Error al enviar la solicitud a la API CRUD para el contenido:", apiResp.Err)
			c.CustomAbort(http.StatusInternalServerError, "Error al enviar la solicitud a la API CRUD para el contenido")
			return
		}
	}

	// Hacer solicitudes POST a la API CRUD para ModuloPublicaion
	for _, moduloID := range moduloPublicaion.IdModulo {
		moduloPublicaionData := models.ModuloPublicacionData{
			Activo: true,
			IdNoticia: struct {
				Id int `json:"Id"`
			}{
				Id: noticiaID,
			},
			RefModuloId: moduloID,
		}

		// Enviar la solicitud POST a la API CRUD para el moduloPublicaion
		apiResp := helpers.SendRequestToCRUDAPI(apiModuloURL, moduloPublicaionData)
		if apiResp.Err != nil {
			logs.Error("Error al enviar la solicitud a la API CRUD para el moduloPublicaion:", apiResp.Err)
			c.CustomAbort(http.StatusInternalServerError, "Error al enviar la solicitud a la API CRUD para el moduloPublicaion")
			return
		}
	}

	// Respondiendo al cliente Angular
	c.Ctx.Output.SetStatus(http.StatusCreated)
	c.Data["json"] = map[string]string{"message": "Noticia creada exitosamente"}
	c.ServeJSON()

}

// GetOne ...
// @Title GetOne
// @Description get Crear_noticia by id
// @Param	id		path 	string	true		"The key for staticblock"
// @Success 200 {object} models.Crear_noticia
// @Failure 403 :id is empty
// @router /:id [get]
func (c *Crear_noticiaController) GetOne() {

}

// GetAll ...
// @Title GetAll
// @Description get Crear_noticia
// @Param	query	query	string	false	"Filter. e.g. col1:v1,col2:v2 ..."
// @Param	fields	query	string	false	"Fields returned. e.g. col1,col2 ..."
// @Param	sortby	query	string	false	"Sorted-by fields. e.g. col1,col2 ..."
// @Param	order	query	string	false	"Order corresponding to each sortby field, if single value, apply to all sortby fields. e.g. desc,asc ..."
// @Param	limit	query	string	false	"Limit the size of result set. Must be an integer"
// @Param	offset	query	string	false	"Start position of result set. Must be an integer"
// @Success 200 {object} models.Crear_noticia
// @Failure 403
// @router / [get]
func (c *Crear_noticiaController) GetAll() {
	defer errorhandler.HandlePanic(&c.Controller)
	respuesta := helpers.GetAllNoticias()

	c.Ctx.Output.SetStatus(respuesta.Status)
	c.Data["json"] = respuesta
	c.ServeJSON()
}

// Put ...
// @Title Put
// @Description update the Crear_noticia
// @Param	id		path 	string	true		"The id you want to update"
// @Param	body		body 	models.Crear_noticia	true		"body for Crear_noticia content"
// @Success 200 {object} models.Crear_noticia
// @Failure 403 :id is not int
// @router /:id [put]
func (c *Crear_noticiaController) Put() {

}

// Delete ...
// @Title Delete
// @Description delete the Crear_noticia
// @Param	id		path 	string	true		"The id you want to delete"
// @Success 200 {string} delete success!
// @Failure 403 id is empty
// @router /:id [delete]
func (c *Crear_noticiaController) Delete() {

}
