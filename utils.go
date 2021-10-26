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
	"time"
)

var (
	cookie string
)

func init() {
	_, err := mydb.Connect()

	if err != nil {
		Logger().Error(err.Error())
	}
}

// CreateOrUpdateToken will issue new bearer token
func CreateOrUpdateToken(user *models.User) (map[string]string, *jwt.Token, error) {
	var f, err = os.ReadFile(os.Getenv("CERTIFICATE_FILE"))
	expireAt := time.Now().Add(2 * time.Minute)
	refreshExpireAt := time.Now().Add(5 * time.Minute)

	if err != nil {
		return map[string]string{}, nil, err
	}

	claim := &models.UserClaims{
		User:    *user,
		Refresh: false,
		StandardClaims: &jwt.StandardClaims{
			ExpiresAt: expireAt.Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	signedString, err := token.SignedString(f)

	if err != nil {
		return map[string]string{}, nil, err
	}

	refreshClaim := &models.UserClaims{
		User:    *user,
		Refresh: true,
		StandardClaims: &jwt.StandardClaims{
			ExpiresAt: refreshExpireAt.Unix(),
		},
	}
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaim)
	refreshSignedString, err := refreshToken.SignedString(f)

	if err != nil {
		return map[string]string{}, nil, err
	}

	return map[string]string{
		"auth_token":    signedString,
		"refresh_token": refreshSignedString,
	}, token, nil
}

func IssueTokenPair(c *gin.Context) (*models.JWT, error) {
	var jwToken models.JWT
	if cookie == "" {
		return nil, errors.New("device id not found")
	}

	if c.GetHeader("Authorization") == "" {
		return nil, errors.New("authorization header missing")
	}

	f, _ := os.ReadFile(os.Getenv("CERTIFICATE_FILE"))
	row := mydb.DB.QueryRow("SELECT refresh_token, device_id from oauth where device_id = ? and auth_token = ?", cookie, c.GetHeader("Authorization"))

	err := row.Scan(&jwToken.Authorization, &cookie)

	if err != nil {
		Logger().Error(err.Error())
		return nil, err
	}

	token, err := verify(jwToken, f)

	if err != nil {
		return nil, nil
	}

	if token.Valid {
		return &jwToken, nil
	}
	return nil, err
}

// ValidateAuth ValidateToken for auth header request
func ValidateAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		var jwtToken models.JWT
		var token string
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

		row := mydb.DB.QueryRow("select auth_token from oauth where device_id = ? and auth_token = ?", cookie, c.GetHeader("Authorization"))

		err = row.Scan(&token)

		fmt.Println(err)

		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, models.Error{
				Message: "Unauthorized, No Token",
				Code:    401,
			})
			return
		}

		user, err := verifyAuthToken(jwtToken)

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

// verifyAuthToken will verify the passed token
func verifyAuthToken(jwToken models.JWT) (*models.UserClaims, error) {
	if cookie == "" {
		return nil, errors.New("device expired or not exist")
	}

	f, _ := os.ReadFile(os.Getenv("CERTIFICATE_FILE"))

	token, err := verify(jwToken, f)
	if err != nil {
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
