package controllers

import (
	"net/http"
	"pubcore/database"
	"pubcore/helpers"
	"pubcore/models"

	"github.com/gin-gonic/gin"
)

var appJSON = "application/json"

// UserRegister godoc
// @Tags Auth
// @Summary Register a new user
// @Description Create a new user account
// @Accept json
// @Produce json
// @Param user body helpers.RegisterInput true "User Registration Input"
// @Success 201 {object} helpers.CommonResponse "Created"
// @Failure 400 {object} helpers.CommonErrorResponse "Bad Request"
// @Router /auth/register [post]
func UserRegister(c *gin.Context) {
	db := database.GetDB()
	contentType := helpers.GetContentType(c)
	User := models.Users{}

	if contentType == appJSON {
		c.ShouldBindJSON(&User)
	} else {
		c.ShouldBind(&User)
	}

	err := db.Create(&User).Error
	if err != nil {
		c.JSON(http.StatusBadRequest, helpers.CommonErrorResponse{
			StatusCode: http.StatusBadRequest,
			Error:      "Bad Request",
			Message:    err.Error(),
		})
		return
	}

	token := helpers.GenerateToken(User.ID, User.Email)

	c.JSON(http.StatusCreated, helpers.CommonResponse{
		StatusCode: http.StatusCreated,
		Message:    "User registered successfully",
		Data: gin.H{
			"token": token,
			"user": gin.H{
				"id":         User.ID,
				"username":   User.Username,
				"email":      User.Email,
				"bio":        User.Bio,
				"created_at": User.CreatedAt,
				"updated_at": User.UpdatedAt,
			},
		},
	})
}

// UserLogin godoc
// @Tags Auth
// @Summary Login with email and password
// @Description Authenticate user and return a JWT token
// @Accept json
// @Produce json
// @Param user body helpers.LoginInput true "User Login Input"
// @Success 200 {object} helpers.CommonResponse "OK"
// @Failure 401 {object} helpers.CommonErrorResponse "Unauthorized"
// @Router /auth/login [post]
func UserLogin(c *gin.Context) {
	db := database.GetDB()
	input := helpers.LoginInput{}

	if err := c.ShouldBindJSON(&input); err != nil {
		// Fallback for non-JSON content types if needed, though login is usually JSON
		if errBind := c.ShouldBind(&input); errBind != nil {
			c.JSON(http.StatusBadRequest, helpers.CommonErrorResponse{
				StatusCode: http.StatusBadRequest,
				Error:      "Bad Request",
				Message:    errBind.Error(),
			})
			return
		}
	}

	User := models.Users{}
	err := db.Where("email = ?", input.Email).Take(&User).Error
	if err != nil {
		c.JSON(http.StatusUnauthorized, helpers.CommonErrorResponse{
			StatusCode: http.StatusUnauthorized,
			Error:      "Unauthorized",
			Message:    "Invalid email/password",
		})
		return
	}

	comparePass := helpers.CompatePass([]byte(User.Password), []byte(input.Password))
	if !comparePass {
		c.JSON(http.StatusUnauthorized, helpers.CommonErrorResponse{
			StatusCode: http.StatusUnauthorized,
			Error:      "Unauthorized",
			Message:    "Invalid email/password",
		})
		return
	}

	token := helpers.GenerateToken(User.ID, User.Email)

	c.JSON(http.StatusOK, helpers.CommonResponse{
		StatusCode: http.StatusOK,
		Message:    "Login successful",
		Data: gin.H{
			"token": token,
			"user": gin.H{
				"id":         User.ID,
				"username":   User.Username,
				"email":      User.Email,
				"bio":        User.Bio,
				"created_at": User.CreatedAt,
				"updated_at": User.UpdatedAt,
			},
		},
	})
}

// UserLogout godoc
// @Tags Auth
// @Summary Logout the current user
// @Description Invalidate the session (client-side token removal)
// @Security BearerAuth
// @Produce json
// @Success 200 {object} helpers.CommonResponse "OK"
// @Failure 401 {object} helpers.CommonErrorResponse "Unauthorized"
// @Router /auth/logout [post]
func UserLogout(c *gin.Context) {
	c.JSON(http.StatusOK, helpers.CommonResponse{
		StatusCode: http.StatusOK,
		Message:    "Successfully logged out. Please remove your token on the client side.",
	})
}
