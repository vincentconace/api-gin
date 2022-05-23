package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/vincentconace/api-gin/internal/domain"
	"github.com/vincentconace/api-gin/internal/product"
	"github.com/vincentconace/api-gin/pkg/web"
)

var (
	ErrInternal = errors.New("internal error")
)

type ProductHandler struct {
	productService product.Service
	redis          *redis.Client
}

func NewProductHandler(productService product.Service, rd *redis.Client) *ProductHandler {
	return &ProductHandler{productService: productService, redis: rd}
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
		key := fmt.Sprintf("product[%d]", idConv)
		productRedis, err := h.redis.Get(c, key).Result()
		if err != nil {
			fmt.Println(err)
		}
		if productRedis != "" {
			var product domain.Product
			err = json.Unmarshal([]byte(productRedis), &product)
			if err != nil {
				fmt.Println(err)
			}
			fmt.Println("Devolvio desde redis")
			web.Success(c, http.StatusOK, product)
			return
		}
		product, err := h.productService.GetById(c, uint(idConv))
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

		// productReflect := reflect.ValueOf(p)
		// var valuesNil []string
		// for i := 0; i < productReflect.NumField(); i++ {
		// 	if e := productReflect.Field(i); e.IsNil() &&
		// 		productReflect.Type().Field(i).Name != "ID" {
		// 		valuesNil = append(valuesNil, productReflect.Type().Field(i).Name)
		// 	}
		// }

		// if len(valuesNil) > 0 {
		// 	web.Error(c, 422, "required fields: %s", strings.Join(valuesNil, ", "))
		// 	return
		// }

		product, err := h.productService.Create(c, p)
		if err != nil {
			web.Error(c, http.StatusBadRequest, err.Error())
			return
		}

		dataByte, err := json.Marshal(product)
		if err != nil {
			fmt.Println(err)
		}

		key := fmt.Sprintf("product[%d]", product.Model.ID)
		productRedis, _ := h.redis.Set(c, key, string(dataByte), 24*time.Hour).Result()
		if productRedis != "" {
			fmt.Println("El producto se guardo correctamente", productRedis)
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

		productUpdated, err := h.productService.Update(c, uint(idConv), p)
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

		err = h.productService.Delete(c, uint(idConv))
		if err != nil {
			web.Error(c, http.StatusNotFound, "failed to delete product")
			return
		}

		key := fmt.Sprintf("product[%d]", idConv)
		result, _ := h.redis.Del(c, key).Result()
		if result > 0 {
			fmt.Println("El producto se elimino correctamente")
		}

		web.Success(c, http.StatusNoContent, "")
	}
}
