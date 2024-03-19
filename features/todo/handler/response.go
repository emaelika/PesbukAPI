package handler

type TodoResponse struct {
	Kegiatan  string `json:"kegiatan" form:"kegiatan" validate:"required,max=25"`
	Deskripsi string `json:"deskripsi" form:"deskripsi" validate:"required,max=250"`
	Deadline  string `json:"deadline" form:"deadline" validate:"required,datetime=2006-01-02"`
}
type GetTodosResponse struct {
	ID        uint   `json:"id" form:"id" `
	Kegiatan  string `json:"kegiatan" form:"kegiatan" validate:"required,max=25"`
	Deskripsi string `json:"deskripsi" form:"deskripsi" validate:"required,max=250"`
	Deadline  string `json:"deadline" form:"deadline" validate:"required,datetime=2006-01-02"`
}
