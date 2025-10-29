package handler

import (
	"net/http"

	"github.com/Aleksey170999/go-loyaty/internal/models"
	"github.com/Aleksey170999/go-loyaty/internal/service"
	"github.com/gin-gonic/gin"
)

type signUpInput struct {
	UserName  string `json:"username" binding:"required"`
	Password  string `json:"password" binding:"required,min=8"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

func (h *Handler) SignUp(c *gin.Context) {
	var input signUpInput

	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid input body: "+err.Error())
		return
	}

	user := models.User{
		UserName:  input.UserName,
		Password:  input.Password, // Password will be hashed in the service layer
		FirstName: input.FirstName,
		LastName:  input.LastName,
	}

	// Use type assertion to call the Register method
	authService, ok := h.services.Authorization.(interface {
		Register(user models.User) (service.AuthResponse, error)
	})
	if !ok {
		newErrorResponse(c, http.StatusInternalServerError, "internal server error")
		return
	}

	resp, err := authService.Register(user)
	if err != nil {
		if err.Error() == "user with this username already exists" {
			newErrorResponse(c, http.StatusConflict, err.Error())
		} else {
			newErrorResponse(c, http.StatusInternalServerError, err.Error())
		}
		return
	}

	c.JSON(http.StatusOK, resp)
}

type signInInput struct {
	UserName string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (h *Handler) SignIn(c *gin.Context) {
	var input signInInput

	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	// Use type assertion to call the Login method
	authService, ok := h.services.Authorization.(interface {
		Login(username, password string) (service.AuthResponse, error)
	})
	if !ok {
		newErrorResponse(c, http.StatusInternalServerError, "internal server error")
		return
	}

	resp, err := authService.Login(input.UserName, input.Password)
	if err != nil {
		// Return unauthorized for invalid credentials
		newErrorResponse(c, http.StatusUnauthorized, "invalid credentials")
		return
	}

	c.JSON(http.StatusOK, resp)
}
