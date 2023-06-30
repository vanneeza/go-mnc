package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

type Auth struct {
	Code    int
	Status  string
	Message string
}

func AuthMiddleware(secretKey string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		tokenString, err := ctx.Cookie("Authorization")
		if err != nil || tokenString == "" {
			result := Auth{
				Code:    http.StatusUnauthorized,
				Status:  "Unauthorized",
				Message: "You are not logged in. Please log in or register first",
			}
			ctx.JSON(http.StatusUnauthorized, result)
			ctx.Abort()
			return
		}

		jwtKeyByte := []byte(secretKey)
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return jwtKeyByte, nil
		})

		if err != nil || !token.Valid {
			result := Auth{
				Code:    http.StatusUnauthorized,
				Status:  "Unauthorized",
				Message: "You are not logged in. Please log in or register first",
			}
			ctx.JSON(http.StatusUnauthorized, result)
			ctx.Abort()
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			result := Auth{
				Code:    http.StatusUnauthorized,
				Status:  "Unauthorized",
				Message: "Invalid token claims",
			}
			ctx.JSON(http.StatusUnauthorized, result)
			ctx.Abort()
			return
		}

		ctx.Set("claims", claims)
		ctx.Next()
	}
}
