package auth

import (
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"

)

func RequireAuth() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			tokenHeader := c.Request().Header.Get("Authorization")
			if tokenHeader == "" {
				return c.JSON(http.StatusUnauthorized, echo.Map{"error": "auth header is empty"})
			}
			parts := strings.Split(tokenHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				return c.JSON(http.StatusUnauthorized, echo.Map{"error": "Inavalid token format"})
			}

			tokenStr := parts[1]
			claims := &Claims{}
			token, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (interface{}, error) {
				return []byte("my-secret"), nil
			})
			if err != nil || !token.Valid {
				return c.JSON(http.StatusUnauthorized, echo.Map{
					"error": "Invalid or expired token",
				})
			}
			c.Set("user_id", claims.UserID)
			return next(c)
		}
	}
}
