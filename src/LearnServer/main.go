package main

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"

	"LearnServer/conf"
	"LearnServer/services"
)

func main() {
	rand.Seed(time.Now().UTC().UnixNano())

	e := echo.New()
	e.Debug = conf.AppConfig.Debug

	e.Use(middleware.Logger())
	//e.Use(middleware.Recover())

	services.RegisterApis(e)

	addr := fmt.Sprintf("%s:%d", conf.AppConfig.Host, conf.AppConfig.Port)
	log.Fatal(e.Start(addr))
}
