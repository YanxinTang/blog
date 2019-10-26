package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type BadRequestMeta struct {
	Url string
}

func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		if c.Errors != nil {
			if err := c.Errors.Last(); err != nil {
				switch err.Type {
				case http.StatusNotFound:
					c.HTML(http.StatusNotFound, "error/404", err.Meta)
				case http.StatusBadRequest:
					c.Redirect(http.StatusFound, err.Meta.(BadRequestMeta).Url)
				}
			}
		}
	}
}
