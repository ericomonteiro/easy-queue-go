package routes

import (
	"easy-queue-go/src/internal/handlers"

	"github.com/gin-gonic/gin"
)

// SetupRouter configura e retorna o router com todas as rotas definidas
func SetupRouter() *gin.Engine {
	router := gin.Default()

	// Grupo de rotas de health
	router.GET("/health", handlers.HealthCheck)

	return router
}
