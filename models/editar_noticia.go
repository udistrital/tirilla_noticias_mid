package models

type EtiquetaResponse struct {
	Data []struct {
		Id         int  `json:"Id"`     //Id de la tabla noticia_tipo_etiqueta
		Activo     bool `json:"Activo"` //Estado de la etiqueta
		IdEtiqueta int  `json:"IdEtiqueta"`
	} `json:"Data"`
}

// ContenidoResponse representa la estructura para decodificar la respuesta de la API CRUD para el contenido
type ContenidoResponse struct {
	Data []struct {
		Id          int    `json:"Id"`
		Activo      bool   `json:"Activo"`
		Dato        string `json:"Dato"`
		IdContenido int    `json:"IdContenido"`
	} `json:"data"`
}
