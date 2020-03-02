package models

import (
	"database/sql"
)

type Comment struct {
	BaseModel
	ArticleID       uint64 `form:"articleID" binding:"required"`
	ParentCommentID uint64 `form:"parentCommentID"`
	Username        string `form:"username"`
	Content         string `form:"content"  binding:"required"`
}

func AddComment(comment *Comment) (sql.Result, error) {
	stmt, err := db.Prepare("INSERT comment SET article_id = ?, username = ?, content = ?")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	return stmt.Exec(comment.ArticleID, comment.Username, comment.Content)
}

func DeleteComment(commentID uint64) (sql.Result, error) {
	stmt, err := db.Prepare("DELETE FROM comment WHERE id = ?")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	return stmt.Exec(commentID)
}

func GetArticleComments(articleID uint64) (*sql.Rows, error) {
	stmt, err := db.Prepare("SELECT * FROM comment WHERE article_id = ?")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	return stmt.Query(articleID)
}

// CommentsCount returns count of comment
func CommentsCount() uint64 {
	var count uint64
	db.QueryRow("SELECT COUNT(*) FROM comment").Scan(&count)
	return count
}
