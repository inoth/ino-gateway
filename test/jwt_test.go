package test

import (
	"fmt"
	"testing"

	jwtauth "github.com/inoth/ino-toybox/utils/jwt_auth"
)

func TestGenJwtStr(t *testing.T) {
	userInfo := map[string]interface{}{
		"name":      "inoth",
		"tenant_id": "10013",
	}
	jwtStr, err := jwtauth.CreateToken(jwtauth.DEFAULT_SIGNKEY, userInfo)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(jwtStr)
	t.Log()
}

func TestParseJwt(t *testing.T) {
	var (
		jwtStr = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJpbm90aCIsImV4cCI6MTY3NzcwNjg2NCwiVXNlckluZm8iOnsibmFtZSI6Imlub3RoIiwidGVuYW50X2lkIjoiMTAwMTMifX0.S8a7Pl1fYnUKtJpc4mXEt7peqoDVlVLMIiDoye_-2Ys"
	)
	customClaims, err := jwtauth.ParseToken(jwtauth.DEFAULT_SIGNKEY, jwtStr)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("%+v\n", customClaims.UserInfo)
	t.Log()
}
