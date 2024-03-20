package main

import (
	"PesbukAPI/config"
	th "PesbukAPI/features/todo/handler"
	tr "PesbukAPI/features/todo/repository"
	ts "PesbukAPI/features/todo/service"
	ur "PesbukAPI/features/user/data"
	uh "PesbukAPI/features/user/handler"
	us "PesbukAPI/features/user/service"
	"PesbukAPI/routes"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	e := echo.New()            // inisiasi echo
	cfg := config.InitConfig() // baca seluruh system variable
	db := config.InitSQL(cfg)  // konek DB

	uq := ur.New(db) // bagian yang menghungkan coding kita ke database / bagian dimana kita ngoding untk ke DB
	us := us.NewService(uq)
	uh := uh.NewUserHandler(us)

	tq := tr.NewTodoQuery(db) // bagian yang menghungkan coding kita ke database / bagian dimana kita ngoding untk ke DB
	ts := ts.NewTodoService(tq)
	th := th.NewHandler(ts)
	// bagian yang menghandle segala hal yang berurusan dengan HTTP / echo

	e.Pre(middleware.RemoveTrailingSlash())
	e.Use(middleware.Logger())
	e.Use(middleware.CORS()) // ini aja cukup
	routes.InitRoute(e, uh, th)
	e.Logger.Fatal(e.Start(":1323"))
}
