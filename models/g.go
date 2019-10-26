package models

import (
	"database/sql"
	"log"
	"time"

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
	db, err = sql.Open("mysql", "root:root@/blog?charset=utf8&parseTime=true")
	if err != nil {
		log.Fatalf("database connect error %v", err)
	}
}
