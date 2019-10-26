package controllers

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/YanxinTang/blog/middleware"
	"github.com/YanxinTang/blog/models"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func DashBoardView(c *gin.Context) {
	session := sessions.Default(c)
	c.HTML(http.StatusOK, "admin/dashboard", gin.H{
		"login": session.Get("login"),
	})
}

func GetArticle(c *gin.Context) {
	var post models.Article
	articleID, err := strconv.ParseUint(c.Param("articleID"), 10, 64)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"message": err.Error(),
		})
		return
	}

	row, err := models.GetArticle(articleID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}
	if err := row.Scan(&post.ID, &post.CategoryID, &post.Title, &post.Content, &post.CreatedAt, &post.UpdatedAt); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "ok",
		"data": gin.H{
			"postID":     post.ID,
			"title":      post.Title,
			"categoryID": post.CategoryID,
			"content":    post.Content,
			"createdAt":  post.CreatedAt,
			"updatedAt":  post.UpdatedAt,
		},
	})
}

func AddArticleView(c *gin.Context) {
	var categories []*models.Category
	categoriesRows, err := models.GetCategories("id", "name")
	if err != nil {
		c.Error(err).SetType(http.StatusBadRequest).SetMeta(middleware.BadRequestMeta{
			Url: "/admin/articles/new",
		})
		return
	}
	for categoriesRows.Next() {
		var category models.Category
		if err := categoriesRows.Scan(&category.ID, &category.Name); err != nil {
			continue
		}
		categories = append(categories, &category)
	}
	session := sessions.Default(c)
	errorMsgs := session.Flashes("errorMsgs")
	successMsgs := session.Flashes("successMsgs")
	session.Save()
	c.HTML(http.StatusOK, "admin/addArticle", gin.H{
		"login":       session.Get("login"),
		"categories":  categories,
		"errorMsgs":   errorMsgs,
		"successMsgs": successMsgs,
	})
}

func AddArticle(c *gin.Context) {
	session := sessions.Default(c)
	var article models.Article
	if err := c.ShouldBind(&article); err != nil {
		switch {
		case article.Title == "":
			session.AddFlash("标题不能为空", "errorMsgs")
		case article.Content == "":
			session.AddFlash("内容不能为空", "errorMsgs")
		}
		session.Save()
		c.Error(err).SetType(http.StatusBadRequest).SetMeta(middleware.BadRequestMeta{
			"/admin/articles/new",
		})
		return
	}
	if _, err := models.InsertArticle(&article); err != nil {
		c.Error(err).SetType(http.StatusBadRequest).SetMeta(middleware.BadRequestMeta{
			"/admin/articles/new",
		})
		return
	}
	session.AddFlash("发布成功", "successMsgs")
	session.Save()
	c.Redirect(http.StatusFound, "/admin/articles/new")
}

func DeleteArticle(c *gin.Context) {
	articleID, err := strconv.ParseUint(c.Param("articleID"), 10, 64)

	if err != nil {
		c.HTML(http.StatusNotFound, "error/404", nil)
		return
	}

	sqlResult, err := models.DeleteArticle(articleID)
	if err != nil {
		c.HTML(http.StatusNotFound, "error/404", nil)
		return
	}

	rowsAffected, err := sqlResult.RowsAffected()
	if err != nil {
		c.HTML(http.StatusNotFound, "error/404", nil)
		return
	}
	if rowsAffected >= 1 {
		c.Redirect(http.StatusFound, "/")
		return
	}
	c.Redirect(http.StatusFound, fmt.Sprintf("/posts/%d", articleID))
}

func UpdateArticleView(c *gin.Context) {
	articleID, err := strconv.ParseUint(c.Param("articleID"), 10, 64)
	if err != nil {
		c.Error(err).SetType(http.StatusNotFound)
		return
	}
	var article models.Article
	sqlRow, err := models.GetArticle(articleID, "id", "category_id", "title", "content")
	if err != nil {
		c.Error(err).SetType(http.StatusNotFound)
		return
	}
	if err := sqlRow.Scan(&article.ID, &article.CategoryID, &article.Title, &article.Content); err != nil {
		c.Error(err).SetType(http.StatusNotFound)
		return
	}
	var categories []*models.Category
	categoriesRows, err := models.GetCategories("id", "name")
	if err != nil {
		c.Error(err).SetType(http.StatusBadRequest).SetMeta(middleware.BadRequestMeta{
			Url: fmt.Sprintf("/admin/articles/%d", articleID),
		})
		return
	}
	for categoriesRows.Next() {
		var category models.Category
		if err := categoriesRows.Scan(&category.ID, &category.Name); err != nil {
			continue
		}
		categories = append(categories, &category)
	}
	session := sessions.Default(c)
	errorMsgs := session.Flashes("errorMsgs")
	successMsgs := session.Flashes("successMsgs")
	session.Save()
	c.HTML(http.StatusOK, "admin/updateArticle", gin.H{
		"login":       session.Get("login"),
		"article":     article,
		"categories":  categories,
		"errorMsgs":   errorMsgs,
		"successMsgs": successMsgs,
	})
}

