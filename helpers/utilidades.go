package helpers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"
	"strings"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/udistrital/tirilla_noticias_mid/models"
	"github.com/udistrital/utils_oas/request"
	"github.com/udistrital/utils_oas/requestresponse"
)

// Envia una petición con datos para cerar la tirilla de noticias
func SendRequestToCRUDAPI(endpoint string, data interface{}, method string) models.APIResponse {
	APICRUDURL := beego.AppConfig.String("router.APICRUD")

	rutaCompleta := APICRUDURL + endpoint
	logs.Info("Ruta completa: ", rutaCompleta)

	// Inicializar la respuesta
	var apiResp models.APIResponse

	// Convertir los datos a JSON
	reqBody, err := json.Marshal(data)
	if err != nil {
		apiResp.Err = err
		return models.APIResponse{Err: err}
	}

	// Configurar la solicitud HTTP
	req, err := http.NewRequest(method, rutaCompleta, bytes.NewBuffer(reqBody))
	if err != nil {
		apiResp.Err = err
		return models.APIResponse{Err: err}
	}
	req.Header.Set("Content-Type", "application/json")

	// Crear cliente HTTP y enviar la solicitud
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		apiResp.Err = err
		return models.APIResponse{Err: err}
	}
	defer resp.Body.Close()

	// Leer la respuesta de la API CRUD
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		apiResp.Err = err
		return models.APIResponse{Err: err}
	}

	// Asignar la respuesta al campo Body de la estructura APIResponse
	apiResp.Body = respBody

	return apiResp
}

func GetAllNoticias() (APIResponseDTO requestresponse.APIResponse) {

	var noticia []map[string]interface{}
	var listado []map[string]interface{}
	errNoticia := request.GetJson("http://"+beego.AppConfig.String("noticiaService")+fmt.Sprintf("/noticia?query=Activo:true&limit=0"), &noticia)
	if errNoticia == nil {
		fmt.Println("http://" + beego.AppConfig.String("noticiaService"))

		for _, noti := range noticia {
			var noticiaContenido = make(map[string]interface{}) // Mover la inicialización aquí
			noticiaContenido["activo"] = noti["Activo"]
			noticiaContenido["estilo"] = noti["IdEstilo"]
			noticiaContenido["prioridad"] = noti["Prioridad"]
			noticiaContenido["FechaInicio"] = noti["FechaInicio"]
			noticiaContenido["FechaFinal"] = noti["FechaFinal"]

			var responseNoticiaContenido []map[string]interface{}
			errNoticiaContenido := request.GetJson("http://"+beego.AppConfig.String("noticiaService")+fmt.Sprintf("/noticia_contenido?query=IdNoticia__id:%v", noti["Id"]), &responseNoticiaContenido)
			if errNoticiaContenido == nil {
				for _, conte := range responseNoticiaContenido {
					fmt.Println(reflect.TypeOf(conte["IdContenido"]))
					fmt.Println(conte["IdContenido"])
					dato := strings.ReplaceAll(conte["Dato"].(string), "{\"dato\": \"", "")
					dato = strings.TrimSuffix(dato, "\"}")
					if conte["IdContenido"] == float64(1) {
						fmt.Println("Entró")
						noticiaContenido["titulo"] = dato
					}
					if conte["IdContenido"] == float64(2) {
						noticiaContenido["descripcion"] = dato
					}
					if conte["IdContenido"] == float64(3) {
						noticiaContenido["link"] = dato
					}
				}
				listado = append(listado, noticiaContenido) // Mover la adición a listado aquí
			} else {
				APIResponseDTO = requestresponse.APIResponseDTO(false, 400, errNoticiaContenido.Error())
				return APIResponseDTO
			}
		}

		APIResponseDTO = requestresponse.APIResponseDTO(true, 200, listado)

	} else {
		fmt.Println(errNoticia.Error())
		APIResponseDTO = requestresponse.APIResponseDTO(false, 400, nil, errNoticia.Error())
	}
	return APIResponseDTO
}

