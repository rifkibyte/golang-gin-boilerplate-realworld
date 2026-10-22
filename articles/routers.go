package articles

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/rifkiandriyanto/golang-gin-boilerplate-realworld/common"
	"github.com/rifkiandriyanto/golang-gin-boilerplate-realworld/users"
	"gorm.io/gorm"
	"net/http"
	"strconv"
)

// ArticleRequest represents the article create/update request body
// swagger:model ArticleRequest
type ArticleRequest struct {
	Article struct {
		Title       string   `json:"title" binding:"required" example:"How to train your dragon"`
		Description string   `json:"description" binding:"required" example:"Ever wonder how?"`
		Body        string   `json:"body" binding:"required" example:"You have to believe"`
		TagList     []string `json:"tagList" example:"dragons,training"`
	} `json:"article"`
}

// SwaggerArticleResponse represents the article response for swagger
// swagger:model ArticleResponse
type SwaggerArticleResponse struct {
	Article struct {
		Slug           string   `json:"slug"`
		Title          string   `json:"title"`
		Description    string   `json:"description"`
		Body           string   `json:"body"`
		TagList        []string `json:"tagList"`
		CreatedAt      string   `json:"createdAt"`
		UpdatedAt      string   `json:"updatedAt"`
		Favorited      bool     `json:"favorited"`
		FavoritesCount int      `json:"favoritesCount"`
		Author         struct {
			Username  string `json:"username"`
			Bio       string `json:"bio"`
			Image     string `json:"image"`
			Following bool   `json:"following"`
		} `json:"author"`
	} `json:"article"`
}

// SwaggerArticlesResponse represents the articles list response for swagger
// swagger:model ArticlesResponse
type SwaggerArticlesResponse struct {
	Articles      []struct {
		Slug           string   `json:"slug"`
		Title          string   `json:"title"`
		Description    string   `json:"description"`
		Body           string   `json:"body"`
		TagList        []string `json:"tagList"`
		CreatedAt      string   `json:"createdAt"`
		UpdatedAt      string   `json:"updatedAt"`
		Favorited      bool     `json:"favorited"`
		FavoritesCount int      `json:"favoritesCount"`
		Author         struct {
			Username  string `json:"username"`
			Bio       string `json:"bio"`
			Image     string `json:"image"`
			Following bool   `json:"following"`
		} `json:"author"`
	} `json:"articles"`
	ArticlesCount int `json:"articlesCount"`
}

// CommentRequest represents the comment create request body
// swagger:model CommentRequest
type CommentRequest struct {
	Comment struct {
		Body string `json:"body" binding:"required" example:"This is a comment"`
	} `json:"comment"`
}

// SwaggerCommentResponse represents the comment response for swagger
// swagger:model CommentResponse
type SwaggerCommentResponse struct {
	Comment struct {
		ID        uint   `json:"id"`
		Body      string `json:"body"`
		CreatedAt string `json:"createdAt"`
		UpdatedAt string `json:"updatedAt"`
		Author    struct {
			Username  string `json:"username"`
			Bio       string `json:"bio"`
			Image     string `json:"image"`
			Following bool   `json:"following"`
		} `json:"author"`
	} `json:"comment"`
}

// SwaggerCommentsResponse represents the comments list response for swagger
// swagger:model CommentsResponse
type SwaggerCommentsResponse struct {
	Comments []struct {
		ID        uint   `json:"id"`
		Body      string `json:"body"`
		CreatedAt string `json:"createdAt"`
		UpdatedAt string `json:"updatedAt"`
		Author    struct {
			Username  string `json:"username"`
			Bio       string `json:"bio"`
			Image     string `json:"image"`
			Following bool   `json:"following"`
		} `json:"author"`
	} `json:"comments"`
}

// SwaggerTagsResponse represents the tags list response for swagger
// swagger:model TagsResponse
type SwaggerTagsResponse struct {
	Tags []string `json:"tags"`
}

// ErrorResponse represents an error response
// swagger:model ErrorResponse
type ErrorResponse struct {
	Errors map[string]interface{} `json:"errors"`
}

// ArticlesRegister godoc
// @Summary Articles endpoints requiring authentication
// @Description Register article routes that require authentication
// @Tags articles
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Router /articles/feed [get]
// @Router /articles [post]
// @Router /articles/{slug} [put]
// @Router /articles/{slug} [delete]
// @Router /articles/{slug}/favorite [post]
// @Router /articles/{slug}/favorite [delete]
// @Router /articles/{slug}/comments [post]
// @Router /articles/{slug}/comments/{id} [delete]
func ArticlesRegister(router *gin.RouterGroup) {
	router.GET("/feed", ArticleFeed)
	router.POST("", ArticleCreate)
	router.POST("/", ArticleCreate)
	router.PUT("/:slug", ArticleUpdate)
	router.PUT("/:slug/", ArticleUpdate)
	router.DELETE("/:slug", ArticleDelete)
	router.POST("/:slug/favorite", ArticleFavorite)
	router.DELETE("/:slug/favorite", ArticleUnfavorite)
	router.POST("/:slug/comments", ArticleCommentCreate)
	router.DELETE("/:slug/comments/:id", ArticleCommentDelete)
}

