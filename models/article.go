package models

import (
	"database/sql"
	"fmt"

	"github.com/YanxinTang/blog/utils"
)

type Article struct {
	BaseModel
	CategoryID uint64 `form:"categoryID" binding:"required"`
	Category   `binding:"-"`
	Title      string `form:"title" binding:"required"`
	Content    string `form:"content" binding:"required"`
}

func InsertArticle(article *Article) (sql.Result, error) {
	stmt, err := db.Prepare("INSERT article SET title=?, category_id = ?, content=?")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	return stmt.Exec(&article.Title, &article.CategoryID, &article.Content)
}

func GetArticle(articleID uint64, columns ...string) (*sql.Row, error) {
	if len(columns) == 0 {
		return db.QueryRow("SELECT article.*, category.name FROM article INNER  JOIN category ON article.category_id = category.id WHERE article.id = ?", articleID), nil
	}
	columnPlaceholder := utils.ColumnPlaceholder(columns...)
	sql := fmt.Sprintf("SELECT %s FROM article WHERE id = ?", columnPlaceholder)
	return db.QueryRow(sql, articleID), nil

}

func DeleteArticle(articleID uint64) (sql.Result, error) {
	stmt, err := db.Prepare("DELETE FROM article WHERE id = ?")
	defer stmt.Close()
	if err != nil {
		return nil, err
	}
	return stmt.Exec(articleID)
}

func UpdateArticle(article *Article) (sql.Result, error) {
	stmt, err := db.Prepare("UPDATE article SET category_id = ?, title = ?, content = ? WHERE id = ?")
	defer stmt.Close()
	if err != nil {
		return nil, err
	}
	return stmt.Exec(&article.CategoryID, &article.Title, &article.Content, &article.ID)
}

func GetPosts(page uint64, perPage uint64) (*sql.Rows, error) {
	start := (page - 1) * perPage
	return db.Query("SELECT * FROM article ORDER BY id DESC LIMIT ?, ?", start, perPage)
}

func ArticlesCount() uint64 {
	var count uint64
	db.QueryRow("SELECT COUNT(*) FROM article").Scan(&count)
	return count
}
