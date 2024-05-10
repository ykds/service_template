package api

import (
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"service_template/docs"
	"service_template/internal/middleware"
	"service_template/internal/service"
)

func InitRouter(engine *gin.Engine, api *Api, srv *service.Service) {
	auth := func(token string) (int, error) {
		return 0, nil
	}
	_ = engine.Group("/api", middleware.Authentication(auth))
}

func RegisterSwagger(engine *gin.Engine, host string) {
	docs.SwaggerInfo.BasePath = "/api"
	docs.SwaggerInfo.Title = ""
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Schemes = []string{"http", "https"}
	docs.SwaggerInfo.Host = host
	engine.GET("/inner/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
}
