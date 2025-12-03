package main

import (
	"golang_system/config"
	"golang_system/controllers"
	"golang_system/database"
	"golang_system/middleware"
	"log"

	_ "golang_system/docs"
	"net/http"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title 博客系统
// @version 0.0.1
// @description  微服务 code  by 田发亮
// @BasePath /
func main() {
	// 加载配置
	cfg := config.Load()

	// 连接数据库
	if err := database.Connect(); err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// 初始化 Gin
	r := gin.Default()

	// 中间件
	r.Use(middleware.CORSMiddleware())
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	// 初始化控制器
	authController := controllers.NewAuthController()
	articleController := controllers.NewArticleController()
	commentController := controllers.NewCommentController()

	// 路由组
	api := r.Group("/api")
	{
		// 认证路由
		auth := api.Group("/auth")
		{
			auth.POST("/register", authController.Register)
			auth.POST("/login", authController.Login)
		}

		// 文章路由
		articles := api.Group("/articles")
		{
			articles.GET("", articleController.GetArticles)
			articles.GET("/:id", articleController.GetArticle)

			// 需要认证的路由
			authRequired := articles.Group("")
			authRequired.Use(middleware.AuthMiddleware())
			{
				authRequired.POST("", articleController.CreateArticle)
				authRequired.PUT("/:id", articleController.UpdateArticle)
				authRequired.DELETE("/:id", articleController.DeleteArticle)
			}
		}

		// 评论路由
		comments := api.Group("/comments")
		comments.Use(middleware.AuthMiddleware())
		{
			comments.POST("/:articleId", commentController.CreateComment)
			comments.GET("/:articleId", commentController.GetArticleComments)
		}
	}
	// http://127.0.0.1:8080/swagger/index.html
	// 仅在开发环境启用
	if gin.Mode() != gin.ReleaseMode {
		r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}
	// 启动服务器
	log.Printf("Server starting on port %s", cfg.Server.Port)
	if err := r.Run(":" + cfg.Server.Port); err != nil {
		log.Fatal("Failed to start server:", err)
	}

}

// 自定义错误处理中间件
func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) > 0 {
			for _, err := range c.Errors {
				log.Printf("Error: %v\n", err)
			}

			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Internal server error",
			})
		}
	}
}
