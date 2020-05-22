package app

import (
	"net/http"
	"strings"
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

func (a *app) login(c *gin.Context) {
	var req usecase.LoginInput
	if err := c.ShouldBindJSON(&req); err != nil {
		abortWithError(c, http.StatusBadRequest, err)
		return
	}

	out, err := a.auth.Login(&req)
	if err != nil {
		abortWithError(c, http.StatusBadRequest, err)
		return
	}

	c.JSON(200, out)
}

func (a *app) register(c *gin.Context) {
	var req usecase.RegisterInput
	if err := c.ShouldBindJSON(&req); err != nil {
		abortWithError(c, http.StatusBadRequest, err)
		return
	}

	out, err := a.auth.Register(&req)
	if err != nil {
		abortWithError(c, http.StatusBadRequest, err)
		return
	}

	c.JSON(http.StatusOK, out)
}

func (a *app) confirm(c *gin.Context) {
	var req usecase.ConfirmInput
	if err := c.ShouldBindJSON(&req); err != nil {
		abortWithError(c, http.StatusBadRequest, err)
		return
	}
	out, err := a.auth.Confirm(&req)
	if err != nil {
		abortWithError(c, http.StatusBadRequest, err)
		return
	}
	c.JSON(http.StatusOK, out)
}

func authRequired(auth usecase.Auth) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		tokenStr := strings.Split(authHeader, "Bearer ")[1]

		token, err := auth.ParseAndVerifyJWT(tokenStr)
		if err != nil {
			abortWithError(c, http.StatusUnauthorized, err)
			return
		}

		println(token)
	}
}

func me() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, "ok")
	}
}
