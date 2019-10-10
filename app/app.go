package app

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"go-rest-api/docs"
	"go-rest-api/env"
	"go-rest-api/router"
	"runtime"
	"github.com/swaggo/gin-swagger"
	swaggerFiles "github.com/swaggo/files"
)

func GetEngine() *gin.Engine {
	route := router.NewRoute()

	gin.SetMode(gin.ReleaseMode)

	engine := gin.Default()

	initCors(engine)
	router := engine.Group("/api")

	route.BuildRoutes(router)

	runtime.SetBlockProfileRate(1)

	initSwagger(engine)

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

func initSwagger(engine *gin.Engine) {
	docs.SwaggerInfo.Title = "GoLang REST API"
	docs.SwaggerInfo.Description = "REST Api com GoLang."
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Host = "localhost:" + env.SystemPort()
	docs.SwaggerInfo.BasePath = "/api/v1"
	docs.SwaggerInfo.Schemes = []string{"http"}

	url := ginSwagger.URL("http://localhost:" + env.SystemPort() + "/swagger/doc.json")
	engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))
}