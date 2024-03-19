package main

import (
	"21-api/config"
	th "21-api/features/todo/handler"
	tr "21-api/features/todo/repository"
	ts "21-api/features/todo/service"
	uh "21-api/features/user/handler"
	ur "21-api/features/user/repository"
	us "21-api/features/user/service"
	"21-api/routes"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	e := echo.New()            // inisiasi echo
	cfg := config.InitConfig() // baca seluruh system variable
	db := config.InitSQL(cfg)  // konek DB

	uq := ur.NewUserQuery(db) // bagian yang menghungkan coding kita ke database / bagian dimana kita ngoding untk ke DB
	us := us.NewUserService(uq)
	uh := uh.NewHandler(us)

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
