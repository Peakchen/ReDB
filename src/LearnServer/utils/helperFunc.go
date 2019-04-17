package utils

import (
	"net/http"

	// "LearnServer/conf"
	"LearnServer/conf"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

// JwtMiddleware 创建jwt中间件
func JwtMiddleware() echo.MiddlewareFunc {
	cfg := middleware.DefaultJWTConfig
	cfg.TokenLookup = "cookie:TOKEN"
	cfg.SigningKey = []byte(conf.AppConfig.Secret)
	jwt := middleware.JWTWithConfig(cfg)
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		handler := jwt(next)
		return func(c echo.Context) error {
			err := handler(c)
			if err != nil {
				e, ok := err.(*echo.HTTPError)
				if !ok {
					return err
				}
				if e.Code == http.StatusBadRequest {
					return echo.NewHTTPError(http.StatusUnauthorized, "Missing jwt")
				}
			}
			return err
		}
	}
}

// NotFound 返回404
func NotFound(msg ...interface{}) error {
	return echo.NewHTTPError(http.StatusNotFound, msg...)
}

// InvalidParams 返回422
func InvalidParams(msg ...interface{}) error {
	return echo.NewHTTPError(http.StatusUnprocessableEntity, msg...)
}

// Unauthorized 返回401
func Unauthorized(msg ...interface{}) error {
	return echo.NewHTTPError(http.StatusUnauthorized, msg...)
}

// Forbidden 返回403
func Forbidden(msg ...interface{}) error {
	return echo.NewHTTPError(http.StatusForbidden, msg...)
}
