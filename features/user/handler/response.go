package handler

type LoginResponse struct {
	Email string `json:"email"`
	Name  string `json:"name"`
	Token string `json:"token"`
}
type ProfilResponse struct {
	Email string `json:"email"`
	Name  string `json:"name"`
}
