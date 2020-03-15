package models

import (
	"database/sql"
)

type Draft struct {
	BaseModel
	CategoryID uint64 `json:"categoryID"`
	Category   `binding:"-"`
	Title      string `json:"title" binding:"required"`
	Content    string `json:"content"`
	Publish    bool
}

func InsertDraft(draft *Draft) (sql.Result, error) {
	stmt, err := db.Prepare("INSERT draft SET title=?, category_id = ?, content=?")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	return stmt.Exec(&draft.Title, &draft.CategoryID, &draft.Content)
}

func GetDraft(draftID uint64) *sql.Row {
	return db.QueryRow(
		"SELECT draft.*, category.name FROM draft INNER JOIN category ON draft.category_id = category.id WHERE draft.id = ?",
		draftID,
	)
}

func DeleteDraft(draftID uint64) (sql.Result, error) {
	stmt, err := db.Prepare("DELETE FROM draft WHERE id = ?")
	defer stmt.Close()
	if err != nil {
		return nil, err
	}
	return stmt.Exec(draftID)
}

func UpdateDraft(draft *Draft) (sql.Result, error) {
	return db.Exec(
		"UPDATE draft SET category_id = ?, title = ?, content = ? WHERE id = ?",
		&draft.CategoryID, &draft.Title, &draft.Content, &draft.ID,
	)
}

func GetUnpublishedDraftsByPage(page uint64, perPage uint64) (*sql.Rows, error) {
	start := (page - 1) * perPage
	return db.Query(
		"SELECT draft.*, category.name FROM `draft` INNER JOIN `category` ON draft.category_id = category.id WHERE draft.published = 0 ORDER BY `updated_at` DESC LIMIT ?, ?",
		start,
		perPage,
	)
}

func UnpublishedDraftsCount() uint64 {
	var count uint64
	db.QueryRow("SELECT COUNT(*) FROM draft WHERE draft.published = 0").Scan(&count)
	return count
}

func PublishDraft(draftID uint64) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	_, err = tx.Exec(
		"INSERT INTO `article` (`category_id`, `title`, `content`) SELECT `category_id`, `title`, `content` FROM `draft` WHERE `draft`.`id` = ?",
		draftID,
	)
	if err != nil {
		_ = tx.Rollback()
	}
	_, err = tx.Exec("UPDATE `draft` SET `published` = 1 WHERE id = ?", draftID)
	if err != nil {
		_ = tx.Rollback()
	}
	if err = tx.Commit(); err != nil {
		return err
	}
	return nil
}
