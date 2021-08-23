package utils

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/thisismyaim/utils/models"
	"net/http"
	"os"
)

// ValidateAuth ValidateToken for auth header request
func ValidateAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		var jwtToken models.JWT
		err := c.BindHeader(&jwtToken)

		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, models.Error{
				Message: "missing authorization in header",
				Code:    400,
			})

			return
		}

		f, _ := os.ReadFile(os.Getenv("CERTIFICATE_FILE"))

		token, err := jwt.Parse(jwtToken.Authorization, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("unexpected signing method")
			}
			return f, nil
		})
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, models.Error{
				Message: err.Error(),
				Code:    401,
			})
		}

		if token.Valid {
			fmt.Println(token.Claims)
			c.Next()
		}
	}
}
