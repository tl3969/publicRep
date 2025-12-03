package controllers

import (
	"golang_system/database"
	"golang_system/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ArticleController struct {
	DB *gorm.DB
}

func NewArticleController() *ArticleController {
	return &ArticleController{DB: database.DB}
}

// @CreateArticle
// @Description CreateArticle
// @Accept  json
// @Produce json
// @Param Authorization header string true "Bearer token" default(Bearer )
// @Param   req      body   models.CreateArticleRequest     true        "req"
// @Success 200 {string} string	"name,helloWorld"
// @Router /api/articles [post]
func (ac *ArticleController) CreateArticle(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)

	var req models.CreateArticleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	article := models.Article{
		Title:    req.Title,
		Content:  req.Content,
		AuthorID: userID,
	}

	if err := ac.DB.Create(&article).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create article"})
		return
	}

	// 加载作者信息
	ac.DB.Preload("Author").First(&article, article.ID)

	c.JSON(http.StatusCreated, gin.H{
		"message": "Article created successfully",
		"article": article,
	})
}

// @GetArticles
// @Description GetArticles
// @Accept  json
// @Produce json
// @Param Authorization header string true "Bearer token" default(Bearer )
// @Success 200 {string} string	"name,helloWorld"
// @Router /api/articles [get]
func (ac *ArticleController) GetArticles(c *gin.Context) {
	var articles []models.Article

	if err := ac.DB.Preload("Author").Order("created_at DESC").Find(&articles).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch articles"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"articles": articles})
}

// @GetArticle
// @Description GetArticle
// @Accept  json
// @Produce json
// @Param Authorization header string true "Bearer token" default(Bearer )
// @Param   id     path    int     true        "id"
// @Success 200 {string} string	"name,helloWorld"
// @Router /api/articles/{id} [get]
func (ac *ArticleController) GetArticle(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid article ID"})
		return
	}

	var article models.Article
	if err := ac.DB.Preload("Author").Preload("Comments.User").First(&article, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Article not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch article"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"article": article})
}

// @UpdateArticle
// @Description UpdateArticle
// @Accept  json
// @Produce json
// @Param Authorization header string true "Bearer token" default(Bearer )
// @Param   id     path    int     true        "id"
// @Param   req     body    models.UpdateArticleRequest     true        "UpdateArticleRequest"
// @Success 200 {string} string	"name,helloWorld"
// @Router /api/articles/{id} [put]
func (ac *ArticleController) UpdateArticle(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)
	articleID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid article ID"})
		return
	}

	var article models.Article
	if err := ac.DB.First(&article, articleID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Article not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch article"})
		return
	}

	// 检查权限
	if article.AuthorID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "You can only update your own articles"})
		return
	}

	var req models.UpdateArticleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.Title != "" {
		article.Title = req.Title
	}
	if req.Content != "" {
		article.Content = req.Content
	}

	if err := ac.DB.Save(&article).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update article"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Article updated successfully",
		"article": article,
	})
}

// @DeleteArticle
// @Description DeleteArticle
// @Accept  json
// @Produce json
// @Param Authorization header string true "Bearer token" default(Bearer )
// @Param   id     path    int     true        "id"
// @Success 200 {string} string	"name,helloWorld"
// @Router /api/articles/{id} [delete]
func (ac *ArticleController) DeleteArticle(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)
	articleID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid article ID"})
		return
	}

	var article models.Article
	if err := ac.DB.First(&article, articleID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Article not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch article"})
		return
	}

	// 检查权限
	if article.AuthorID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "You can only delete your own articles"})
		return
	}

	// 删除相关评论
	ac.DB.Where("article_id = ?", articleID).Delete(&models.Comment{})

	// 删除文章
	if err := ac.DB.Delete(&article).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete article"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Article deleted successfully"})
}
