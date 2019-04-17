package validation

import (
	// "LearnServer/utils"
	"LearnServer/utils"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
)

// ValidateUser 验证用户是否正常登陆
func ValidateUser(c echo.Context, id *string) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	idInJwt := claims["id"].(string)
	userType := claims["userType"].(string)
	if userType != "students" || idInJwt == "" {
		return utils.Unauthorized()
	}
	*id = idInJwt
	return nil
}
