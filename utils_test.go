package utils

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/thisismyaim/utils/utils@v1.0/models"
	"os"
	"testing"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func TestValidateAuth(t *testing.T) {
	jwToken := models.JWT{Authorization: os.Getenv("TEST_ACCESS_TOKEN")}

	u, err := verifyAuthToken(jwToken)

	if err != nil {
		t.Error(err)
	}

	fmt.Println(u.User.Email)
}
