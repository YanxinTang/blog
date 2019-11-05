package controllers

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/YanxinTang/blog/config"
	"github.com/YanxinTang/blog/middleware"
	"github.com/YanxinTang/blog/models"
	"github.com/YanxinTang/blog/utils"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

const PerPage uint64 = 10

func Articles(c *gin.Context) {
	pageParam := c.Param("page")
	var page uint64
	if pageParam == "" {
		page = 1
	} else {
		var err error
		page, err = strconv.ParseUint(c.Param("page"), 10, 64)
		if page == 0 || err != nil {
			c.HTML(http.StatusNotFound, "error/404", nil)
			return
		}
	}
	rows, err := models.GetPosts(page, PerPage)
	if err != nil {
		c.Error(err).SetType(http.StatusNotFound)
	}
	var articles []*models.Article
	for rows.Next() {
		var article models.Article
		err := rows.Scan(&article.ID, &article.CategoryID, &article.Title, &article.Content, &article.CreatedAt, &article.UpdatedAt)
		if err != nil {
			continue
		}
		sqlRow, err := models.GetCategory(article.CategoryID)
		if err != nil {
			continue
		}
		var category models.Category
		if err := sqlRow.Scan(&category.ID, &category.Name); err != nil {
			continue
		}
		article.Category = category
		articles = append(articles, &article)
	}
	count := models.ArticlesCount()
	session := sessions.Default(c)
	pagination := NewPagination(count, page, PerPage, "/page/%d").Init()
	c.HTML(http.StatusOK, "blog/index", gin.H{
		"title":      siteName,
		"login":      session.Get("login"),
		"articles":   articles,
		"count":      count,
		"page":       page,
		"perPage":    PerPage,
		"pagination": pagination.Render(),
	})
}

func ArticleView(c *gin.Context) {
	var article models.Article
	var comments []*models.Comment
	articleID, err := strconv.ParseUint(c.Param("articleID"), 10, 64)

	if err != nil {
		c.HTML(http.StatusNotFound, "error/404", nil)
		return
	}

	row, err := models.GetArticle(articleID)
	if err != nil {
		c.HTML(http.StatusNotFound, "error/404", nil)
		return
	}
	var categoryName string
	if err := row.Scan(&article.ID, &article.CategoryID, &article.Title, &article.Content, &article.CreatedAt, &article.UpdatedAt, &categoryName); err != nil {
		c.Error(err).SetType(http.StatusNotFound)
		return
	}

	commentsRows, err := models.GetArticleComments(articleID)
	if err != nil {
		c.Error(err).SetType(http.StatusNotFound)
		return
	}
	for commentsRows.Next() {
		var comment models.Comment
		if err := commentsRows.Scan(
			&comment.ID,
			&comment.ArticleID,
			&comment.ParentCommentID,
			&comment.Username,
			&comment.Content,
			&comment.CreatedAt,
			&comment.UpdatedAt,
		); err != nil {
			continue
		}
		comments = append(comments, &comment)
	}
	session := sessions.Default(c)
	errorMsgs := session.Flashes("errorMsgs")
	successMsgs := session.Flashes("successMsgs")
	session.Save()

	c.HTML(http.StatusOK, "blog/article", gin.H{
		"title":        utils.SiteTitle(article.Title, siteName),
		"login":        session.Get("login"),
		"username":     config.Config.Auth.Username,
		"article":      article,
		"categoryName": categoryName,
		"comments":     comments,
		"errorMsgs":    errorMsgs,
		"successMsgs":  successMsgs,
	})
}

func AddComment(c *gin.Context) {
	articleID, err := strconv.ParseUint(c.Param("articleID"), 10, 64)
	if err != nil {
		log.Println("error: ", err)
	}
	var comment models.Comment
	comment.ArticleID = articleID
	session := sessions.Default(c)
	if err := c.ShouldBind(&comment); err != nil {
		session.AddFlash("评论内容不能为空", "errorMsgs")
		session.Save()
		c.Error(err).SetType(http.StatusBadRequest).SetMeta(middleware.BadRequestMeta{
			Url: fmt.Sprintf("/articles/%d", articleID),
		})
		return
	}

	comment.Username = strings.TrimSpace(comment.Username)
	if comment.Username == "" {
		comment.Username = "匿名"
	}

	_, err = models.AddComment(&comment)
	if err != nil {
		c.Error(err).SetType(http.StatusNotFound)
		return
	}
	session.AddFlash("评论添加成功", "successMsgs")
	session.Save()
	c.Redirect(http.StatusFound, fmt.Sprintf("/articles/%d", articleID))
}

func DeleteComment(c *gin.Context) {
	articleID, err := strconv.ParseUint(c.Param("articleID"), 10, 64)
	if err != nil {
		c.Error(err).SetType(http.StatusNotFound)
		return
	}
	commentID, err := strconv.ParseUint(c.Param("commentID"), 10, 64)
	if err != nil {
		c.Error(err).SetType(http.StatusNotFound)
		return
	}
	if _, err = models.DeleteComment(commentID); err != nil {
		c.Error(err).SetType(http.StatusNotFound)
		return
	}
	session := sessions.Default(c)
	session.AddFlash("删除成功", "successMsgs")
	session.Save()
	c.Redirect(http.StatusFound, fmt.Sprintf("/articles/%d", articleID))
}
