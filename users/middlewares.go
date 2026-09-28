package users

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/rifkiandriyanto/golang-gin-boilerplate-realworld/common"
)

// Extract token from Authorization header or query parameter
func extractToken(c *gin.Context) string {
	// Check Authorization header first
	bearerToken := c.GetHeader("Authorization")
	if len(bearerToken) > 6 && strings.ToUpper(bearerToken[0:6]) == "TOKEN " {
		return bearerToken[6:]
	}

	// Check query parameter
	token := c.Query("access_token")
	if token != "" {
		return token
	}

	return ""
}

// A helper to write user_id and user_model to the context
func UpdateContextUserModel(c *gin.Context, my_user_id uint) {
	var myUserModel UserModel
	if my_user_id != 0 {
		db := common.GetDB()
		db.First(&myUserModel, my_user_id)
	}
	c.Set("my_user_id", my_user_id)
	c.Set("my_user_model", myUserModel)
}

// You can custom middlewares yourself as the doc: https://github.com/gin-gonic/gin#custom-middleware
//
//	r.Use(AuthMiddleware(true))
func AuthMiddleware(auto401 bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		UpdateContextUserModel(c, 0)
		tokenString := extractToken(c)

		if tokenString == "" {
			if auto401 {
				c.AbortWithStatus(http.StatusUnauthorized)
			}
			return
		}

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Validate the signing method
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return common.JWTSecret(), nil
		})

		if err != nil {
			if auto401 {
				c.AbortWithStatus(http.StatusUnauthorized)
			}
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			// id is encoded as a float64 by MapClaims; guard against missing/malformed claims
			id, ok := claims["id"].(float64)
			if !ok {
				if auto401 {
					c.AbortWithStatus(http.StatusUnauthorized)
				}
				return
			}
			UpdateContextUserModel(c, uint(id))
		}
	}
}
