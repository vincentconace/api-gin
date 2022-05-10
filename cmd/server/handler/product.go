package handler

import (
	"errors"
	"net/http"
	"reflect"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/vincentconace/api-gin/internal/domain"
	"github.com/vincentconace/api-gin/internal/product"
	"github.com/vincentconace/api-gin/pkg/web"
)

var (
	ErrInternal = errors.New("internal error")
)

type ProductHandler struct {
	productService product.Service
}

func NewProductHandler(productService product.Service) *ProductHandler {
	return &ProductHandler{productService: productService}
}

func (h *ProductHandler) Get() gin.HandlerFunc {
	return func(c *gin.Context) {
		products, err := h.productService.Get(c)
		if err != nil {
			web.Error(c, http.StatusNotFound, "product not found")
			return
		}
		web.Success(c, http.StatusOK, products)
	}
}

func (h *ProductHandler) GetById() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		idConv, err := strconv.Atoi(id)
		if err != nil {
			web.Error(c, http.StatusBadRequest, "invalid id")
		}
		product, err := h.productService.GetById(c, idConv)
		if err != nil {
			web.Error(c, http.StatusNotFound, "product not found")
			return
		}
		web.Success(c, http.StatusOK, product)
	}
}

func (h *ProductHandler) Create() gin.HandlerFunc {
	return func(c *gin.Context) {
		var p domain.Product
		if err := c.ShouldBindJSON(&p); err != nil {
			web.Error(c, http.StatusUnprocessableEntity, "invalid request")
			return
		}

		productReflect := reflect.ValueOf(p)
		var valuesNil []string
		for i := 0; i < productReflect.NumField(); i++ {
			if e := productReflect.Field(i); e.IsNil() &&
				productReflect.Type().Field(i).Name != "ID" {
				valuesNil = append(valuesNil, productReflect.Type().Field(i).Name)
			}
		}

		if len(valuesNil) > 0 {
			web.Error(c, 422, "required fields: %s", strings.Join(valuesNil, ", "))
			return
		}

		product, err := h.productService.Create(c, p)
		if err != nil {
			web.Error(c, http.StatusBadRequest, err.Error())
			return
		}
		web.Success(c, http.StatusOK, product)
	}
}

func (h *ProductHandler) Update() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		idConv, err := strconv.Atoi(id)
		if err != nil {
			web.Error(c, http.StatusBadRequest, "invalid id")
			return
		}
		var p domain.Product
		if err := c.ShouldBindJSON(&p); err != nil {
			web.Error(c, http.StatusUnprocessableEntity, "unprocesable request")
			return
		}

		productUpdated, err := h.productService.Update(c, idConv, p)
		if err != nil {
			web.Error(c, http.StatusNotFound, err.Error())
			return
		}

		web.Success(c, http.StatusOK, productUpdated)
	}
}

func (h *ProductHandler) Delete() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		idConv, err := strconv.Atoi(id)
		if err != nil {
			web.Error(c, http.StatusBadRequest, "invalid id")
			return
		}

		err = h.productService.Delete(c, idConv)
		if err != nil {
			web.Error(c, http.StatusNotFound, "failed to delete product")
			return
		}
		web.Success(c, http.StatusNoContent, "")
	}
}
