package utils

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/thisismyaim/utils/dbmodel"
	"os"
	"testing"
)

func init()  {
	gin.SetMode(gin.TestMode)
}

func TestValidateAuth(t *testing.T) {
	jwToken := dbmodel.JWT{Authorization: os.Getenv("TEST_ACCESS_TOKEN")}

	u, err := getToken(jwToken)

	if err != nil {
		t.Error(err)
	}

	fmt.Println(u.User.Email)
}