// ArticlesAnonymousRegister godoc
// @Summary Public articles endpoints
// @Description Register public article routes (no authentication required)
// @Tags articles
// @Accept json
// @Produce json
// @Router /articles [get]
// @Router /articles/{slug} [get]
// @Router /articles/{slug}/comments [get]
func ArticlesAnonymousRegister(router *gin.RouterGroup) {
	router.GET("", ArticleList)
	router.GET("/", ArticleList)
	router.GET("/:slug", ArticleRetrieve)
	router.GET("/:slug/comments", ArticleCommentList)
}

// TagsAnonymousRegister godoc
// @Summary Public tags endpoint
// @Description Register public tags route (no authentication required)
// @Tags tags
// @Accept json
// @Produce json
// @Router /tags [get]
func TagsAnonymousRegister(router *gin.RouterGroup) {
	router.GET("", TagList)
	router.GET("/", TagList)
}

// ArticleCreate godoc
// @Summary Create an article
// @Description Create a new article (authentication required)
// @Tags articles
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param article body ArticleRequest true "Article data"
// @Success 201 {object} SwaggerArticleResponse
// @Failure 401 {object} ErrorResponse
// @Failure 422 {object} ErrorResponse
// @Router /articles [post]
func ArticleCreate(c *gin.Context) {
	articleModelValidator := NewArticleModelValidator()
	if err := articleModelValidator.Bind(c); err != nil {
		c.JSON(http.StatusUnprocessableEntity, common.NewValidatorError(err))
		return
	}
	//fmt.Println(articleModelValidator.articleModel.Author.UserModel)

	if err := SaveOne(&articleModelValidator.articleModel); err != nil {
		c.JSON(http.StatusUnprocessableEntity, common.NewError("database", err))
		return
	}
	serializer := ArticleSerializer{c, articleModelValidator.articleModel}
	c.JSON(http.StatusCreated, gin.H{"article": serializer.Response()})
}

// ArticleList godoc
// @Summary List articles
// @Description Get a list of articles with optional filters
// @Tags articles
// @Accept json
// @Produce json
// @Param tag query string false "Filter by tag"
// @Param author query string false "Filter by author"
// @Param favorited query string false "Filter by favorited user"
// @Param limit query int false "Limit number of articles" default(20)
// @Param offset query int false "Offset for pagination" default(0)
// @Success 200 {object} SwaggerArticlesResponse
// @Failure 404 {object} ErrorResponse
// @Router /articles [get]
func ArticleList(c *gin.Context) {
	//condition := ArticleModel{}
	tag := c.Query("tag")
	author := c.Query("author")
	favorited := c.Query("favorited")
	limit := c.Query("limit")
	offset := c.Query("offset")
	articleModels, modelCount, err := FindManyArticle(tag, author, limit, offset, favorited)
	if err != nil {
		c.JSON(http.StatusNotFound, common.NewError("articles", errors.New("Invalid param")))
		return
	}
	serializer := ArticlesSerializer{c, articleModels}
	c.JSON(http.StatusOK, gin.H{"articles": serializer.Response(), "articlesCount": modelCount})
}

// ArticleFeed godoc
// @Summary Get article feed
// @Description Get articles from users you follow (authentication required)
// @Tags articles
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param limit query int false "Limit number of articles" default(20)
// @Param offset query int false "Offset for pagination" default(0)
// @Success 200 {object} SwaggerArticlesResponse
// @Failure 401 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /articles/feed [get]
func ArticleFeed(c *gin.Context) {
	limit := c.Query("limit")
	offset := c.Query("offset")
	myUserModel := c.MustGet("my_user_model").(users.UserModel)
	if myUserModel.ID == 0 {
		c.JSON(http.StatusUnauthorized, common.NewError("articles", errors.New("Require auth!")))
		return
	}
	articleUserModel := GetArticleUserModel(myUserModel)
	articleModels, modelCount, err := articleUserModel.GetArticleFeed(limit, offset)
	if err != nil {
		c.JSON(http.StatusNotFound, common.NewError("articles", errors.New("Invalid param")))
		return
	}
	serializer := ArticlesSerializer{c, articleModels}
	c.JSON(http.StatusOK, gin.H{"articles": serializer.Response(), "articlesCount": modelCount})
}

