package users

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/rifkibyte/golang-gin-boilerplate-realworld/common"
	"net/http"
)

// UserRegistrationRequest represents the registration request body
// swagger:model UserRegistrationRequest
type UserRegistrationRequest struct {
	User struct {
		Username string `json:"username" binding:"required" example:"jake"`
		Email    string `json:"email" binding:"required,email" example:"jake@jake.jake"`
		Password string `json:"password" binding:"required" example:"jakejake"`
	} `json:"user"`
}

// UserLoginRequest represents the login request body
// swagger:model UserLoginRequest
type UserLoginRequest struct {
	User struct {
		Email    string `json:"email" binding:"required,email" example:"jake@jake.jake"`
		Password string `json:"password" binding:"required" example:"jakejake"`
	} `json:"user"`
}

// SwaggerUserResponse represents the user response
// swagger:model UserResponse
type SwaggerUserResponse struct {
	User struct {
		Email    string `json:"email"`
		Token    string `json:"token"`
		Username string `json:"username"`
		Bio      string `json:"bio"`
		Image    string `json:"image"`
	} `json:"user"`
}

// SwaggerProfileResponse represents the profile response
// swagger:model ProfileResponse
type SwaggerProfileResponse struct {
	Profile struct {
		Username  string `json:"username"`
		Bio       string `json:"bio"`
		Image     string `json:"image"`
		Following bool   `json:"following"`
	} `json:"profile"`
}

// UserUpdateRequest represents the user update request body
// swagger:model UserUpdateRequest
type UserUpdateRequest struct {
	User struct {
		Email    string `json:"email" example:"jake@jake.jake"`
		Username string `json:"username" example:"jake"`
		Password string `json:"password" example:"newpassword"`
		Bio      string `json:"bio" example:"I like to skateboard"`
		Image    string `json:"image" example:"https://example.com/image.png"`
	} `json:"user"`
}

// ErrorResponse represents an error response
// swagger:model ErrorResponse
type ErrorResponse struct {
	Errors map[string]interface{} `json:"errors"`
}

// Register godoc
// @Summary Register a new user
// @Description Register a new user account
// @Tags users
// @Accept json
// @Produce json
// @Param user body UserRegistrationRequest true "User registration info"
// @Success 201 {object} SwaggerUserResponse
// @Failure 422 {object} ErrorResponse
// @Router /users [post]
func UsersRegister(router *gin.RouterGroup) {
	router.POST("", UsersRegistration)
	router.POST("/", UsersRegistration)
	router.POST("/login", UsersLogin)
}

// UserRegister godoc
// @Summary Get/Update current user
// @Description Get or update the currently authenticated user
// @Tags users
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} SwaggerUserResponse
// @Failure 401 {object} ErrorResponse
// @Failure 422 {object} ErrorResponse
// @Router /user [get]
func UserRegister(router *gin.RouterGroup) {
	router.GET("", UserRetrieve)
	router.GET("/", UserRetrieve)
	router.PUT("", UserUpdate)
	router.PUT("/", UserUpdate)
}

// ProfileRetrieveRegister godoc
// @Summary Get profile by username
// @Description Get a user's profile by username
// @Tags profiles
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param username path string true "Username"
// @Success 200 {object} SwaggerProfileResponse
// @Failure 404 {object} ErrorResponse
// @Router /profiles/{username} [get]
func ProfileRetrieveRegister(router *gin.RouterGroup) {
	router.GET("/:username", ProfileRetrieve)
}

// ProfileRegister godoc
// @Summary Follow/Unfollow a user
// @Description Follow or unfollow a user by username
// @Tags profiles
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param username path string true "Username"
// @Success 200 {object} SwaggerProfileResponse
// @Failure 401 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 422 {object} ErrorResponse
// @Router /profiles/{username}/follow [post]
// @Router /profiles/{username}/follow [delete]
func ProfileRegister(router *gin.RouterGroup) {
	router.POST("/:username/follow", ProfileFollow)
	router.DELETE("/:username/follow", ProfileUnfollow)
}

// ProfileRetrieve godoc
// @Summary Get profile by username
// @Description Get a user's profile by username
// @Tags profiles
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param username path string true "Username"
// @Success 200 {object} SwaggerProfileResponse
// @Failure 404 {object} ErrorResponse
// @Router /profiles/{username} [get]
func ProfileRetrieve(c *gin.Context) {
	username := c.Param("username")
	userModel, err := FindOneUser(&UserModel{Username: username})
	if err != nil {
		c.JSON(http.StatusNotFound, common.NewError("profile", errors.New("Invalid username")))
		return
	}
	profileSerializer := ProfileSerializer{c, userModel}
	c.JSON(http.StatusOK, gin.H{"profile": profileSerializer.Response()})
}

