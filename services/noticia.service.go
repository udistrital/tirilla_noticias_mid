package services

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/astaxie/beego"
	"github.com/udistrital/tirilla_noticias_mid/tirilla_noticias_mid/models"
	"github.com/udistrital/utils_oas/request"
	"github.com/udistrital/utils_oas/requestresponse"
	"github.com/udistrital/utils_oas/time_bogota"
)

func PostNoticia(data []byte) (APIResponseDTO requestresponse.APIResponse) {
	//Almacena la nueva noticia
	var nuevaNoticia map[string]interface{}
	var nuevoContenido map[string]interface{}
	var nuevaEtiqueta map[string]interface{}
	var errSaveAll bool
	//respuesta a la petición
	var respuesta map[string]interface{}
	//timestamp
	date := time_bogota.TiempoBogotaFormato()

	if err := json.Unmarshal(data, &nuevaNoticia); err == nil {

		dataNoticia := map[string]interface{}{
			"Activo":            true,
			"FechaCreacion":     date,
			"FechaModificacion": date,
			"Prioridad":         nuevaNoticia["Prioridad"].(float64),
			"IdEstilo":          nuevaNoticia["IdEstilo"].(float64),
			"FechaInicio":       fmt.Sprintf("%v", nuevaNoticia["FechaInicio"]),
			"FechaFinal":        fmt.Sprintf("%v", nuevaNoticia["FechaFinal"]),
		}
		//var guardada map[string]interface{}
		nuevaNoticia = dataNoticia

		errNoticia := request.SendJson("http://"+beego.AppConfig.String("noticiaService")+"/noticia/", "POST", &nuevaNoticia, dataNoticia)
		if errNoticia == nil {

			//noticiaId := nuevaNoticia["Data"].(map[string]interface{})["Id"].(float64)
			if err := json.Unmarshal(data, &nuevaNoticia); err == nil {
				//fmt.Println(nuevaNoticia["Contenido"])
				contenido, contenidoExist := nuevaNoticia["Contenido"]
				noti := nuevaNoticia["Data"]
				if contenidoExist {

					contenidos := contenido.([]interface{})
					//fmt.Println(len(contenidos))
					for _, c := range contenidos {
						contenidoMap := c.(map[string]interface{})
						dato := contenidoMap["Dato"].(string)
						contenidoId := contenidoMap["IdContenido"].(float64)
						dataContenido := map[string]interface{}{
							"Dato":              dato,
							"Activo":            true,
							"FechaCreacion":     date,
							"FechaModificacion": date,
							"IdNoticia":         noti,
							"IdContenido":       contenidoId,
						}

						errContenido := request.SendJson("http://"+beego.AppConfig.String("noticiaService")+"/noticia_contenido/", "POST", &nuevoContenido, dataContenido)
						if errContenido != nil {
							errSaveAll = true
						}
					}
					if !errSaveAll {
						if err := json.Unmarshal(data, &nuevaNoticia); err == nil {
							//fmt.Println(nuevaNoticia["Etiqueta"])
							etiqueta, etiquetaExist := nuevaNoticia["Etiqueta"]
							noti := nuevaNoticia["Data"]
							if etiquetaExist {

								etiquetas := etiqueta.([]interface{})
								//fmt.Println(len(etiquetas))
								for _, c := range etiquetas {
									etiquetaMap := c.(map[string]interface{})
									etiquetaId := etiquetaMap["IdEtiqueta"].(float64)
									dataEtiquetas := map[string]interface{}{
										"Activo":            true,
										"FechaCreacion":     date,
										"FechaModificacion": date,
										"IdNoticia":         noti,
										"IdEtiqueta":        etiquetaId,
									}

									errEtiqueta := request.SendJson("http://"+beego.AppConfig.String("noticiaService")+"/noticia_etiqueta/", "POST", &nuevaEtiqueta, dataEtiquetas)
									if errEtiqueta != nil {
										errSaveAll = true
									}
								}
								if !errSaveAll {

									APIResponseDTO = requestresponse.APIResponseDTO(true, 200, nuevaEtiqueta)
									return APIResponseDTO
								} else {
									models.SetInactivo(fmt.Sprintf("http://"+beego.AppConfig.String("noticiaService")+"/noticia_etiqueta/%.f", nuevaNoticia["Data"].(map[string]interface{})["Id"].(float64)))
								}
							}
						}
						APIResponseDTO = requestresponse.APIResponseDTO(true, 200, nuevoContenido)
						return APIResponseDTO
					} else {
						models.SetInactivo(fmt.Sprintf("http://"+beego.AppConfig.String("noticiaService")+"/noticia_contenido/%.f", nuevaNoticia["Data"].(map[string]interface{})["Id"].(float64)))
					}
				}
			}

			APIResponseDTO = requestresponse.APIResponseDTO(true, 200, nuevaNoticia)
			return APIResponseDTO
		} else {
			models.SetInactivo(fmt.Sprintf("http://"+beego.AppConfig.String("noticiaService")+"/noticia/%.f", nuevaNoticia["Data"].(map[string]interface{})["Id"].(float64)))
		}

		APIResponseDTO = requestresponse.APIResponseDTO(true, 500, respuesta, nuevaNoticia)
		return APIResponseDTO
	}

	APIResponseDTO = requestresponse.APIResponseDTO(true, 200, respuesta, nuevaNoticia)
	return APIResponseDTO
}
func GetAllNoticias() (APIResponseDTO requestresponse.APIResponse) {
	fmt.Println("GetAll")
	var noticia []map[string]interface{}
	var listado []map[string]interface{}
	errNoticia := request.GetJson("http://"+beego.AppConfig.String("noticiaService")+fmt.Sprintf("/noticia?query=Activo:true&limit=0"), &noticia)
	if errNoticia == nil {
		//fmt.Println("http://" + beego.AppConfig.String("noticiaService"))

		for _, noti := range noticia {
			var noticiaContenido = make(map[string]interface{}) // Mover la inicialización aquí
			noticiaContenido["activo"] = noti["Activo"]
			noticiaContenido["estilo"] = noti["IdEstilo"]
			noticiaContenido["prioridad"] = noti["Prioridad"]
			noticiaContenido["fechaInicio"] = noti["FechaInicio"]
			noticiaContenido["fechaFinal"] = noti["FechaFinal"]
			noticiaContenido["id"] = noti["Id"]

			var responseNoticiaContenido []map[string]interface{}
			errNoticiaContenido := request.GetJson("http://"+beego.AppConfig.String("noticiaService")+fmt.Sprintf("/noticia_contenido?query=IdNoticia__id:%v", noti["Id"]), &responseNoticiaContenido)
			if errNoticiaContenido == nil {
				//fmt.Println(noti["Id"])
				//fmt.Println(responseNoticiaContenido)
				for _, conte := range responseNoticiaContenido {
					//fmt.Println(reflect.TypeOf(conte["IdContenido"]))
					//fmt.Println(conte["IdContenido"])
					dato := strings.ReplaceAll(conte["Dato"].(string), "{\"dato\": \"", "")
					dato = strings.TrimSuffix(dato, "\"}")

					if conte["IdContenido"] == float64(1) {
						//fmt.Println("Entró")
						noticiaContenido["idTitulo"] = conte["Id"]
						noticiaContenido["titulo"] = dato
					}
					if conte["IdContenido"] == float64(2) {
						noticiaContenido["idDesc"] = conte["Id"]
						noticiaContenido["descripcion"] = dato
					}
					if conte["IdContenido"] == float64(3) {
						noticiaContenido["idLink"] = conte["Id"]
						noticiaContenido["link"] = dato
					}
				}
				//listado = append(listado, noticiaContenido)
			} else {
				APIResponseDTO = requestresponse.APIResponseDTO(false, 400, errNoticiaContenido.Error())
				return APIResponseDTO
			}
			var responseNoticiaEtiqueta []map[string]interface{}
			errNoticiaEtiqueta := request.GetJson("http://"+beego.AppConfig.String("noticiaService")+fmt.Sprintf("/noticia_etiqueta?query=IdNoticia__id:%v", noti["Id"]), &responseNoticiaEtiqueta)
			if errNoticiaEtiqueta == nil {
				//fmt.Println(noti["Id"])
				//fmt.Println(responseNoticiaEtiqueta)
				for _, eti := range responseNoticiaEtiqueta {
					noticiaContenido["idEtiqueta"] = eti["Id"]
					noticiaContenido["etiqueta"] = eti["IdEtiqueta"]
				}
				listado = append(listado, noticiaContenido)
			} else {
				APIResponseDTO = requestresponse.APIResponseDTO(false, 400, errNoticiaEtiqueta.Error())
				return APIResponseDTO
			}
		}

		APIResponseDTO = requestresponse.APIResponseDTO(true, 200, listado)

	} else {
		//fmt.Println(errNoticia.Error())
		APIResponseDTO = requestresponse.APIResponseDTO(false, 400, nil, errNoticia.Error())
	}
	return APIResponseDTO
}

