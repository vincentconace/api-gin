package main

import (
	"github.com/gin-gonic/gin"
	"github.com/vincentconace/api-gin/cmd/server/router"
	"github.com/vincentconace/api-gin/pkg/db"
)

func main() {
	// Init database connection
	db, err := db.Init()
	if err != nil {
		panic(err)
	}
	defer db.Close()
	r := gin.Default()

	// Run server
	router := router.NewRouter(r, db)
	router.MapaRuter()

	if err := r.Run(":8080"); err != nil {
		panic(err)
	}
}
