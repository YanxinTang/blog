package main

import (
	"database/sql"
	"log"

	"github.com/YanxinTang/blog/router"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	_, err := sql.Open("mysql", "root:root@/blog")
	if err != nil {
		log.Fatal(err)
	}
	engine := router.SetupRouter()

	// Listen and Server in 0.0.0.0:8080
	engine.Run(":8080")
}
