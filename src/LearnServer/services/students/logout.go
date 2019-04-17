package students

import (
	"net/http"

	"github.com/labstack/echo"
)

func logoutHandler(c echo.Context) error {
	c.SetCookie(&http.Cookie{
		Name:     "TOKEN",
		Path:     "/api/v3/students/",
		Value:    "",
		MaxAge:   -1,
		HttpOnly: true,
	})
	return c.NoContent(http.StatusOK)
}
