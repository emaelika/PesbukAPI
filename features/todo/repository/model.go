package repository

import "gorm.io/gorm"

type TodoModel struct {
	gorm.Model
	Kegiatan  string `json:"nama" form:"nama" validate:"required"`
	Deskripsi string `json:"hp" form:"hp" validate:"required,max=13,min=10"`
	Deadline  string `json:"password" form:"password" validate:"required"`
	UserID    uint
}
