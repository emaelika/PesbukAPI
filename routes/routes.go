package routes

import (
	"PesbukAPI/config"
	"PesbukAPI/features/comment"
	"PesbukAPI/features/post"
	"PesbukAPI/features/user"

	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
)

func InitRoute(c *echo.Echo, ctl user.UserController, pc post.PostController, cc comment.CommentController) {
	c.POST("/register", ctl.Add()) // register -> umum (boleh diakses semua orang)
	c.POST("/login", ctl.Login())
	c.GET("/users", ctl.Profile(), echojwt.WithConfig(echojwt.Config{
		SigningKey: []byte(config.JWTSECRET),
	})) 
	c.PUT("/users", ctl.Update(), echojwt.WithConfig(echojwt.Config{
		SigningKey: []byte(config.JWTSECRET),
	}))
	c.DELETE("/users", ctl.Delete(), echojwt.WithConfig(echojwt.Config{
		SigningKey: []byte(config.JWTSECRET),
	}))
	c.GET("/users/:id", ctl.GetUserByIDParam(), echojwt.WithConfig(echojwt.Config{
		SigningKey: []byte(config.JWTSECRET),
	}))
	c.GET("/avatar", ctl.Avatar(), echojwt.WithConfig(echojwt.Config{
		SigningKey: []byte(config.JWTSECRET),
	}))



	c.POST("/posts", pc.Add(), echojwt.WithConfig(echojwt.Config{
		SigningKey: []byte(config.JWTSECRET),
	}))
	c.GET("/posts", pc.ShowAllPosts())

	c.GET("/posts/:id", pc.ShowPostByID())

	c.PUT("/posts/:id", pc.Update(), echojwt.WithConfig(echojwt.Config{
		SigningKey: []byte(config.JWTSECRET),
	}))
	c.DELETE("/posts/:id", pc.Delete(), echojwt.WithConfig(echojwt.Config{
		SigningKey: []byte(config.JWTSECRET),
	}))






	c.POST("/posts/:post_id/comments", cc.Add(), echojwt.WithConfig(echojwt.Config{
		SigningKey: []byte(config.JWTSECRET),
	}))
	
	c.GET("/comments", cc.ShowAllComments())
	
	c.DELETE("/comments/:comments_id", cc.Delete(), echojwt.WithConfig(echojwt.Config{
		SigningKey: []byte(config.JWTSECRET),
	}))
	c.PUT("/comments/:comments_id", cc.Update(), echojwt.WithConfig(echojwt.Config{
		SigningKey: []byte(config.JWTSECRET),
	}))
}
