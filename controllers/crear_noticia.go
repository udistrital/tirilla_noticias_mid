package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/udistrital/tirilla_noticias_mid/helpers"
	"github.com/udistrital/tirilla_noticias_mid/models"
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
	c.Mapping("GetAllLista", c.GetAllLista)
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

	apiResp := helpers.SendRequestToCRUDAPI(apiNoticiaURL, noticia, "POST")
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
			IdEtiqueta: etiquetaID,
		}

		// Enviar la solicitud POST a la API CRUD para la etiqueta
		apiResp := helpers.SendRequestToCRUDAPI(apiEtiquetaURL, etiquetaData, "POST")
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
			IdContenido: contenido.Id[i],
		}

		// Enviar la solicitud POST a la API CRUD para el contenido
		apiResp := helpers.SendRequestToCRUDAPI(apiContenidoURL, contenidoData, "POST")
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
		apiResp := helpers.SendRequestToCRUDAPI(apiModuloURL, moduloPublicaionData, "POST")
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

// GetAllLista ...
// @Title GetAllLista
// @Description get all Noticias with Etiquetas and Contenido
// @Success 200 {object} interface{}
// @Failure 403
// @router /lista [get]
func (c *Crear_noticiaController) GetAllLista() {
	defer errorhandler.HandlePanic(&c.Controller)
	respuesta := helpers.GetAllLista()

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

	//apiNoticiaURL := beego.AppConfig.String("router.noticia")
	apiEtiquetaURL := beego.AppConfig.String("router.etiqueta")
	apiContenidoURL := beego.AppConfig.String("router.contenido")
	// apiModuloURL := beego.AppConfig.String("router.modulo")

	etiquetas_desactivar_id := []int{}
	etiquetas_desactivar_fk := []int{}

	contenido_desactivar_id := []int{}
	contenido_desactivar_fk := []int{}
	//contenido_desactivar_id_tipo_contenido := []int{}

	// Obtener el ID de la noticia a actualizar del parámetro de la URL
	idStr := c.Ctx.Input.Param(":id")
	id, err := strconv.Atoi(idStr)

	if err != nil {
		logs.Error("Error al convertir el ID de la noticia:", err)
		c.CustomAbort(http.StatusBadRequest, "ID de la noticia inválido")
		return
	}

	logs.Info("ID de la noticia:", id)

	// Decodificar el cuerpo de la solicitud en un objeto NoticiaRequest
	var reqBody models.NoticiaRequest
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &reqBody); err != nil {
		logs.Error("Error al decodificar el cuerpo de la solicitud:", err)
		c.CustomAbort(http.StatusBadRequest, "Error al decodificar el cuerpo de la solicitud")
		return
	}

	logs.Info("Cuerpo de la solicitud:", reqBody)
	noticia := reqBody.Noticia
	etiqueta := reqBody.Etiqueta
	contenido := reqBody.Contenido
	moduloPublicaion := reqBody.ModuloPublicacion

	logs.Info("Noticiaaaaaaa:", noticia)
	logs.Info("Modulo de publicación:", moduloPublicaion)

	// Actualizar la noticia en la tabla Noticia
	apiNoticiaURL := fmt.Sprintf("%s/%d", beego.AppConfig.String("router.noticia"), id)
	apiResp := helpers.SendRequestToCRUDAPI(apiNoticiaURL, noticia, "PUT")
	if apiResp.Err != nil {
		logs.Error("Error al enviar la solicitud a la API CRUD para actualizar la noticia:", apiResp.Err)
		c.CustomAbort(http.StatusInternalServerError, "Error al enviar la solicitud a la API CRUD para actualizar la noticia")
		return
	}

	// Ahora, después de actualizar la noticia, vamos a obtener todas las etiquetas asociadas a esta noticia
	apiEtiquetasURL := fmt.Sprintf("%s/etiquetas/%d", beego.AppConfig.String("router.etiqueta"), id)
	etiquetasResp := helpers.SendRequestToCRUDAPI(apiEtiquetasURL, nil, "GET")
	if etiquetasResp.Err != nil {
		logs.Error("Error al obtener las etiquetas asociadas a la noticia:", etiquetasResp.Err)
		c.CustomAbort(http.StatusInternalServerError, "Error al obtener las etiquetas asociadas a la noticia")
		return
	}

	// Decodificar la respuesta JSON en la estructura definida
	var etiquetaRespuesta models.EtiquetaResponse
	if err := json.Unmarshal(etiquetasResp.Body, &etiquetaRespuesta); err != nil {
		logs.Error("Error al decodificar la respuesta JSON de etiquetas:", err)
		c.CustomAbort(http.StatusInternalServerError, "Error al decodificar la respuesta JSON de etiquetas")
		return
	}

	//############################################################################################################# Etiqueta

	// Crear nuevas etiquetas para las etiquetas en la solicitud PUT que no estén en la respuesta de la API
	for _, etiquetaPUT := range etiqueta.IdTipoEtiqueta {
		etiquetaEncontrada := false

		// Iterar sobre las etiquetas de la respuesta de la API
		for _, etiquetaAPI := range etiquetaRespuesta.Data {
			// Si la etiqueta de la solicitud PUT coincide con una etiqueta de la respuesta de la API
			if etiquetaAPI.IdEtiqueta == etiquetaPUT {
				etiquetaEncontrada = true

				// Si la etiqueta de la API está desactivada, activarla
				if !etiquetaAPI.Activo {
					logs.Info("Activando la etiqueta:", etiquetaAPI)
					activarEtiqueta := models.EtiquetaData{
						Activo: true,
						IdNoticia: struct {
							Id int `json:"Id"`
						}{
							Id: id, // Id de la noticia
						},
						IdEtiqueta: etiquetaAPI.IdEtiqueta,
					}

					// Enviar la solicitud para activar la etiqueta
					apiResp := helpers.SendRequestToCRUDAPI(apiEtiquetaURL+"/"+strconv.Itoa(etiquetaAPI.Id), activarEtiqueta, "PUT")
					if apiResp.Err != nil {
						logs.Error("Error al enviar la solicitud a la API CRUD para activar la etiqueta:", apiResp.Err)
						c.CustomAbort(http.StatusInternalServerError, "Error al enviar la solicitud a la API CRUD para activar la etiqueta")
						return
					}
					break
				}
			}
		}

		// Si la etiqueta de la solicitud PUT no está en la respuesta de la API, crear una nueva etiqueta
		if !etiquetaEncontrada {
			nuevaEtiqueta := models.EtiquetaData{
				Activo: true,
				IdNoticia: struct {
					Id int `json:"Id"`
				}{
					Id: id, // Id de la noticia
				},
				IdEtiqueta: etiquetaPUT,
			}

			// Enviar la solicitud para crear la nueva etiqueta
			apiResp := helpers.SendRequestToCRUDAPI(apiEtiquetaURL, nuevaEtiqueta, "POST")
			if apiResp.Err != nil {
				logs.Error("Error al enviar la solicitud a la API CRUD para la etiqueta:", apiResp.Err)
				c.CustomAbort(http.StatusInternalServerError, "Error al enviar la solicitud a la API CRUD para la etiqueta")
				return
			}
		}
	}

	// Iterar sobre las etiquetas de la respuesta de la API
	for _, etiquetaAPI := range etiquetaRespuesta.Data {
		etiquetaEncontrada := false

		// Iterar sobre las etiquetas de la solicitud PUT
		for _, etiquetaPUT := range etiqueta.IdTipoEtiqueta {
			// Si la etiqueta de la solicitud PUT coincide con una etiqueta de la respuesta de la API
			if etiquetaAPI.IdEtiqueta == etiquetaPUT {
				etiquetaEncontrada = true
				break
			}
		}

		// Si la etiqueta de la API no está en la solicitud PUT, desactivarla
		if !etiquetaEncontrada {
			etiquetaAPI.Activo = false
			etiquetas_desactivar_id = append(etiquetas_desactivar_id, etiquetaAPI.Id)
			etiquetas_desactivar_fk = append(etiquetas_desactivar_fk, etiquetaAPI.IdEtiqueta)
		}
	}

	// Iterar sobre las etiquetas para desactivar
	for i, idEtiqueta := range etiquetas_desactivar_id {
		// Obtener el IdTipoEtiqueta correspondiente del arreglo etiquetas_desactivar_fk
		idTipoEtiqueta := etiquetas_desactivar_fk[i]

		// Crear una estructura EtiquetaData con el campo Activo establecido en false
		etiquetaDesactivar := models.EtiquetaData{
			Activo: false,
			IdNoticia: struct {
				Id int `json:"Id"`
			}{
				Id: id, // ID de la noticia
			},
			IdEtiqueta: idTipoEtiqueta, // ID del tipo de etiqueta

		}

		// Enviar la solicitud para desactivar la etiqueta
		apiResp := helpers.SendRequestToCRUDAPI(apiEtiquetaURL+"/"+strconv.Itoa(idEtiqueta), etiquetaDesactivar, "PUT")
		if apiResp.Err != nil {
			logs.Error("Error al enviar la solicitud a la API CRUD para desactivar la etiqueta:", apiResp.Err)
			c.CustomAbort(http.StatusInternalServerError, "Error al enviar la solicitud a la API CRUD para desactivar la etiqueta")
			return
		}
	}

	//############################################################################################################# Contenido

	// Ahora, después de actualizar la etiqueta, vamos a obtener todo el contenido asociado a esta noticia
	apiContenidosURL := fmt.Sprintf("%s/contenido/%d", beego.AppConfig.String("router.contenido"), id)
	contenidoResp := helpers.SendRequestToCRUDAPI(apiContenidosURL, nil, "GET")
	if contenidoResp.Err != nil {
		logs.Error("Error al obtener las etiquetas asociadas a la noticia:", contenidoResp.Err)
		c.CustomAbort(http.StatusInternalServerError, "Error al obtener las etiquetas asociadas a la noticia")
		return
	}

	// Decodificar la respuesta JSON en la estructura definida
	var contenidoRespuesta models.ContenidoResponse
	if err := json.Unmarshal(contenidoResp.Body, &contenidoRespuesta); err != nil {
		logs.Error("Error al decodificar la respuesta JSON de contenido:", err)
		c.CustomAbort(http.StatusInternalServerError, "Error al decodificar la respuesta JSON de contenido")
		return
	}

	logs.Info("Respuesta decodificada de la API CRUD para contenido:", contenidoRespuesta)

	// Crear nuevas etiquetas para las etiquetas en la solicitud PUT que no estén en la respuesta de la API
	for i, contenidoPUT := range contenido.Id {
		ContenidoEncontrado := false

		// Iterar sobre las etiquetas de la respuesta de la API
		for _, contenidoAPI := range contenidoRespuesta.Data {
			// Si el contenido de la solicitud PUT coincide con un contenido de la respuesta de la API
			if contenidoAPI.IdContenido == contenidoPUT {
				ContenidoEncontrado = true

				// Si el contenido de la API está desactivado, activarlo
				if !contenidoAPI.Activo {
					datoJSON := fmt.Sprintf(`{"dato": "%s"}`, contenido.Dato[i])
					// logs.Info("Activando la etiqueta:", contenidoAPI)
					//logs.Info("Contenido:", contenido.Dato[contenidoPUT])
					activarContenido := models.ContenidoData{
						Activo: true,
						Dato:   datoJSON,
						IdNoticia: struct {
							Id int `json:"Id"`
						}{
							Id: id, // Id de la noticia
						},
						IdContenido: contenidoAPI.IdContenido,
					}

					logs.Info("Activando contenido:", activarContenido)

					// Enviar la solicitud para activar el contenido
					apiResp := helpers.SendRequestToCRUDAPI(apiContenidoURL+"/"+strconv.Itoa(contenidoAPI.Id), activarContenido, "PUT")
					if apiResp.Err != nil {
						logs.Error("Error al enviar la solicitud a la API CRUD para activar el contenido:", apiResp.Err)
						c.CustomAbort(http.StatusInternalServerError, "Error al enviar la solicitud a la API CRUD para activar el contenido")
						return
					}
					break
				} else {
					datoJSON := fmt.Sprintf(`{"dato": "%s"}`, contenido.Dato[i])
					// logs.Info("Activando la etiqueta:", contenidoAPI)
					//logs.Info("Contenido:", contenido.Dato[contenidoPUT])
					activarContenido := models.ContenidoData{
						Activo: true,
						Dato:   datoJSON,
						IdNoticia: struct {
							Id int `json:"Id"`
						}{
							Id: id, // Id de la noticia
						},
						IdContenido: contenidoAPI.IdContenido,
					}

					logs.Info("Activando contenido:", activarContenido)

					// Enviar la solicitud para activar el contenido
					apiResp := helpers.SendRequestToCRUDAPI(apiContenidoURL+"/"+strconv.Itoa(contenidoAPI.Id), activarContenido, "PUT")
					if apiResp.Err != nil {
						logs.Error("Error al enviar la solicitud a la API CRUD para activar el contenido:", apiResp.Err)
						c.CustomAbort(http.StatusInternalServerError, "Error al enviar la solicitud a la API CRUD para activar el contenido")
						return
					}
					break
				}
			}
		}

		// Si el contenido de la solicitud PUT no está en la respuesta de la API, crear un nuevo contenido
		if !ContenidoEncontrado {
			datoJSON := fmt.Sprintf(`{"dato": "%s"}`, contenido.Dato[contenidoPUT])
			nuevoContenido := models.ContenidoData{
				Activo: true,
				Dato:   datoJSON,
				IdNoticia: struct {
					Id int `json:"Id"`
				}{
					Id: id, // Id de la noticia
				},
				IdContenido: contenidoPUT,
			}

			logs.Info("Creando nuevo contenido:", nuevoContenido)

			// Enviar la solicitud para crear la nueva etiqueta
			apiResp := helpers.SendRequestToCRUDAPI(apiContenidoURL, nuevoContenido, "POST")
			if apiResp.Err != nil {
				logs.Error("Error al enviar la solicitud a la API CRUD para la etiqueta:", apiResp.Err)
				c.CustomAbort(http.StatusInternalServerError, "Error al enviar la solicitud a la API CRUD para la etiqueta")
				return
			}
		}
	}

	// Iterar sobre las etiquetas de la respuesta de la API
	for _, contenidoAPI := range contenidoRespuesta.Data {
		contenidoEncontrado := false

		// Iterar sobre las etiquetas de la solicitud PUT
		for _, contenidoPUT := range contenido.Id {
			// Si la etiqueta de la solicitud PUT coincide con una etiqueta de la respuesta de la API
			if contenidoAPI.IdContenido == contenidoPUT {
				contenidoEncontrado = true
				break
			}
		}

		// Si la etiqueta de la API no está en la solicitud PUT, desactivarla
		if !contenidoEncontrado {
			contenidoAPI.Activo = false
			contenido_desactivar_id = append(contenido_desactivar_id, contenidoAPI.Id)
			contenido_desactivar_fk = append(contenido_desactivar_fk, contenidoAPI.IdContenido)
			//contenido_desactivar_id_tipo_contenido = append(contenido_desactivar_id_tipo_contenido, contenidoAPI.IdContenido.Id)
		}
	}

	// Iterar sobre las etiquetas para desactivar
	for i, idContenido := range contenido_desactivar_id {
		// Obtener el IdTipoEtiqueta correspondiente del arreglo etiquetas_desactivar_fk
		IdContenido := contenido_desactivar_fk[i]

		logs.Info("IdContenido:", IdContenido)
		logs.Info("valor a guardar en dato:", contenidoRespuesta.Data[IdContenido-1].Dato)

		// Crear una estructura EtiquetaData con el campo Activo establecido en false
		contenidoDesactivar := models.ContenidoData{
			Activo: false,
			Dato:   contenidoRespuesta.Data[IdContenido-1].Dato,
			IdNoticia: struct {
				Id int `json:"Id"`
			}{
				Id: id, // ID de la noticia
			},
			IdContenido: IdContenido, // ID del tipo de etiqueta

		}

		// Enviar la solicitud para desactivar la etiqueta
		apiResp := helpers.SendRequestToCRUDAPI(apiContenidoURL+"/"+strconv.Itoa(idContenido), contenidoDesactivar, "PUT")
		if apiResp.Err != nil {
			logs.Error("Error al enviar la solicitud a la API CRUD para desactivar la etiqueta:", apiResp.Err)
			c.CustomAbort(http.StatusInternalServerError, "Error al enviar la solicitud a la API CRUD para desactivar la etiqueta")
			return
		}
	}

	// Respondiendo al cliente Angular
	c.Ctx.Output.SetStatus(http.StatusOK)
	c.Data["json"] = map[string]string{"message": "Noticia actualizada exitosamente"}
	c.ServeJSON()

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
