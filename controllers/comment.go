package controllers

import (
	"golang_system/database"
	"golang_system/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type CommentController struct {
	DB *gorm.DB
}

func NewCommentController() *CommentController {
	return &CommentController{DB: database.DB}
}

// @CreateComment
// @Description CreateComment
// @Accept  json
// @Produce json
// @Param Authorization header string true "Bearer token" default(Bearer )
// @Param   articleId     path    int     true        "articleId"
// @Param   req     body    models.CreateCommentRequest    true        "req"
// @Success 200 {string} string	"name,helloWorld"
// @Router /api/comments/{articleId} [post]
func (cc *CommentController) CreateComment(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)
	articleID, err := strconv.Atoi(c.Param("articleId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid article ID"})
		return
	}

	// 检查文章是否存在
	var article models.Article
	if err := cc.DB.First(&article, articleID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Article not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch article"})
		return
	}

	var req models.CreateCommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	comment := models.Comment{
		Content:   req.Content,
		ArticleID: uint(articleID),
		UserID:    userID,
	}

	if err := cc.DB.Create(&comment).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create comment"})
		return
	}

	// 加载用户信息
	cc.DB.Preload("User").First(&comment, comment.ID)

	c.JSON(http.StatusCreated, gin.H{
		"message": "Comment created successfully",
		"comment": comment,
	})
}

// @GetArticleComments
// @Description GetArticleComments
// @Accept  json
// @Produce json
// @Param Authorization header string true "Bearer token" default(Bearer )
// @Param   articleId     path    int     true        "articleId"
// @Success 200 {string} string	"name,helloWorld"
// @Router /api/comments/{articleId} [get]
func (cc *CommentController) GetArticleComments(c *gin.Context) {
	articleID, err := strconv.Atoi(c.Param("articleId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid article ID"})
		return
	}

	var comments []models.Comment
	if err := cc.DB.Preload("User").Where("article_id = ?", articleID).Order("created_at DESC").Find(&comments).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch comments"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"comments": comments})
}
