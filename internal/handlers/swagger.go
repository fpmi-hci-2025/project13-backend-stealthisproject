package handlers

import (
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/gin-gonic/gin"
)

// SwaggerHandler returns swagger handler
func SwaggerHandler() gin.HandlerFunc {
	return ginSwagger.WrapHandler(swaggerFiles.Handler)
}

