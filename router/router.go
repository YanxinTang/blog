package router

import (
	"html/template"

	"github.com/YanxinTang/blog/controllers"
	"github.com/YanxinTang/blog/middleware"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	store := cookie.NewStore([]byte("secret"))
	r.Use(sessions.Sessions("sessionid", store))
	r.Use(middleware.Method(r), middleware.ErrorHandler())
	r.SetFuncMap(template.FuncMap{
		"Date":     Date,
		"Safe":     Safe,
		"Markdown": Markdown,
		"Summary":  Summary,
		"Config":   Config,
	})
	r.LoadHTMLGlob("dist/views/*.tmpl")
	r.Static("static", "dist/static")

	public := r.Group("")
	{
		public.GET("/", controllers.Articles)
		public.GET("/page/:page", controllers.Articles)
		public.GET("/login/", controllers.LoginView)
		public.POST("/login", controllers.Login)
		public.POST("/logout", controllers.Logout)
		public.GET("/articles/:articleID", controllers.ArticleView)
		public.POST("/articles/:articleID/comment", controllers.AddComment)
	}

	protected := r.Group("")
	{
		protected.Use(middleware.Auth())
		// protected.GET("/app", controllers.App)
		// protected.GET("/app/*route", controllers.App)
		protected.DELETE("/articles/:articleID", controllers.DeleteArticle)
		protected.DELETE("/articles/:articleID/comment/:commentID", controllers.DeleteComment)
		admin := protected.Group("admin")
		{
			admin.GET("/", controllers.DashBoardView)
			admin.GET("/articles/new", controllers.AddArticleView)
			admin.POST("/articles/new", controllers.AddArticle)
			admin.GET("/articles/update/:articleID", controllers.UpdateArticleView)
			admin.PUT("/articles/update/:articleID", controllers.UpdateArticle)
			admin.GET("/categories", controllers.CategoriesView)
			admin.POST("/categories", controllers.AddCategory)
			admin.DELETE("/categories/:categoryID", controllers.DeleteCategory)
			admin.GET("/drafts/", controllers.DraftsView)
			admin.GET("/drafts/page/:page", controllers.DraftsView)
			admin.GET("/drafts/edit/:draftID", controllers.UpdateDraftView)
			admin.DELETE("/drafts/:draftID", controllers.DeleteDraft)
			admin.PUT("/drafts/publish/:draftID", controllers.PublishDraft)
		}
	}

	api := r.Group("api")
	{
		admin := api.Group("admin")
		{
			admin.Use(middleware.AuthAPI())
			admin.PUT("/drafts/", controllers.AddDraft)
			admin.PATCH("/drafts/", controllers.UpdateDraft)
		}
	}

	return r
}
