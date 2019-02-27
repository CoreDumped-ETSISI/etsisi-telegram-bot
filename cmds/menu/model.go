package menu

type MenuDia struct {
	PrimerPlato  []string `json:"primer"`
	SegundoPlato []string `json:"segundo"`
}

type ErrorMessage struct {
	Message string `json:"message"`
}
