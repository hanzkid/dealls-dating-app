package middleware

// parse the token from the request

import (
	"main/helpers"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

func AuthMiddleware(JWTSecret string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			token := c.Request().Header.Get("Authorization")
			if token == "" {
				helpers.ResponseWithError(c, http.StatusUnauthorized, "Unauthorized")
				return nil
			}

			parts := strings.Split(token, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				helpers.ResponseWithError(c, http.StatusUnauthorized, "Unauthorized")
				return nil
			}
			tokenValidated, err := helpers.ValidateToken(parts[1], JWTSecret)
			if err != nil {
				helpers.ResponseWithError(c, http.StatusUnauthorized, "Unauthorized")
				return nil
			}
			if !tokenValidated.Valid {
				helpers.ResponseWithError(c, http.StatusUnauthorized, "Unauthorized")
				return nil
			}

			claims := tokenValidated.Claims.(jwt.MapClaims)
			userID, ok := claims["user_id"].(float64)
			if !ok {
				helpers.ResponseWithError(c, http.StatusUnauthorized, "Unauthorized")
				return nil
			}
			profileId, ok := claims["profile_id"].(float64)
			if !ok {
				helpers.ResponseWithError(c, http.StatusUnauthorized, "Unauthorized")
				return nil
			}
			c.Set("user_id", int(userID))
			c.Set("profile_id", int(profileId))
			c.Set("name", claims["name"])
			c.Set("email", claims["email"])
			return next(c)
		}
	}

}
