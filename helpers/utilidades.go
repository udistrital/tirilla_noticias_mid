package helpers

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/udistrital/noticias_mid/models"
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
