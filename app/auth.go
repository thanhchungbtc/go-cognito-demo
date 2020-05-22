package app

import (
	"net/http"
	"strings"
	"thanhchungbtc/go-cognito-demo/domain/usecase"

	"github.com/gin-gonic/gin"
)

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
