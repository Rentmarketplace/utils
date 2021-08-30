package utils

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/thisismyaim/utils/models"
	"testing"
)

func init()  {
	gin.SetMode(gin.TestMode)
}

func TestValidateAuth(t *testing.T) {
	jwToken := models.JWT{Authorization: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJVc2VyIjp7ImlkIjo1ODYsImVtYWlsIjoidGhpc2lzbXlhaW1AZ21haWwuY29tIiwiZmlyc3RuYW1lIjoiIiwibGFzdG5hbWUiOiIiLCJQaG9uZSI6IiJ9LCJleHAiOjE2MzAzMDAyNjF9.BF7ysY5YegHW_ilR3R49n2oLclJRAmH2v8RSvw8rW4k"}

	u, err := getToken(jwToken)

	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(u.User.Email)
}

