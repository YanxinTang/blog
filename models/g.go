package models

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/YanxinTang/blog/config"
	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

type BaseModel struct {
	ID        uint64 `json:"ID"`
	UpdatedAt time.Time
	CreatedAt time.Time
}

func init() {
	var err error
	mysql := &config.Config.Mysql
	connect := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=true",
		mysql.User,
		mysql.Password,
		mysql.Host,
		mysql.Port,
		mysql.Database,
	)
	db, err = sql.Open("mysql", connect)
	if err != nil {
		log.Fatalf("database connect error %v", err)
	}
}
