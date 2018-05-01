package presentation

import (
	"net/http"
	"net/mail"

	"amalhanaja.com/user-service/data"
	"amalhanaja.com/user-service/domain"
	"github.com/gin-gonic/gin"
)

// RegisterController ...
func RegisterController(c *gin.Context) {
	var jsonParams Register
	if err := c.ShouldBindJSON(&jsonParams); err == nil {
		useCase := domain.NewUserUseCase(data.GetInstance())
		var errorModel *ErrorModel
		if len(jsonParams.FullName) < 8 {
			errorModel = &ErrorModel{
				4001,
				"Full Name must be at least 8 characters",
			}
		} else if _, err := mail.ParseAddress(jsonParams.Email); err != nil {
			errorModel = &ErrorModel{
				4002,
				"Invalid Email Address",
			}
		}
		if errorModel != nil {
			c.JSON(http.StatusBadRequest, errorModel)
		} else {
			if token, err := useCase.Register(domain.NewUser{
				Email:    jsonParams.Email,
				FullName: jsonParams.FullName,
				Password: jsonParams.Password,
			}); err == nil {
				c.JSON(http.StatusCreated, gin.H{
					"token": token,
				})
			} else {
				c.Status(http.StatusFound)
			}
		}
	} else {
		c.Status(http.StatusBadRequest)
	}
}

// CheckAvailableEmailAddress ...
func CheckAvailableEmailAddress(c *gin.Context) {
	var jsonParams Email
	if err := c.ShouldBindJSON(&jsonParams); err == nil {
		useCase := domain.NewUserUseCase(data.GetInstance())
		if _, err := mail.ParseAddress(jsonParams.Email); err != nil {
			errorModel := &ErrorModel{
				4002,
				"Invalid Email Address",
			}
			c.JSON(http.StatusBadRequest, errorModel)
		} else if useCase.IsEmailUsed(jsonParams.Email) {
			c.Status(http.StatusFound)
		} else {
			c.Status(http.StatusOK)
		}
	} else {
		c.Status(http.StatusBadRequest)
	}
}

// ActivateAccount ...
func ActivateAccount(c *gin.Context) {
	var jsonParams ActivationToken
	if err := c.ShouldBindJSON(&jsonParams); err == nil {
		useCase := domain.NewUserUseCase(data.GetInstance())
		if err := useCase.ActivateAccount(jsonParams.Token); err != nil {
			c.Status(http.StatusForbidden)
		} else {
			c.Status(http.StatusOK)
		}
	} else {
		c.Status(http.StatusBadRequest)
	}
}

// Login controllers fro Login
func Login(c *gin.Context) {
	var jsonParams EmailLogin
	if err := c.ShouldBindJSON(&jsonParams); err == nil {
		useCase := domain.NewUserUseCase(data.GetInstance())
		if token, err := useCase.DoLogin(jsonParams.Email, jsonParams.Password); err != nil {
			c.Status(http.StatusUnauthorized)
		} else {
			c.JSON(http.StatusOK, token)
		}
	} else {
		c.Status(http.StatusBadRequest)
	}
}

// Profile ...
func Profile(c *gin.Context) {
	user, exist := c.Get("user")
	if !exist {
		c.Status(http.StatusInternalServerError)
	}
	c.JSON(http.StatusOK, user)
	// c.String(http.StatusCreated, "Wow")
}
