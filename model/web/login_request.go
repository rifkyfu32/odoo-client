package web

type LoginRequest struct {
	DB       string `json:"db"`
	Login    string `json:"login"`
	Password string `json:"password"`
}
