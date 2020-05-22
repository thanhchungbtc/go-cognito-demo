package app

import (
	"net/http"
	"thanhchungbtc/go-cognito-demo/domain/usecase"

	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"
)

type app struct {
	r    *gin.Engine
	auth usecase.Auth
}

func New(auth usecase.Auth) *app {
	r := gin.Default()
	a := &app{r: r, auth: auth}

	r.Use(static.Serve("/", static.LocalFile("./public", true)))

	a.r.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})
	api := r.Group("/api")
	{
		api.POST("/login", a.login)
		api.POST("/register", a.register)
		api.POST("/confirm", a.confirm)

		protected := api.Use(authRequired(auth))
		{
			protected.GET("/me", me())
		}
	}

	return a
}

func (a *app) Run(addr string) error {
	return a.r.Run(addr)
}

func abortWithError(c *gin.Context, code int, err error) {
	c.AbortWithStatusJSON(code, gin.H{
		"error": err.Error(),
	})
}
