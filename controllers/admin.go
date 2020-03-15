package controllers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/YanxinTang/blog/middleware"
	"github.com/YanxinTang/blog/models"
	"github.com/YanxinTang/blog/utils"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func DashBoardView(c *gin.Context) {
	session := sessions.Default(c)

	articlesCount := models.ArticlesCount()
	categoriesCount := models.CategoriesCount()
	commentsCount := models.CommentsCount()

	type Card struct {
		Title string
		Body  string
	}

	cards := []Card{}
	cards = append(
		cards,
		Card{"分类数量", strconv.FormatUint(categoriesCount, 10)},
		Card{"文章数量", strconv.FormatUint(articlesCount, 10)},
		Card{"评论数量", strconv.FormatUint(commentsCount, 10)},
	)

	c.HTML(http.StatusOK, "admin/dashboard", gin.H{
		"title": utils.SiteTitle("总览", siteName),
		"login": session.Get("login"),
		"menu":  "dashboard",
		"cards": cards,
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
		"title":       utils.SiteTitle("新增文章", siteName),
		"login":       session.Get("login"),
		"menu":        "addArticle",
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
		"title":       utils.SiteTitle("更新文章", siteName),
		"login":       session.Get("login"),
		"menu":        "none",
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
		return
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
		"title":       utils.SiteTitle("分类管理", siteName),
		"login":       session.Get("login"),
		"menu":        "categories",
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
		session.AddFlash("删除失败", "errorMsgs")
		session.Save()
		c.Error(err).SetType(http.StatusBadRequest).SetMeta(middleware.BadRequestMeta{referer})
		return
	}

	session.AddFlash("删除成功", "successMsgs")
	session.Save()
	c.Redirect(http.StatusFound, referer)
}

func DraftsView(c *gin.Context) {
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
	rows, err := models.GetUnpublishedDraftsByPage(page, PerPage)
	if err != nil {
		c.Error(err).SetType(http.StatusNotFound)
		return
	}
	var drafts []*models.Draft
	for rows.Next() {
		var draft models.Draft
		err := rows.Scan(
			&draft.ID,
			&draft.CategoryID,
			&draft.Title,
			&draft.Content,
			&draft.Publish,
			&draft.CreatedAt,
			&draft.UpdatedAt,
			&draft.Category.Name,
		)
		if err != nil {
			continue
		}
		drafts = append(drafts, &draft)
	}
	count := models.UnpublishedDraftsCount()
	pagination := NewPagination(count, page, PerPage, "/admin/drafts/page/%d").Init()

	session := sessions.Default(c)
	errorMsgs := session.Flashes("errorMsgs")
	successMsgs := session.Flashes("successMsgs")
	c.HTML(http.StatusOK, "admin/drafts", gin.H{
		"title":       utils.SiteTitle("草稿箱", siteName),
		"login":       session.Get("login"),
		"menu":        "drafts",
		"errorMsgs":   errorMsgs,
		"successMsgs": successMsgs,
		"drafts":      drafts,
		"count":       count,
		"page":        page,
		"perPage":     PerPage,
		"pagination":  pagination.Render(),
	})
}

func UpdateDraftView(c *gin.Context) {
	draftID, err := strconv.ParseUint(c.Param("draftID"), 10, 64)
	if err != nil || draftID <= 0 {
		c.HTML(http.StatusNotFound, "error/404", nil)
		return
	}
	var draft models.Draft
	sqlRes := models.GetDraft(draftID)
	if err = sqlRes.Scan(
		&draft.ID,
		&draft.CategoryID,
		&draft.Title,
		&draft.Content,
		&draft.Publish,
		&draft.CreatedAt,
		&draft.UpdatedAt,
		&draft.Category.Name,
	); err != nil {
		c.HTML(http.StatusNotFound, "error/404", nil)
		return
	}

	var categories []*models.Category
	categoriesRows, err := models.GetCategories("id", "name")
	if err != nil {
		c.Error(err).SetType(http.StatusBadRequest).SetMeta(middleware.BadRequestMeta{
			Url: "/admin/deafts",
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
	c.HTML(http.StatusOK, "admin/editDraft", gin.H{
		"title":       utils.SiteTitle("编辑草稿", siteName),
		"login":       session.Get("login"),
		"menu":        "drafts",
		"draft":       draft,
		"categories":  categories,
		"errorMsgs":   errorMsgs,
		"successMsgs": successMsgs,
	})
}

func AddDraft(c *gin.Context) {
	var draft models.Draft
	if err := c.ShouldBind(&draft); err != nil {
		var msg string
		switch {
		case draft.Title == "":
			msg = "标题不能为空"
		default:
			msg = "内部错误"
		}
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"code": "1",
			"msg":  msg,
		})
		return
	}
	if _, err := models.InsertDraft(&draft); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"code": "1",
			"msg":  "内部错误",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": "0",
		"msg":  "OK",
		"data": gin.H{
			"title":   draft.Title,
			"content": draft.Content,
		},
	})
}

func UpdateDraft(c *gin.Context) {
	var draft models.Draft
	if err := c.ShouldBind(&draft); err != nil {
		var msg string
		switch {
		case draft.Title == "":
			msg = "标题不能为空"
		default:
			msg = err.Error()
		}
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"code": "1",
			"msg":  msg,
		})
		return
	}
	if _, err := models.UpdateDraft(&draft); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"code": "1",
			"msg":  err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": "0",
		"msg":  "OK",
		"data": gin.H{
			"title":   draft.Title,
			"content": draft.Content,
		},
	})
}

func DeleteDraft(c *gin.Context) {
	var draftID uint64
	var err error
	if draftID, err = strconv.ParseUint(c.Param("draftID"), 10, 64); err != nil || draftID < 0 {
		c.HTML(http.StatusNotFound, "error/404", nil)
		return
	}
	session := sessions.Default(c)
	sqlRes, err := models.DeleteDraft(draftID)
	if err != nil {
		session.AddFlash("errorMsg", "参数错误，请重试")
		c.Redirect(http.StatusFound, fmt.Sprintf("/admin/drafts/edit/%d", draftID))
		return
	}

	if affectedRows, err := sqlRes.RowsAffected(); affectedRows <= 0 || err != nil {
		session.AddFlash("errorMsg", "参数错误，请重试")
		c.Redirect(http.StatusFound, fmt.Sprintf("/admin/drafts/edit/%d", draftID))
		return
	}

	c.Redirect(http.StatusFound, "/admin/drafts/")
	return
}

func PublishDraft(c *gin.Context) {
	draftID, err := strconv.ParseUint(c.Param("draftID"), 10, 64)
	if err != nil {
		c.HTML(http.StatusNotFound, "error/404", nil)
		return
	}
	err = models.PublishDraft(draftID)
	session := sessions.Default(c)
	if err != nil {
		session.AddFlash("errorMsgs", "服务器开小差了，请稍后重试")
		c.Redirect(http.StatusFound, "/drafts")
		return
	}
	c.Redirect(http.StatusFound, "/admin/drafts")
}
