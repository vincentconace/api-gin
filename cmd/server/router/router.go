package router

import (
	"database/sql"

	"github.com/gin-gonic/gin"
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
}

func NewRouter(r *gin.Engine, db *sql.DB) Router {
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
	handler := handler.NewProductHandler(service)

	// Product routes
	r.rg.POST("/products", handler.Create())
	r.rg.GET("/products", handler.Get())
	r.rg.GET("/products/:id", handler.GetById())
	r.rg.PATCH("/products/:id", handler.Update())
	r.rg.DELETE("/products/:id", handler.Delete())
}
