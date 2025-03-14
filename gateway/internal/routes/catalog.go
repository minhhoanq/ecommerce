package routes

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/minhhoanq/ecommerce/common/logger"
	"github.com/minhhoanq/ecommerce/gateway/internal/generated/catalog_service"
	"github.com/minhhoanq/ecommerce/gateway/internal/modules/catalog"
	"go.uber.org/zap"
)

type catalogHandlerFunc struct {
	l                         logger.Interface
	catalogManagementOperator catalog.CatalogManagementOperator
}

func newCatalogRouter(handler *echo.Group, l logger.Interface, catalogManagemntOperator catalog.CatalogManagementOperator) {
	c := &catalogHandlerFunc{
		l:                         l,
		catalogManagementOperator: catalogManagemntOperator,
	}

	handler.POST("", c.createProduct)
}

type CreateProductParams struct {
	Name        string            `json:"name"`
	Description string            `json:"description"`
	Image       string            `json:"image"`
	CategoryID  int32             `json:"category_id"`
	BrandID     int32             `json:"brand_id"`
	SKUs        []CreateSKUParams `json:"skus"`
}

type CreateSKUParams struct {
	Name       string                  `json:"name"`
	Slug       string                  `json:"slug"`
	Price      int32                   `json:"original_price"`
	Inventory  int32                   `json:"initial_stock"`
	Attributes []CreateAttributeParams `json:"attributes"`
}

type CreateAttributeParams struct {
	AttributeID int32  `json:"attribute_id"`
	Value       string `json:"value"`
}

func (c *catalogHandlerFunc) createProduct(e echo.Context) error {
	var req CreateProductParams

	if err := e.Bind(&req); err != nil {
		return e.JSON(http.StatusBadRequest, err.Error())
	}

	c.l.Info("order params: ", zap.String("order: ", req.Name))
	c.l.Info("order params: ", zap.Int32("order: ", req.SKUs[0].Attributes[0].AttributeID))

	arg := &catalog_service.CreateProductRequest{
		Name:        req.Name,
		Description: req.Description,
		Image:       "",
		CategoryId:  req.CategoryID,
		BrandId:     req.BrandID,
		Skus:        make([]*catalog_service.SKUToCreate, 0, len(req.SKUs)),
	}

	for _, skuDetail := range req.SKUs {
		sku := &catalog_service.SKUToCreate{
			Name:          skuDetail.Name,
			Slug:          skuDetail.Slug,
			OriginalPrice: skuDetail.Price,
			InitialStock:  skuDetail.Inventory,
			Attributes:    make([]*catalog_service.AttributeValue, 0, len(skuDetail.Attributes)),
		}

		for _, attr := range skuDetail.Attributes {
			sku.Attributes = append(sku.Attributes, &catalog_service.AttributeValue{
				AttributeId: attr.AttributeID,
				Value:       attr.Value,
			})
		}

		fmt.Println("sku", len(sku.Attributes))

		arg.Skus = append(arg.Skus, sku)
	}

	fmt.Println("dawjdaiw", arg.Skus[0].Name)

	product, err := c.catalogManagementOperator.CreateProduct(e.Request().Context(), arg)
	fmt.Println("dawjdaiw 2")

	if err != nil {
		return e.JSON(http.StatusBadRequest, err.Error())
	}

	return e.JSON(http.StatusCreated, product)
}
