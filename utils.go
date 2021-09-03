package utils

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/thisismyaim/utils/models"
	"github.com/thisismyaim/utils/mydb"
	"net/http"
	"os"
)

var (
	cookie string
)

func init()  {
	_, err := mydb.Connect()

	if err != nil {
		Logger().Error(err.Error())
	}
}

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

func checkIfRefreshTokenNotExpired() (*models.JWT, error) {
	var jwToken models.JWT
	row := mydb.DB.QueryRow("SELECT refresh_token, device_id from oauth where refresh_token=?", cookie)

	err := row.Scan(&jwToken.Authorization, &cookie)

	if err != nil {
		Logger().Error(err.Error())
		return nil, err
	}

	return &jwToken, nil
}

func getToken(jwToken models.JWT) (*models.UserClaims, error) {
	f, _ := os.ReadFile(os.Getenv("CERTIFICATE_FILE"))

	token, err := verify(jwToken, f)
	if err != nil {
		refreshToken, err := checkIfRefreshTokenNotExpired()
		if err != nil {
			return nil, err
		}

		r, refreshTokenErr := verify(*refreshToken, f)

		if refreshTokenErr != nil {
			return nil, refreshTokenErr
		}

		if r.Valid {
			return r.Claims.(*models.UserClaims), nil
		}

		return &models.UserClaims{}, err
	}

	if token.Valid {
		return token.Claims.(*models.UserClaims), nil
	}

	return &models.UserClaims{}, nil
}

func verify(jwToken models.JWT, f []byte) (*jwt.Token, error) {
	token, err := jwt.ParseWithClaims(jwToken.Authorization, &models.UserClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return f, nil
	})
	return token, err
}