package middleware

import (
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		if session.Get("login") != true {
			c.Redirect(http.StatusTemporaryRedirect, "/login")
			return
		}
		c.Next()
	}
}
