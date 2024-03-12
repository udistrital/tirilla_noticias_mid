package controllers

import (
	"github.com/astaxie/beego"
	"github.com/udistrital/tirilla_noticias_mid/tirilla_noticias_mid/services"
	"github.com/udistrital/utils_oas/errorhandler"
)

// NoticiaController operations for Noticia
type NoticiaController struct {
	beego.Controller
}

// URLMapping ...
func (c *NoticiaController) URLMapping() {
	c.Mapping("PostNoticia", c.PostNoticia)
	c.Mapping("GetOne", c.GetOne)
	c.Mapping("GetAllNoticias", c.GetAllNoticias)
	c.Mapping("PutNoticia", c.PutNoticia)
	c.Mapping("Delete", c.Delete)
}

// PostNoticia ...
// @Title PostNoticia
// @Description create Noticia
// @Param	body		body 	models.Noticia	true		"body for Noticia content"
// @Success 201 {object} models.Noticia
// @Failure 403 body is empty
// @router / [post]
func (c *NoticiaController) PostNoticia() {
	defer errorhandler.HandlePanic(&c.Controller)
	data := c.Ctx.Input.RequestBody
	respuesta := services.PostNoticia(data)
	c.Ctx.Output.SetStatus(respuesta.Status)
	c.Data["json"] = respuesta
	c.ServeJSON()

}

// GetOne ...
// @Title GetOne
// @Description get Noticia by id
// @Param	id		path 	string	true		"The key for staticblock"
// @Success 200 {object} models.Noticia
// @Failure 403 :id is empty
// @router /:id [get]
func (c *NoticiaController) GetOne() {

}

// GetAllNoticias ...
// @Title GetAllNoticias
// @Description get Noticia
// @Param	query	query	string	false	"Filter. e.g. col1:v1,col2:v2 ..."
// @Param	fields	query	string	false	"Fields returned. e.g. col1,col2 ..."
// @Param	sortby	query	string	false	"Sorted-by fields. e.g. col1,col2 ..."
// @Param	order	query	string	false	"Order corresponding to each sortby field, if single value, apply to all sortby fields. e.g. desc,asc ..."
// @Param	limit	query	string	false	"Limit the size of result set. Must be an integer"
// @Param	offset	query	string	false	"Start position of result set. Must be an integer"
// @Success 200 {object} models.Noticia
// @Failure 403
// @router / [get]
func (c *NoticiaController) GetAllNoticias() {
	defer errorhandler.HandlePanic(&c.Controller)
	respuesta := services.GetAllNoticias()
	c.Ctx.Output.SetStatus(respuesta.Status)
	c.Data["json"] = respuesta
	c.ServeJSON()
}

// PutNoticia ...
// @Title PutNoticia
// @Description update the Noticia
// @Param	id		path 	string	true		"The id you want to update"
// @Param	body		body 	models.Noticia	true		"body for Noticia content"
// @Success 200 {object} models.Noticia
// @Failure 403 :id is not int
// @router /:id [put]
func (c *NoticiaController) PutNoticia() {
	defer errorhandler.HandlePanic(&c.Controller)
	data := c.Ctx.Input.RequestBody
	respuesta := services.PutNoticia(data)
	c.Ctx.Output.SetStatus(respuesta.Status)
	c.Data["json"] = respuesta
	c.ServeJSON()
}

// Delete ...
// @Title Delete
// @Description delete the Noticia
// @Param	id		path 	string	true		"The id you want to delete"
// @Success 200 {string} delete success!
// @Failure 403 id is empty
// @router /:id [delete]
func (c *NoticiaController) Delete() {

}
