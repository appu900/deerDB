package auth

import (
	"fmt"
	"net/http"
	"github.com/labstack/echo/v4"
)

type AuthHttpHandler struct {
	authService *AuthService
}

func NewAuthHttpHandler(authService *AuthService) *AuthHttpHandler {
	return &AuthHttpHandler{
		authService: authService,
	}
}

func (h *AuthHttpHandler) HandleUserRegistration(c echo.Context) error {
	var req RegisterRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error:": err.Error()})
	}
	user, err := h.authService.RegisterUser(req.UserName, req.Email, req.Password)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": err.Error(),
		})
	}
	resp := RegisterResponse{
		UserID:    user.ID,
		UserName:  user.Username,
		UserEmail: user.Email,
	}
	return c.JSON(http.StatusCreated, resp)
}

func (h *AuthHttpHandler) HandleUserLogin(c echo.Context) error {
	var req LoginRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error:": err.Error()})
	}

	user, err := h.authService.Login(req.Email, req.Password)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": err.Error(),
		})
	}
	token, err := h.authService.GenerateToken(user)
	if err != nil {
		fmt.Println("error occured in generating token")
		return c.JSON(http.StatusInternalServerError,map[string]string{
			"error":"Internal server error",
		})
	}
	return c.JSON(http.StatusOK,map[string]interface{}{
		"user":ToUserResponse(user),
		"token":token,
	})
}
