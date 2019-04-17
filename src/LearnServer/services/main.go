package services

import (
	"LearnServer/services/staffs"
	"LearnServer/services/students"
	"github.com/labstack/echo"
)

// RegisterApis 注册所有路由
func RegisterApis(e *echo.Echo) {
	apiGroup := e.Group("/api/v3")

	students.RegisterStudentsApis(apiGroup)
	staffs.RegisterStaffsApis(apiGroup)
}
