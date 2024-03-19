package routes

import (
	"21-api/config"
	"21-api/features/todo"
	"21-api/features/user"

	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
)

func InitRoute(c *echo.Echo, ctl user.UserController, tc todo.TodoController) {
	c.POST("/users", ctl.Register()) // register -> umum (boleh diakses semua orang)
	c.POST("/login", ctl.Login())
	// c.GET("/users", ctl.ListUser(), echojwt.WithConfig(echojwt.Config{
	// 	SigningKey: []byte(config.JWTSECRET),
	// })) // get all user -> butuh penanda khusus
	c.GET("/users/", ctl.Profile(), echojwt.WithConfig(echojwt.Config{
		SigningKey: []byte(config.JWTSECRET),
	})) // get profile -> butuh penanda khusus
	// c.PUT("/users/:id", ctl.Update(), echojwt.WithConfig(echojwt.Config{
	// 	SigningKey: []byte(config.JWTSECRET),
	// })) // update user -> butuh penanda khusus

	c.POST("/todos", tc.AddTodo(), echojwt.WithConfig(echojwt.Config{
		SigningKey: []byte(config.JWTSECRET),
	}))
	c.GET("/todos", tc.GetTodos(), echojwt.WithConfig(echojwt.Config{
		SigningKey: []byte(config.JWTSECRET),
	}))
	c.GET("/todos/:id", tc.GetTodo(), echojwt.WithConfig(echojwt.Config{
		SigningKey: []byte(config.JWTSECRET),
	}))
	c.PUT("/todos/:id", tc.UpdateTodo(), echojwt.WithConfig(echojwt.Config{
		SigningKey: []byte(config.JWTSECRET),
	}))
}