func GetAllLista() (APIResponseDTO requestresponse.APIResponse) {
	var noticia []map[string]interface{}
	var listado []map[string]interface{}
	errNoticia := request.GetJson("http://"+beego.AppConfig.String("noticiaService")+fmt.Sprintf("/noticia?query=Activo:true&limit=0"), &noticia)
	if errNoticia == nil {
		fmt.Println("http://" + beego.AppConfig.String("noticiaService"))

		for _, noti := range noticia {
			var noticiaContenido = make(map[string]interface{}) // Mover la inicialización aquí
			noticiaContenido["activo"] = noti["Activo"]
			noticiaContenido["estilo"] = noti["IdEstilo"]
			noticiaContenido["prioridad"] = noti["Prioridad"]
			noticiaContenido["FechaInicio"] = noti["FechaInicio"]
			noticiaContenido["FechaFinal"] = noti["FechaFinal"]

			var responseNoticiaContenido []map[string]interface{}
			errNoticiaContenido := request.GetJson("http://"+beego.AppConfig.String("noticiaService")+fmt.Sprintf("/noticia_contenido?query=IdNoticia__id:%v", noti["Id"]), &responseNoticiaContenido)
			if errNoticiaContenido == nil {
				for _, conte := range responseNoticiaContenido {
					fmt.Println(reflect.TypeOf(conte["IdContenido"]))
					fmt.Println(conte["IdContenido"])
					dato := strings.ReplaceAll(conte["Dato"].(string), "{\"dato\": \"", "")
					dato = strings.TrimSuffix(dato, "\"}")
					if conte["IdContenido"] == float64(1) {
						fmt.Println("Entró")
						noticiaContenido["titulo"] = dato
					}
					if conte["IdContenido"] == float64(2) {
						noticiaContenido["descripcion"] = dato
					}
					if conte["IdContenido"] == float64(3) {
						noticiaContenido["link"] = dato
					}
				}
				var responseNoticiaEtiqueta []map[string]interface{}
				errNoticiaEtiqueta := request.GetJson("http://"+beego.AppConfig.String("noticiaService")+fmt.Sprintf("/noticia_etiqueta?query=IdNoticia__id:%v", noti["Id"]), &responseNoticiaEtiqueta)
				for _, eti := range responseNoticiaEtiqueta {
					fmt.Println(eti["IdEtiqueta"])
					noticiaContenido["IdEtiqueta"] = eti["IdEtiqueta"]
				}
				if errNoticiaEtiqueta == nil {
					listado = append(listado, noticiaContenido) // Mover la adición a listado aquí
				} else {
					APIResponseDTO = requestresponse.APIResponseDTO(false, 400, errNoticiaEtiqueta.Error())
					return APIResponseDTO
				}

			} else {
				APIResponseDTO = requestresponse.APIResponseDTO(false, 400, errNoticiaContenido.Error())
				return APIResponseDTO
			}
		}

		APIResponseDTO = requestresponse.APIResponseDTO(true, 200, listado)

	} else {
		fmt.Println(errNoticia.Error())
		APIResponseDTO = requestresponse.APIResponseDTO(false, 400, nil, errNoticia.Error())
	}
	return APIResponseDTO
}

func obtenerContenidoPorNoticia(noticiaID int) (models.Contenido, error) {
	var contenidoRespuesta models.ContenidoResponse
	var contenido models.Contenido

	// Construir la URL para obtener el contenido asociado a la noticia
	apiContenidoURL := fmt.Sprintf("%s/contenido/%d", beego.AppConfig.String("router.contenido"), noticiaID)
	contenidoResp := SendRequestToCRUDAPI(apiContenidoURL, nil, "GET")
	if contenidoResp.Err != nil {
		logs.Error("Error al obtener el contenido asociado a la noticia:", contenidoResp.Err)
		return contenido, contenidoResp.Err
	}

	// Decodificar la respuesta JSON en una estructura de contenido
	if err := json.Unmarshal(contenidoResp.Body, &contenidoRespuesta); err != nil {
		logs.Error("Error al decodificar la respuesta JSON de contenido:", err)
		return contenido, err
	}

	// Convertir la estructura de contenido a la estructura deseada por el cliente
	for _, item := range contenidoRespuesta.Data {
		var datoMap map[string]string
		err := json.Unmarshal([]byte(item.Dato), &datoMap)
		if err != nil {
			logs.Error("Error al decodificar el dato JSON:", err)
			return contenido, err
		}
		contenido.Id = append(contenido.Id, item.IdContenido)
		contenido.Dato = append(contenido.Dato, datoMap["dato"])
	}

	return contenido, nil
}
