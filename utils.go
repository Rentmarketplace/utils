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

var (
	cookie string
)

// ValidateAuth ValidateToken for auth header request
func ValidateAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		var jwtToken models.JWT
		deviceCookie, err := c.Cookie("device")

		if err != nil {
			c.AbortWithStatusJSON(404, gin.H{
				"error": err.Error(),
			})
			return
		}

		cookie = deviceCookie

		err = c.BindHeader(&jwtToken)

		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, models.Error{
				Message: "missing authorization in header",
				Code:    400,
			})
			return
		}

		user, err := getToken(jwtToken)

		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.Set("user", user.User)
		c.Next()
	}
}

func RegenerateToken()  {
	fmt.Println("Test")
}

func checkIfRefreshTokenNotExpired()  {
	fmt.Println(cookie)
}

func getToken(jwToken models.JWT) (*models.UserClaims, error) {
	f, _ := os.ReadFile(os.Getenv("CERTIFICATE_FILE"))

	token, err := jwt.ParseWithClaims(jwToken.Authorization, &models.UserClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return f, nil
	})
	if err != nil {
		checkIfRefreshTokenNotExpired()
		return &models.UserClaims{}, err
	}

	if token.Valid {
		return token.Claims.(*models.UserClaims), nil
	}

	return &models.UserClaims{}, nil
}