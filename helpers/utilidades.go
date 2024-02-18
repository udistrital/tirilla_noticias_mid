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
	"github.com/udistrital/tirilla_noticias_mid/tirilla_noticias_mid/models"
	"github.com/udistrital/utils_oas/request"
	"github.com/udistrital/utils_oas/requestresponse"
)

// Envia una petición con datos para cerar la tirilla de noticias
func SendRequestToCRUDAPI(endpoint string, data interface{}) models.APIResponse {
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
	req, err := http.NewRequest("POST", rutaCompleta, bytes.NewBuffer(reqBody))
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
			noticiaContenido := map[string]interface{}{
				"titulo":      "",
				"descripcion": "",
				"link":        "",
				"activo":      noti["Activo"],
				"estilo":      noti["IdTipoEstilo"].(map[string]interface{})["Id"],
				"prioridad":   noti["IdTipoPrioridad"].(map[string]interface{})["Id"],
			}
			var responseNoticiaContenido []map[string]interface{}
			errNoticiaContenido := request.GetJson("http://"+beego.AppConfig.String("noticiaService")+fmt.Sprintf("/noticia_tipo_contenido?query=IdNoticia__id:%v", noti["Id"]), &responseNoticiaContenido)
			if errNoticiaContenido == nil {
				for _, conte := range responseNoticiaContenido {
					fmt.Println(reflect.TypeOf(conte["IdTipoContenido"].(map[string]interface{})["Id"]))
					fmt.Println(conte["IdTipoContenido"].(map[string]interface{})["Id"])
					dato := strings.ReplaceAll(conte["Dato"].(string), "{\"dato\": \"", "")
					dato = strings.TrimSuffix(dato, "\"}")
					if conte["IdTipoContenido"].(map[string]interface{})["Id"] == float64(1) {
						fmt.Println("Entró")
						noticiaContenido["titulo"] = dato
					}
					if conte["IdTipoContenido"].(map[string]interface{})["Id"] == float64(2) {
						noticiaContenido["descripcion"] = dato
					}
					if conte["IdTipoContenido"].(map[string]interface{})["Id"] == float64(3) {
						noticiaContenido["link"] = dato
					}

					listado = append(listado, noticiaContenido)
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
