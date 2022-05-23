package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/vincentconace/api-gin/cmd/server/router"
	"github.com/vincentconace/api-gin/internal/domain"
	"github.com/vincentconace/api-gin/pkg/db"
	"github.com/vincentconace/api-gin/pkg/redis"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	// Init database connection
	db := db.InitMysqlDB()
	db.AutoMigrate(&domain.Product{})

	r := gin.Default()

	// Init redis connection
	rd := redis.RedisClient()

	// Run server
	router := router.NewRouter(r, db, rd)
	router.MapaRuter()

	if err := r.Run(":8080"); err != nil {
		panic(err)
	}
}
