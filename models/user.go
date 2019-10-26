package models

type User struct {
	Username string `json:"username" form:"username" binding:"required"`
	Role     uint
	Password string `json:"password" form:"password" binding:"required"`
}
