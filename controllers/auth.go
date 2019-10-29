package controllers

import (
	"errors"
	"net/http"

	"github.com/YanxinTang/blog/config"
	"github.com/YanxinTang/blog/middleware"
	"github.com/YanxinTang/blog/models"
	"github.com/YanxinTang/blog/utils"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func LoginView(c *gin.Context) {
	session := sessions.Default(c)
	if session.Get("login") == true {
		c.Redirect(http.StatusTemporaryRedirect, "/app")
		return
	}
	errorMsgs := session.Flashes("errorMsgs")
	successMsgs := session.Flashes("successMsgs")
	session.Save()
	c.HTML(http.StatusOK, "blog/login", gin.H{
		"title":       utils.SiteTitle("登录", siteName),
		"errorMsgs":   errorMsgs,
		"successMsgs": successMsgs,
	})
}

func Login(c *gin.Context) {
	var user models.User
	session := sessions.Default(c)
	if err := c.ShouldBind(&user); err != nil {
		session.AddFlash("无效的用户信息", "errorMsgs")
		session.Save()
		c.Error(err).SetType(http.StatusBadRequest).SetMeta(middleware.BadRequestMeta{"/login/"})
		return
	}

	if user.Username != config.Config.Auth.Username || user.Password != config.Config.Auth.Password {
		session.AddFlash("密码错误", "errorMsgs")
		session.Save()
		c.Error(errors.New("密码错误")).SetType(http.StatusBadRequest).SetMeta(middleware.BadRequestMeta{"/login/"})
		return
	}

	session.Set("login", true)
	session.Save()
	c.Redirect(http.StatusSeeOther, "/admin/")
}

func Logout(c *gin.Context) {
	session := sessions.Default(c)
	session.Clear()
	session.Save()
	c.Redirect(http.StatusFound, "/")
}
