// Common tools and helper functions
package common

import (
	"crypto/rand"
	"errors"
	"fmt"
	"math/big"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v5"
)

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

// A helper function to generate random string
func RandString(n int) string {
	b := make([]rune, n)
	for i := range b {
		randIdx, err := rand.Int(rand.Reader, big.NewInt(int64(len(letters))))
		if err != nil {
			panic(err)
		}
		b[i] = letters[randIdx.Int64()]
	}
	return string(b)
}

// A helper function to generate random int
func RandInt() int {
	randNum, err := rand.Int(rand.Reader, big.NewInt(1000000))
	if err != nil {
		panic(err)
	}
	return int(randNum.Int64())
}

// Keep this two config private, it should not expose to open source
const RandomPassword = "A String Very Very Very Random!!@##$!@#4" // #nosec G101

// devJWTSecret is only used as a fallback for local development.
// Always set JWT_SECRET in production via environment variables or .env.
const devJWTSecret = "dev-only-secret-change-me-in-production" // #nosec G101

// JWTSecret returns the JWT signing secret. It reads from the JWT_SECRET
// environment variable and falls back to a development-only value.
func JWTSecret() []byte {
	if s := os.Getenv("JWT_SECRET"); s != "" {
		return []byte(s)
	}
	return []byte(devJWTSecret)
}

// CheckSecret fails fast when the development secret would be used in
// production mode.
func CheckSecret() error {
	if os.Getenv("GIN_MODE") == gin.ReleaseMode && os.Getenv("JWT_SECRET") == "" {
		return errors.New("JWT_SECRET must be set when running in release mode")
	}
	return nil
}

// A Util function to generate jwt_token which can be used in the request header
func GenToken(id uint) string {
	jwt_token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":  id,
		"exp": time.Now().Add(time.Hour * 24).Unix(),
	})
	// Sign and get the complete encoded token as a string
	token, err := jwt_token.SignedString(JWTSecret())
	if err != nil {
		fmt.Printf("failed to sign JWT token for id %d: %v\n", id, err)
		return ""
	}
	return token
}

// My own Error type that will help return my customized Error info
//
//	{"database": {"hello":"no such table", error: "not_exists"}}
type CommonError struct {
	Errors map[string]interface{} `json:"errors"`
}

// To handle the error returned by c.Bind in gin framework
// https://github.com/go-playground/validator/blob/v9/_examples/translations/main.go
func NewValidatorError(err error) CommonError {
	res := CommonError{}
	res.Errors = make(map[string]interface{})
	errs := err.(validator.ValidationErrors)
	for _, v := range errs {
		// can translate each error one at a time.
		//fmt.Println("gg",v.NameNamespace)
		if v.Param() != "" {
			res.Errors[v.Field()] = fmt.Sprintf("{%v: %v}", v.Tag(), v.Param())
		} else {
			res.Errors[v.Field()] = fmt.Sprintf("{key: %v}", v.Tag())
		}

	}
	return res
}

// Wrap the error info in an object
func NewError(key string, err error) CommonError {
	res := CommonError{}
	res.Errors = make(map[string]interface{})
	res.Errors[key] = err.Error()
	return res
}

// Changed the c.MustBindWith() ->  c.ShouldBindWith().
// I don't want to auto return 400 when error happened.
// origin function is here: https://github.com/gin-gonic/gin/blob/master/context.go
func Bind(c *gin.Context, obj interface{}) error {
	b := binding.Default(c.Request.Method, c.ContentType())
	return c.ShouldBindWith(obj, b)
}
