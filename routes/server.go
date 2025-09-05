package routes

import (
	"backend/config"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

func InitServer(db *pgxpool.Pool) {
	globalEnv := config.GetEnv()

	config.GetLogger().Info("v0.1.0")

	config.GetLogger().Info("Initializing API routes")
	router := SetupRoutes(db)

	config.GetLogger().Info("Starting up Gin server")

	err := router.Run(fmt.Sprintf(":%s", globalEnv.GinPort))
	if err != nil {
		config.GetLogger().Fatal("failed to run HTTP server: %s", zap.Error(err))
	}
}
