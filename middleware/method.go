package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
)

type MethodSet []string

var overrideMethods = MethodSet{"PUT", "PATCH", "DELETE"}

func Method(r *gin.Engine) gin.HandlerFunc {
	return func(c *gin.Context) {
		method := strings.ToUpper(c.PostForm("_method"))
		if c.Request.Method == "POST" && overrideMethods.include(method) {
			c.Request.Method = method
			r.HandleContext(c)
			c.Abort()
			return
		}
		c.Next()
	}
}

func (s *MethodSet) include(method string) bool {
	for _, m := range *s {
		if m == method {
			return true
		}
	}
	return false
}
