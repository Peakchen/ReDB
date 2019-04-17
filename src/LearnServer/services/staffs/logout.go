package staffs

import (
	"net/http"

	"github.com/labstack/echo"
)

func logoutHandler(c echo.Context) error {
	c.SetCookie(&http.Cookie{
		Name:     "TOKEN",
		Path:     "/api/v3/staffs/",
		Value:    "",
		MaxAge:   -1,
		HttpOnly: true,
	})
	return c.JSON(http.StatusOK, "Successfully logged out!")
}
