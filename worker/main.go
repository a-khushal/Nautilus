package main

import (
	"context"
	"fmt"
	"log"

	"github.com/a-khushal/Nautilus/worker/config"
	"github.com/a-khushal/Nautilus/worker/db"
	"github.com/a-khushal/Nautilus/worker/processor"
	"github.com/a-khushal/Nautilus/worker/redis"
)

var (
	ctx  = context.Background()
	JOBS = "JOBS"
)

func main() {
	fmt.Println("worker started")
	cfg := config.Get()

	gormDB := db.InitDB(cfg.DBConn)

	redisClient := redis.InitRedis(cfg.RedisAddr)

	for {
		result, err := redisClient.BRPop(ctx, 0, JOBS).Result()
		if err != nil {
			log.Printf("error popping job: %v", err)
			continue
		}

		jobID := result[1]
		processor.Process(ctx, jobID, gormDB, redisClient)
	}
}
