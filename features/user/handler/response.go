package handler

type LoginResponse struct {
	ID    uint   `json:"id"`
	Email string `json:"username"`
	Token string `json:"token"`
}