func PutNoticia(data []byte) (APIResponseDTO requestresponse.APIResponse) {
	//var body map[string]interface{}
	var errSaveAll bool
	var noticiaActualizada map[string]interface{}
	var contenidoActualizado map[string]interface{}
	var date = time_bogota.TiempoBogotaFormato()

	if err := json.Unmarshal(data, &noticiaActualizada); err == nil {

		dataActualizada := map[string]interface{}{
			"Activo":            noticiaActualizada["Activo"],
			"FechaCreacion":     noticiaActualizada["FechaCreacion"],
			"FechaModificacion": date,
			"FechaFinal":        noticiaActualizada["FechaFinal"],
			"FechaInicio":       noticiaActualizada["FechaInicio"],
			"Id":                noticiaActualizada["Id"],
			"IdEstilo":          noticiaActualizada["IdEstilo"],
			"Prioridad":         noticiaActualizada["Prioridad"],
		}

		noticiaActualizada = dataActualizada

		idNoticia := noticiaActualizada["Id"].(float64)
		errActualizarNoticia := request.SendJson("http://"+beego.AppConfig.String("noticiaService")+"/noticia/"+fmt.Sprintf("%.f", idNoticia), "PUT", &noticiaActualizada, dataActualizada)
		if errActualizarNoticia == nil {
			if err := json.Unmarshal(data, &noticiaActualizada); err == nil {
				//fmt.Println(nuevaNoticia["Contenido"])
				contenido, contenidoExist := noticiaActualizada["Contenido"]
				noti := noticiaActualizada["Data"]
				if contenidoExist {
					contenidos := contenido.([]interface{})
					//fmt.Println(len(contenidos))
					for _, c := range contenidos {
						contenidoMap := c.(map[string]interface{})
						dato := contenidoMap["Dato"].(string)
						contenidoId := contenidoMap["IdContenido"].(float64)
						idConte := contenidoMap["Id"].(float64)
						dataContenido := map[string]interface{}{
							"Dato":              dato,
							"Activo":            true,
							"FechaCreacion":     noticiaActualizada["FechaCreacion"],
							"FechaModificacion": date,
							"IdNoticia":         noti,
							"IdContenido":       contenidoId,
							"Id":                idConte,
						}
						contenidoActualizado = dataContenido

						errContenido := request.SendJson("http://"+beego.AppConfig.String("noticiaService")+"/noticia_contenido/"+fmt.Sprintf("%.f", idConte), "PUT", &contenidoActualizado, dataContenido)
						if errContenido != nil {
							errSaveAll = true
						}
					}
					if !errSaveAll {
						if err := json.Unmarshal(data, &noticiaActualizada); err == nil {
							//fmt.Println(nuevaNoticia["Etiqueta"])
							etiqueta, etiquetaExist := noticiaActualizada["Etiqueta"]
							noti := noticiaActualizada["Data"]
							if etiquetaExist {

								etiquetas := etiqueta.([]interface{})
								//fmt.Println(len(etiquetas))
								for _, c := range etiquetas {
									etiquetaMap := c.(map[string]interface{})
									etiquetaId := etiquetaMap["IdEtiqueta"].(float64)
									idEtiqueta := etiquetaMap["Id"].(float64)
									dataEtiquetas := map[string]interface{}{
										"Activo":            true,
										"FechaCreacion":     noticiaActualizada["FechaCreacion"],
										"FechaModificacion": date,
										"IdNoticia":         noti,
										"IdEtiqueta":        etiquetaId,
										"Id":                idEtiqueta,
									}

									errContenido := request.SendJson("http://"+beego.AppConfig.String("noticiaService")+"/noticia_etiqueta/"+fmt.Sprintf("%.f", idEtiqueta), "PUT", &contenidoActualizado, dataEtiquetas)
									if errContenido != nil {
										errSaveAll = true
									}
								}
								if !errSaveAll {

									APIResponseDTO = requestresponse.APIResponseDTO(true, 200, noticiaActualizada)
									return APIResponseDTO
								}
							}
						}
						APIResponseDTO = requestresponse.APIResponseDTO(true, 200, noticiaActualizada)
						return APIResponseDTO
					}
				}
			}
		}
		return requestresponse.APIResponseDTO(false, 200, noticiaActualizada, dataActualizada)

	} else {
		return requestresponse.APIResponseDTO(false, 400, nil, "Error al decodificar datos JSON")
	}

}
