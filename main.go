package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"github.com/rifkiandriyanto/golang-gin-boilerplate-realworld/articles"
	"github.com/rifkiandriyanto/golang-gin-boilerplate-realworld/common"
	"github.com/rifkiandriyanto/golang-gin-boilerplate-realworld/users"
	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) {
	users.AutoMigrate()
	db.AutoMigrate(&articles.ArticleModel{})
	db.AutoMigrate(&articles.TagModel{})
	db.AutoMigrate(&articles.FavoriteModel{})
	db.AutoMigrate(&articles.ArticleUserModel{})
	db.AutoMigrate(&articles.CommentModel{})
}

func pingHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "pong"})
}

func main() {
	// Load .env file if present (real env vars take precedence)
	if err := godotenv.Load(); err != nil && !os.IsNotExist(err) {
		log.Printf("warning: could not load .env file: %v", err)
	}

	if err := common.CheckSecret(); err != nil {
		log.Fatal(err)
	}

	db := common.Init()
	Migrate(db)
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatal("failed to get sql.DB:", err)
	}
	defer sqlDB.Close()

	r := gin.Default()

	// Disable automatic redirect for trailing slashes
	// This prevents POST body from being lost during redirects
	r.RedirectTrailingSlash = false

	v1 := r.Group("/api")
	users.UsersRegister(v1.Group("/users"))
	v1.Use(users.AuthMiddleware(false))
	articles.ArticlesAnonymousRegister(v1.Group("/articles"))
	articles.TagsAnonymousRegister(v1.Group("/tags"))
	users.ProfileRetrieveRegister(v1.Group("/profiles"))

	v1.Use(users.AuthMiddleware(true))
	users.UserRegister(v1.Group("/user"))
	users.ProfileRegister(v1.Group("/profiles"))

	articles.ArticlesRegister(v1.Group("/articles"))

	ping := r.Group("/api/ping")

	ping.GET("", pingHandler)
	ping.GET("/", pingHandler)

	// Get port from environment variable or use default
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	if err := r.Run(":" + port); err != nil {
		log.Fatal("failed to start server:", err)
	}
}
