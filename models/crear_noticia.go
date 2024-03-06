package models

type NoticiaRequest struct {
	Noticia           Noticia           `json:"Noticia"`
	Etiqueta          Etiqueta          `json:"Etiqueta"`
	Contenido         Contenido         `json:"Contenido"`
	ModuloPublicacion ModuloPublicacion `json:"ModuloPublicacion"`
}

// NoticiaSend representa la estructura para enviar la noticia al cliente en la ventana de listar noticias
type NoticiaSend struct {
	Noticia   NoticiaGetAll `json:"Noticia"`
	Etiquetas []Etiqueta    `json:"Etiquetas"`
	Contenido []Contenido   `json:"Contenido"`
}

type Noticia struct {
	Activo       bool `json:"Activo"`
	IdTipoEstilo struct {
		Id int `json:"Id"`
	} `json:"IdTipoEstilo"`
	IdTipoPrioridad struct {
		Id int `json:"Id"`
	} `json:"IdTipoPrioridad"`
}

type Etiqueta struct {
	Activo    bool `json:"Activo"`
	IdNoticia struct {
		Id int `json:"Id"`
	} `json:"IdNoticia"`
	IdTipoEtiqueta []int `json:"IdTipoEtiqueta"`
}

// EtiquetaData representa la estructura para enviar datos de etiqueta a la API CRUD
type EtiquetaData struct {
	Activo    bool `json:"Activo"`
	IdNoticia struct {
		Id int `json:"Id"`
	} `json:"IdNoticia"`
	IdEtiqueta int `json:"IdEtiqueta"`
}

type Contenido struct {
	Id   []int    `json:"Id"`
	Dato []string `json:"Dato"`
}

// ContenidoData representa la estructura para enviar datos del contenido a la API CRUD
type ContenidoData struct {
	Activo    bool   `json:"Activo"`
	Dato      string `json:"Dato"`
	IdNoticia struct {
		Id int `json:"Id"`
	} `json:"IdNoticia"`
	IdContenido int `json:"IdContenido"`
}

type ModuloPublicacion struct {
	IdModulo []string `json:"IdModulo"`
}

// ModuloPublicaionData representa la estructura para enviar datos del contenido a la API CRUD
type ModuloPublicacionData struct {
	Activo    bool `json:"Activo"`
	IdNoticia struct {
		Id int `json:"Id"`
	} `json:"IdNoticia"`
	RefModuloId string `json:"RefModuloId"`
}

type APIResponse struct {
	Body []byte
	Err  error
}

type NoticiaResponse struct {
	Data struct {
		ID int `json:"Id"`
	} `json:"Data"`
}

type NoticiaGetAll struct {
	Id        int  `json:"Id"`
	Activo    bool `json:"Activo"`
	IdEstilo  int  `json:"IdEstilo"`
	Prioridad int  `json:"Prioridad"`
}
