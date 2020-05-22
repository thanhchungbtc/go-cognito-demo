package app

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func me() gin.HandlerFunc {
	return func(c *gin.Context) {
		claims, _ := c.Get("claims")
		c.JSON(http.StatusOK, claims)
	}
}