// ProfileFollow godoc
// @Summary Follow a user
// @Description Follow a user by username
// @Tags profiles
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param username path string true "Username"
// @Success 200 {object} SwaggerProfileResponse
// @Failure 401 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 422 {object} ErrorResponse
// @Router /profiles/{username}/follow [post]
func ProfileFollow(c *gin.Context) {
	username := c.Param("username")
	userModel, err := FindOneUser(&UserModel{Username: username})
	if err != nil {
		c.JSON(http.StatusNotFound, common.NewError("profile", errors.New("Invalid username")))
		return
	}
	myUserModel := c.MustGet("my_user_model").(UserModel)
	err = myUserModel.following(userModel)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, common.NewError("database", err))
		return
	}
	serializer := ProfileSerializer{c, userModel}
	c.JSON(http.StatusOK, gin.H{"profile": serializer.Response()})
}

// ProfileUnfollow godoc
// @Summary Unfollow a user
// @Description Unfollow a user by username
// @Tags profiles
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param username path string true "Username"
// @Success 200 {object} SwaggerProfileResponse
// @Failure 401 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 422 {object} ErrorResponse
// @Router /profiles/{username}/follow [delete]
func ProfileUnfollow(c *gin.Context) {
	username := c.Param("username")
	userModel, err := FindOneUser(&UserModel{Username: username})
	if err != nil {
		c.JSON(http.StatusNotFound, common.NewError("profile", errors.New("Invalid username")))
		return
	}
	myUserModel := c.MustGet("my_user_model").(UserModel)

	err = myUserModel.unFollowing(userModel)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, common.NewError("database", err))
		return
	}
	serializer := ProfileSerializer{c, userModel}
	c.JSON(http.StatusOK, gin.H{"profile": serializer.Response()})
}

// UsersRegistration godoc
// @Summary Register a new user
// @Description Register a new user account
// @Tags users
// @Accept json
// @Produce json
// @Param user body UserRegistrationRequest true "User registration info"
// @Success 201 {object} SwaggerUserResponse
// @Failure 422 {object} ErrorResponse
// @Router /users [post]
func UsersRegistration(c *gin.Context) {
	userModelValidator := NewUserModelValidator()
	if err := userModelValidator.Bind(c); err != nil {
		c.JSON(http.StatusUnprocessableEntity, common.NewValidatorError(err))
		return
	}

	if err := SaveOne(&userModelValidator.userModel); err != nil {
		c.JSON(http.StatusUnprocessableEntity, common.NewError("database", err))
		return
	}
	c.Set("my_user_model", userModelValidator.userModel)
	serializer := UserSerializer{c}
	c.JSON(http.StatusCreated, gin.H{"user": serializer.Response()})
}

// UsersLogin godoc
// @Summary Login user
// @Description Login with email and password
// @Tags users
// @Accept json
// @Produce json
// @Param user body UserLoginRequest true "User login credentials"
// @Success 200 {object} SwaggerUserResponse
// @Failure 401 {object} ErrorResponse
// @Failure 422 {object} ErrorResponse
// @Router /users/login [post]
func UsersLogin(c *gin.Context) {
	loginValidator := NewLoginValidator()
	if err := loginValidator.Bind(c); err != nil {
		c.JSON(http.StatusUnprocessableEntity, common.NewValidatorError(err))
		return
	}
	userModel, err := FindOneUser(&UserModel{Email: loginValidator.userModel.Email})

	if err != nil {
		c.JSON(http.StatusUnauthorized, common.NewError("login", errors.New("Not Registered email or invalid password")))
		return
	}

	if userModel.checkPassword(loginValidator.User.Password) != nil {
		c.JSON(http.StatusUnauthorized, common.NewError("login", errors.New("Not Registered email or invalid password")))
		return
	}
	UpdateContextUserModel(c, userModel.ID)
	serializer := UserSerializer{c}
	c.JSON(http.StatusOK, gin.H{"user": serializer.Response()})
}

// UserRetrieve godoc
// @Summary Get current user
// @Description Get the currently authenticated user
// @Tags users
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} SwaggerUserResponse
// @Failure 401 {object} ErrorResponse
// @Router /user [get]
func UserRetrieve(c *gin.Context) {
	serializer := UserSerializer{c}
	c.JSON(http.StatusOK, gin.H{"user": serializer.Response()})
}

// UserUpdate godoc
// @Summary Update current user
// @Description Update the currently authenticated user
// @Tags users
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param user body UserUpdateRequest true "User update info"
// @Success 200 {object} SwaggerUserResponse
// @Failure 401 {object} ErrorResponse
// @Failure 422 {object} ErrorResponse
// @Router /user [put]
func UserUpdate(c *gin.Context) {
	myUserModel := c.MustGet("my_user_model").(UserModel)
	userModelValidator := NewUserModelValidatorFillWith(myUserModel)
	if err := userModelValidator.Bind(c); err != nil {
		c.JSON(http.StatusUnprocessableEntity, common.NewValidatorError(err))
		return
	}

	userModelValidator.userModel.ID = myUserModel.ID
	if err := myUserModel.Update(userModelValidator.userModel); err != nil {
		c.JSON(http.StatusUnprocessableEntity, common.NewError("database", err))
		return
	}
	UpdateContextUserModel(c, myUserModel.ID)
	serializer := UserSerializer{c}
	c.JSON(http.StatusOK, gin.H{"user": serializer.Response()})
}