// ArticleRetrieve godoc
// @Summary Get an article
// @Description Get a single article by slug
// @Tags articles
// @Accept json
// @Produce json
// @Param slug path string true "Article slug"
// @Success 200 {object} SwaggerArticleResponse
// @Failure 404 {object} ErrorResponse
// @Router /articles/{slug} [get]
func ArticleRetrieve(c *gin.Context) {
	slug := c.Param("slug")
	articleModel, err := FindOneArticle(&ArticleModel{Slug: slug})
	if err != nil {
		c.JSON(http.StatusNotFound, common.NewError("articles", errors.New("Invalid slug")))
		return
	}
	serializer := ArticleSerializer{c, articleModel}
	c.JSON(http.StatusOK, gin.H{"article": serializer.Response()})
}

// ArticleUpdate godoc
// @Summary Update an article
// @Description Update an article by slug (authentication required, must be author)
// @Tags articles
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param slug path string true "Article slug"
// @Param article body ArticleRequest true "Article data"
// @Success 200 {object} SwaggerArticleResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 422 {object} ErrorResponse
// @Router /articles/{slug} [put]
func ArticleUpdate(c *gin.Context) {
	slug := c.Param("slug")
	articleModel, err := FindOneArticle(&ArticleModel{Slug: slug})
	if err != nil {
		c.JSON(http.StatusNotFound, common.NewError("articles", errors.New("Invalid slug")))
		return
	}
	// Check if current user is the author
	myUserModel := c.MustGet("my_user_model").(users.UserModel)
	articleUserModel := GetArticleUserModel(myUserModel)
	if articleModel.AuthorID != articleUserModel.ID {
		c.JSON(http.StatusForbidden, common.NewError("article", errors.New("you are not the author")))
		return
	}

	articleModelValidator := NewArticleModelValidatorFillWith(articleModel)
	if err := articleModelValidator.Bind(c); err != nil {
		c.JSON(http.StatusUnprocessableEntity, common.NewValidatorError(err))
		return
	}

	articleModelValidator.articleModel.ID = articleModel.ID
	if err := articleModel.Update(articleModelValidator.articleModel); err != nil {
		c.JSON(http.StatusUnprocessableEntity, common.NewError("database", err))
		return
	}
	serializer := ArticleSerializer{c, articleModel}
	c.JSON(http.StatusOK, gin.H{"article": serializer.Response()})
}

// ArticleDelete godoc
// @Summary Delete an article
// @Description Delete an article by slug (authentication required, must be author)
// @Tags articles
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param slug path string true "Article slug"
// @Success 200 {object} map[string]string
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 422 {object} ErrorResponse
// @Router /articles/{slug} [delete]
func ArticleDelete(c *gin.Context) {
	slug := c.Param("slug")
	articleModel, err := FindOneArticle(&ArticleModel{Slug: slug})
	if err == nil {
		// Article exists, check authorization
		myUserModel := c.MustGet("my_user_model").(users.UserModel)
		articleUserModel := GetArticleUserModel(myUserModel)
		if articleModel.AuthorID != articleUserModel.ID {
			c.JSON(http.StatusForbidden, common.NewError("article", errors.New("you are not the author")))
			return
		}
	}
	// Delete regardless of existence (idempotent)
	if err := DeleteArticleModel(&ArticleModel{Slug: slug}); err != nil {
		c.JSON(http.StatusUnprocessableEntity, common.NewError("database", err))
		return
	}
	c.JSON(http.StatusOK, gin.H{"article": "delete success"})
}

// ArticleFavorite godoc
// @Summary Favorite an article
// @Description Favorite an article by slug (authentication required)
// @Tags articles
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param slug path string true "Article slug"
// @Success 200 {object} SwaggerArticleResponse
// @Failure 401 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 422 {object} ErrorResponse
// @Router /articles/{slug}/favorite [post]
func ArticleFavorite(c *gin.Context) {
	slug := c.Param("slug")
	articleModel, err := FindOneArticle(&ArticleModel{Slug: slug})
	if err != nil {
		c.JSON(http.StatusNotFound, common.NewError("articles", errors.New("Invalid slug")))
		return
	}
	myUserModel := c.MustGet("my_user_model").(users.UserModel)
	if err = articleModel.favoriteBy(GetArticleUserModel(myUserModel)); err != nil {
		c.JSON(http.StatusUnprocessableEntity, common.NewError("database", err))
		return
	}
	serializer := ArticleSerializer{c, articleModel}
	c.JSON(http.StatusOK, gin.H{"article": serializer.Response()})
}

