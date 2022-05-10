package router

import (
	"database/sql"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/vincentconace/api-gin/cmd/server/handler"
	"github.com/vincentconace/api-gin/internal/product"
)

type Router interface {
	MapaRuter()
}

type router struct {
	r  *gin.Engine
	rg *gin.RouterGroup
	db *sql.DB
	rd *redis.Client
}

func NewRouter(r *gin.Engine, db *sql.DB, rd *redis.Client) Router {
	return &router{r: r, db: db, rd: rd}
}

func (r *router) MapaRuter() {
	r.setGroup()

	r.buildProductRoutes()
}

func (r *router) setGroup() {
	// General routes
	r.rg = r.r.Group("/api/v1")
}

func (r *router) buildProductRoutes() {
	// Repository, service and handler
	repository := product.NewRepository(r.db)
	service := product.NewService(repository)
	handler := handler.NewProductHandler(service, r.rd)

	// Product routes
	r.rg.POST("/products", handler.Create())
	r.rg.GET("/products", handler.Get())
	r.rg.GET("/products/:id", handler.GetById())
	r.rg.PATCH("/products/:id", handler.Update())
	r.rg.DELETE("/products/:id", handler.Delete())
}
