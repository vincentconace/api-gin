package router

import (
	"github.com/gin-gonic/gin"
	"github.com/vincentconace/api-gin/cmd/server/handler"
	"github.com/vincentconace/api-gin/internal/product"
	"github.com/vincentconace/api-gin/pkg/redis"
	"gorm.io/gorm"
)

type Router interface {
	MapaRuter()
}

type router struct {
	r  *gin.Engine
	rg *gin.RouterGroup
	db *gorm.DB
}

func NewRouter(r *gin.Engine, db *gorm.DB) Router {
	return &router{r: r, db: db}
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
	redisClient := redis.NewRedisClient()
	handler := handler.NewProductHandler(service, redisClient)

	// Product routes
	r.rg.POST("/products", handler.Create())
	r.rg.GET("/products", handler.Get())
	r.rg.GET("/products/:id", handler.GetById())
	r.rg.PUT("/products/:id", handler.Update())
	r.rg.DELETE("/products/:id", handler.Delete())
}
