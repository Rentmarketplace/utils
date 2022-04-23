package utils

import (
	"fmt"
	"github.com/Rentmarketplace/utils/models"
	"github.com/gin-gonic/gin"
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
