package models

import (
	"database/sql"
	"fmt"
	"log"
)

type Category struct {
	BaseModel
	Name string `form:"name" binding:"required"`
}

func InsertCategory(category *Category) (sql.Result, error) {
	stmt, err := db.Prepare("INSERT category SET name=?")
	if err != nil {
		log.Printf("insert category error: %v", err)
	}
	return stmt.Exec(&category.Name)
}

func GetCategory(ID uint64) (*sql.Row, error) {
	stmt, err := db.Prepare("SELECT `id`, `name` FROM category WHERE id = ?")
	if err != nil {
		return nil, err
	}
	return stmt.QueryRow(ID), nil
}

func GetCategories(args ...interface{}) (*sql.Rows, error) {
	if len(args) == 0 {
		return db.Query("SELECT * FROM category")
	}

	var colum string
	length := len(args)
	for i := 0; i < length; i++ {
		if i == length-1 {
			colum += fmt.Sprintf("`%s`", args[i])
		} else {
			colum += fmt.Sprintf("`%s`, ", args[i])
		}
	}
	return db.Query(fmt.Sprintf("SELECT %s FROM category", colum))
}

func DeleteCategory(categoryID uint64) (sql.Result, error) {
	return db.Exec("DELETE FROM category WHERE id = ?", categoryID)
}

// CategoriesCount returns count of categories
func CategoriesCount() uint64 {
	var count uint64
	db.QueryRow("SELECT COUNT(*) FROM category").Scan(&count)
	return count
}
