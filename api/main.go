package main

import (
	"github.com/a-khushal/Nautilus/api/config"
	"github.com/a-khushal/Nautilus/api/db"
	"github.com/a-khushal/Nautilus/api/redis"
	"github.com/a-khushal/Nautilus/api/routes"
	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.Get()

	gormDB := db.InitDB(cfg.DBConn)

	redisClient := redis.InitRedis(cfg.RedisAddr)

	r := gin.Default()
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	routes.RegisterJobRoutes(r, gormDB, redisClient)

	r.Run(":3000")
}