// ArticleUnfavorite godoc
// @Summary Unfavorite an article
// @Description Unfavorite an article by slug (authentication required)
// @Tags articles
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param slug path string true "Article slug"
// @Success 200 {object} SwaggerArticleResponse
// @Failure 401 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 422 {object} ErrorResponse
// @Router /articles/{slug}/favorite [delete]
func ArticleUnfavorite(c *gin.Context) {
	slug := c.Param("slug")
	articleModel, err := FindOneArticle(&ArticleModel{Slug: slug})
	if err != nil {
		c.JSON(http.StatusNotFound, common.NewError("articles", errors.New("Invalid slug")))
		return
	}
	myUserModel := c.MustGet("my_user_model").(users.UserModel)
	if err = articleModel.unFavoriteBy(GetArticleUserModel(myUserModel)); err != nil {
		c.JSON(http.StatusUnprocessableEntity, common.NewError("database", err))
		return
	}
	serializer := ArticleSerializer{c, articleModel}
	c.JSON(http.StatusOK, gin.H{"article": serializer.Response()})
}

// ArticleCommentCreate godoc
// @Summary Add a comment to an article
// @Description Add a comment to an article by slug (authentication required)
// @Tags comments
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param slug path string true "Article slug"
// @Param comment body CommentRequest true "Comment data"
// @Success 201 {object} SwaggerCommentResponse
// @Failure 401 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 422 {object} ErrorResponse
// @Router /articles/{slug}/comments [post]
func ArticleCommentCreate(c *gin.Context) {
	slug := c.Param("slug")
	articleModel, err := FindOneArticle(&ArticleModel{Slug: slug})
	if err != nil {
		c.JSON(http.StatusNotFound, common.NewError("comment", errors.New("Invalid slug")))
		return
	}
	commentModelValidator := NewCommentModelValidator()
	if err := commentModelValidator.Bind(c); err != nil {
		c.JSON(http.StatusUnprocessableEntity, common.NewValidatorError(err))
		return
	}
	commentModelValidator.commentModel.Article = articleModel

	if err := SaveOne(&commentModelValidator.commentModel); err != nil {
		c.JSON(http.StatusUnprocessableEntity, common.NewError("database", err))
		return
	}
	serializer := CommentSerializer{c, commentModelValidator.commentModel}
	c.JSON(http.StatusCreated, gin.H{"comment": serializer.Response()})
}

// ArticleCommentDelete godoc
// @Summary Delete a comment
// @Description Delete a comment by ID (authentication required, must be author)
// @Tags comments
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param slug path string true "Article slug"
// @Param id path int true "Comment ID"
// @Success 200 {object} map[string]string
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 422 {object} ErrorResponse
// @Router /articles/{slug}/comments/{id} [delete]
func ArticleCommentDelete(c *gin.Context) {
	id64, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusNotFound, common.NewError("comment", errors.New("Invalid id")))
		return
	}
	id := uint(id64)
	commentModel, err := FindOneComment(&CommentModel{Model: gorm.Model{ID: id}})
	if err == nil {
		// Comment exists, check authorization
		myUserModel := c.MustGet("my_user_model").(users.UserModel)
		articleUserModel := GetArticleUserModel(myUserModel)
		if commentModel.AuthorID != articleUserModel.ID {
			c.JSON(http.StatusForbidden, common.NewError("comment", errors.New("you are not the author")))
			return
		}
	}
	// Delete regardless of existence (idempotent)
	if err := DeleteCommentModel([]uint{id}); err != nil {
		c.JSON(http.StatusUnprocessableEntity, common.NewError("database", err))
		return
	}
	c.JSON(http.StatusOK, gin.H{"comment": "delete success"})
}

// ArticleCommentList godoc
// @Summary Get comments for an article
// @Description Get all comments for an article by slug
// @Tags comments
// @Accept json
// @Produce json
// @Param slug path string true "Article slug"
// @Success 200 {object} SwaggerCommentsResponse
// @Failure 404 {object} ErrorResponse
// @Router /articles/{slug}/comments [get]
func ArticleCommentList(c *gin.Context) {
	slug := c.Param("slug")
	articleModel, err := FindOneArticle(&ArticleModel{Slug: slug})
	if err != nil {
		c.JSON(http.StatusNotFound, common.NewError("comments", errors.New("Invalid slug")))
		return
	}
	err = articleModel.getComments()
	if err != nil {
		c.JSON(http.StatusNotFound, common.NewError("comments", errors.New("Database error")))
		return
	}
	serializer := CommentsSerializer{c, articleModel.Comments}
	c.JSON(http.StatusOK, gin.H{"comments": serializer.Response()})
}

// TagList godoc
// @Summary Get all tags
// @Description Get a list of all tags
// @Tags tags
// @Accept json
// @Produce json
// @Success 200 {object} SwaggerTagsResponse
// @Failure 404 {object} ErrorResponse
// @Router /tags [get]
func TagList(c *gin.Context) {
	tagModels, err := getAllTags()
	if err != nil {
		c.JSON(http.StatusNotFound, common.NewError("articles", errors.New("Invalid param")))
		return
	}
	serializer := TagsSerializer{c, tagModels}
	c.JSON(http.StatusOK, gin.H{"tags": serializer.Response()})
}
