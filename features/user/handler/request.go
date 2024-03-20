package handler

type LoginRequest struct {
	Email    string `json:"email" form:"email" validate:"required;email"`
	Password string `json:"password" form:"password" validate:"required"`
}

type RegisterRequest struct {
	Name         string `json:"name" form:"name" validate:"required,alpha"`
	Email        string `json:"email" form:"email" validate:"required,email"`
	Username     string `json:"username" form:"username" validate:"required"`
	Placeofbirth string `json:"placeofbirth" form:"placeofbirth" validate:"required,alpha"`
	Dateofbirth  string `json:"dateofbirth" form:"dateofbirth" validate:"required,datetime"`
	Password     string `json:"password" form:"password" validate:"required"`
}
