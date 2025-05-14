package controller

import (
	"net/http"

	"product-api/internal/dto"
	"product-api/internal/service"

	"github.com/labstack/echo/v4"
)

type ProductController struct {
	service service.ProductService
}

func NewProductController(service service.ProductService) *ProductController {
	return &ProductController{
		service: service,
	}
}

// RegisterRoutes - Route'ları kaydeder
func (c *ProductController) RegisterRoutes(e *echo.Echo) {
	products := e.Group("/api/products")
	products.GET("", c.GetAllProducts)
	products.GET("/:id", c.GetProduct)
	products.POST("", c.CreateProduct)
	products.PUT("/:id", c.UpdateProduct)
	products.PATCH("/:id", c.UpdateProduct)
	products.DELETE("/:id", c.DeleteProduct)
	products.GET("/search", c.SearchProducts)
}

func (c *ProductController) GetAllProducts(ctx echo.Context) error {
	products, err := c.service.GetAllProducts(ctx.Request().Context())
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return ctx.JSON(http.StatusOK, products)
}

func (c *ProductController) GetProduct(ctx echo.Context) error {
	id := ctx.Param("id")
	product, err := c.service.GetProduct(ctx.Request().Context(), id)
	if err != nil {
		return ctx.JSON(http.StatusNotFound, map[string]string{"error": "Product not found"})
	}

	return ctx.JSON(http.StatusOK, product)
}

func (c *ProductController) CreateProduct(ctx echo.Context) error {
	var req dto.CreateProductRequest
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request format"})
	}

	product, err := c.service.CreateProduct(ctx.Request().Context(), &req)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return ctx.JSON(http.StatusCreated, product)
}

func (c *ProductController) UpdateProduct(ctx echo.Context) error {
	id := ctx.Param("id")
	var req dto.UpdateProductRequest
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request format"})
	}

	product, err := c.service.UpdateProduct(ctx.Request().Context(), id, &req)
	if err != nil {
		return ctx.JSON(http.StatusNotFound, map[string]string{"error": "Product not found"})
	}

	return ctx.JSON(http.StatusOK, product)
}

func (c *ProductController) DeleteProduct(ctx echo.Context) error {
	id := ctx.Param("id")
	if err := c.service.DeleteProduct(ctx.Request().Context(), id); err != nil {
		return ctx.JSON(http.StatusNotFound, map[string]string{"error": "Product not found"})
	}

	return ctx.NoContent(http.StatusNoContent)
}

func (c *ProductController) SearchProducts(ctx echo.Context) error {
	var req dto.SearchProductRequest
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request format"})
	}

	products, err := c.service.SearchProducts(ctx.Request().Context(), &req)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return ctx.JSON(http.StatusOK, products)
}
