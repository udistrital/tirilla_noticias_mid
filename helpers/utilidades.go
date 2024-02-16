package helpers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

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
	errNoticia := request.GetJson("http://"+beego.AppConfig.String("noticiaService")+fmt.Sprintf("/noticia"), &noticia)
	if errNoticia == nil {
		fmt.Println("http://" + beego.AppConfig.String("noticiaService"))

		for _, noti := range noticia {
			noticiaContenido := map[string]interface{}{
				"titulo":      "",
				"descripcion": "",
				"link":        "",
			}
			var responseNoticiaContenido []map[string]interface{}
			errNoticiaContenido := request.GetJson("http://"+beego.AppConfig.String("noticiaService")+fmt.Sprintf("/noticia_tipo_contenido?query=IdNoticia__id:%v", noti["Id"]), &responseNoticiaContenido)
			if errNoticiaContenido == nil {
				for _, conte := range responseNoticiaContenido {
					if conte["IdTipoContenido"].(map[string]interface{})["Id"] == 1 {
						noticiaContenido["titulo"] = conte["Dato"]
					}
					if conte["IdTipoContenido"].(map[string]interface{})["Id"] == 2 {
						noticiaContenido["descripcion"] = conte["Dato"]
					}
					if conte["IdTipoContenido"].(map[string]interface{})["Id"] == 3 {
						noticiaContenido["link"] = conte["Dato"]
					}

					listado = append(listado, noticiaContenido)
				}
			}
		}

		APIResponseDTO = requestresponse.APIResponseDTO(true, 200, listado)

	} else {
		APIResponseDTO = requestresponse.APIResponseDTO(false, 400, nil, errNoticia.Error())
	}
	return APIResponseDTO
}
