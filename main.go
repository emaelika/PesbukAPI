package main

import (
	"PesbukAPI/config"
	cr "PesbukAPI/features/comment/data"
	ch "PesbukAPI/features/comment/handler"
	cs "PesbukAPI/features/comment/service"
	pr "PesbukAPI/features/post/data"
	ph "PesbukAPI/features/post/handler"
	ps "PesbukAPI/features/post/service"
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

	cq := cr.New(db) // bagian yang menghungkan coding kita ke database / bagian dimana kita ngoding untk ke DB
	cs := cs.NewCommentService(cq)
	ch := ch.NewHandler(cs)
	// bagian yang menghandle segala hal yang berurusan dengan HTTP / echo

	pq := pr.New(db)
	ps := ps.NewPostService(pq)
	ph := ph.NewHandler(ps)

	e.Pre(middleware.RemoveTrailingSlash())
	e.Use(middleware.Logger())
	e.Use(middleware.CORS()) // ini aja cukup
	routes.InitRoute(e, uh, ph, ch)
	e.Logger.Fatal(e.Start(":1323"))
}
