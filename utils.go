package utils

import (
	"errors"
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

func init()  {
	_, err := mydb.Connect()

	if err != nil {
		Logger().Error(err.Error())
	}
}

// CreateOrUpdateToken will issue new bearer token
func CreateOrUpdateToken(user dbmodel.User) (map[string]string, *jwt.Token, error) {
	var f, err = os.ReadFile(os.Getenv("CERTIFICATE_FILE"))
	expireAt := time.Now().Add(25 * time.Minute)
	refreshExpireAt := time.Now().Add(50 * time.Minute)

	if err != nil {
		return map[string]string{}, nil, err
	}

	claim := &dbmodel.UserClaims{
		User: user,
		StandardClaims: &jwt.StandardClaims{
			ExpiresAt: expireAt.Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	signedString, err := token.SignedString(f)

	if err != nil {
		return map[string]string{}, nil, err
	}

	refreshClaim := &dbmodel.UserClaims{
		User: user,
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
		"auth_token":   signedString,
		"refresh_token": refreshSignedString,
	}, token, nil
}

// ValidateAuth ValidateToken for auth header request
func ValidateAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		var jwtToken dbmodel.JWT
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
			c.AbortWithStatusJSON(http.StatusBadRequest, dbmodel.Error{
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

func checkIfRefreshTokenNotExpired() (*dbmodel.JWT, error) {
	var jwToken dbmodel.JWT
	row := mydb.DB.QueryRow("SELECT refresh_token, device_id from oauth where device_id=?", cookie)

	err := row.Scan(&jwToken.Authorization, &cookie)

	if err != nil {
		Logger().Error(err.Error())
		return nil, err
	}

	return &jwToken, nil
}

func getToken(jwToken dbmodel.JWT) (*dbmodel.UserClaims, error) {
	if cookie == "" {
		return nil, errors.New("cookie expired or not exist")
	}

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
			return r.Claims.(*dbmodel.UserClaims), nil
		}

		return &dbmodel.UserClaims{}, err
	}

	if token.Valid {
		return token.Claims.(*dbmodel.UserClaims), nil
	}

	return &dbmodel.UserClaims{}, nil
}

func verify(jwToken dbmodel.JWT, f []byte) (*jwt.Token, error) {
	token, err := jwt.ParseWithClaims(jwToken.Authorization, &dbmodel.UserClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return f, nil
	})
	return token, err
}