func UpdateArticle(c *gin.Context) {
	articleID, err := strconv.ParseUint(c.Param("articleID"), 10, 64)
	if err != nil {
		c.Error(err).SetType(http.StatusNotFound)
		return
	}
	referer := fmt.Sprintf("/admin/articles/update/%d", articleID)
	session := sessions.Default(c)
	var article models.Article
	if err := c.ShouldBind(&article); err != nil {
		log.Println(err)
		switch {
		case article.Title == "":
			session.AddFlash("标题不能为空", "errorMsgs")
		case article.Content == "":
			session.AddFlash("内容不能为空", "errorMsgs")
		}
		session.Save()
		c.Error(err).SetType(http.StatusBadRequest).SetMeta(middleware.BadRequestMeta{referer})
		return
	}

	article.ID = articleID

	if _, err := models.UpdateArticle(&article); err != nil {
		session.AddFlash("更新失败", "errorMsgs")
		session.Save()
		c.Error(err).SetType(http.StatusBadRequest).SetMeta(middleware.BadRequestMeta{referer})
	}

	session.AddFlash("更新成功", "successMsgs")
	session.Save()
	c.Redirect(http.StatusFound, referer)
}

func CategoriesView(c *gin.Context) {
	var categories []*models.Category
	categoryRows, err := models.GetCategories()
	if err != nil {
		c.Error(err).SetType(http.StatusNotFound)
		return
	}
	for categoryRows.Next() {
		var category models.Category
		if err := categoryRows.Scan(&category.ID, &category.Name, &category.CreatedAt, &category.UpdatedAt); err != nil {
			continue
		}
		categories = append(categories, &category)
	}
	session := sessions.Default(c)
	errorMsgs := session.Flashes("errorMsgs")
	successMsgs := session.Flashes("successMsgs")
	session.Save()
	c.HTML(http.StatusOK, "admin/categories", gin.H{
		"categories":  categories,
		"errorMsgs":   errorMsgs,
		"successMsgs": successMsgs,
	})
}

// AddCategory add category
func AddCategory(c *gin.Context) {
	var category models.Category
	session := sessions.Default(c)
	referer := "/admin/categories"
	if err := c.ShouldBind(&category); err != nil {
		if category.Name == "" {
			session.AddFlash("分类名称不能为空", "errorMsgs")
		}
		session.Save()
		c.Error(err).SetType(http.StatusBadRequest).SetMeta(middleware.BadRequestMeta{referer})
		return
	}

	if _, err := models.InsertCategory(&category); err != nil {
		session.AddFlash("添加失败", "errorMsgs")
		session.Save()
		c.Error(err).SetType(http.StatusBadRequest).SetMeta(middleware.BadRequestMeta{referer})
		return
	}
	session.AddFlash("添加成功", "successMsgs")
	session.Save()
	c.Redirect(http.StatusFound, referer)
}

func DeleteCategory(c *gin.Context) {
	categoryID, err := strconv.ParseUint(c.Param("categoryID"), 10, 64)
	if err != nil {
		c.Error(err).SetType(http.StatusNotFound)
		return
	}
	session := sessions.Default(c)
	referer := "/admin/categories/"
	if _, err := models.DeleteCategory(categoryID); err != nil {
		log.Println(err)
		session.AddFlash("删除失败", "errorMsgs")
		session.Save()
		c.Error(err).SetType(http.StatusBadRequest).SetMeta(middleware.BadRequestMeta{referer})
		return
	}

	session.AddFlash("删除成功", "successMsgs")
	session.Save()
	c.Redirect(http.StatusFound, referer)
}
