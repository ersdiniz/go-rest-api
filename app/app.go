package app

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"go-rest-api/router"
	"runtime"
)

func GetEngine() *gin.Engine {
	route := router.NewRoute()

	gin.SetMode(gin.ReleaseMode)

	engine := gin.Default()

	initCors(engine)
	router := engine.Group("/api")

	route.BuildRoutes(router)

	runtime.SetBlockProfileRate(1)

	return engine
}

func initCors(engine *gin.Engine) {
	config := cors.DefaultConfig()
	config.AllowCredentials = true
	config.AllowOriginFunc = func(string) bool {
		return true
	}
	engine.Use(cors.New(config))
}
