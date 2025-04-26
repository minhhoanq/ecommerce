package routes

import (
	"bytes"
	"fmt"
	"io"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/minhhoanq/ecommerce/common/logger"
	"github.com/minhhoanq/ecommerce/gateway/internal/generated/catalog_service"
	"github.com/minhhoanq/ecommerce/gateway/internal/modules/catalog"
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
	handler.GET("", c.getListProduct)
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

func (catalogHandler *catalogHandlerFunc) createProduct(c echo.Context) error {
	catalogHandler.l.Info("routes.catalog.createProduct create product")
	req := new(CreateProductParams)
	if err := parseJSONFormField(c, "dawd", req); err != nil {
		return Error(c, err, http.StatusBadRequest)
	}

	catalogHandler.l.Info("routes.catalog.createProduct create product 1")
	image, err := c.FormFile("image")
	if err != nil {
		return Error(c, err, http.StatusBadRequest)
	}
	catalogHandler.l.Info("routes.catalog.createProduct create product 2")

	file, err := image.Open()
	if err != nil {
		return Error(c, err, http.StatusBadRequest)
	}
	defer file.Close()

	buf := new(bytes.Buffer)
	if _, err := io.Copy(buf, file); err != nil {
		return Error(c, fmt.Errorf("failed to copy buffer file image err: %w", err), http.StatusBadRequest)
	}

	metadata := &catalog_service.CreateProductRequest{
		Name:        req.Name,
		Description: req.Description,
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

		metadata.Skus = append(metadata.Skus, sku)
	}

	imageInfo := &catalog_service.ImageInfo{
		OriginalFileName: image.Filename,
		ImageData:        buf.Bytes(),
	}

	arg := &catalog_service.CreateProductWithImageRequest{
		Metadata:  metadata,
		ImageInfo: imageInfo,
	}

	product, err := catalogHandler.catalogManagementOperator.CreateProduct(c.Request().Context(), arg)
	if err != nil {
		return Error(c, err, http.StatusBadRequest)
	}

	return Success(c, product, http.StatusCreated, nil)
}

type ListProductRequest struct {
	Page     int32
	PageSize int32
}

func (catalogHandler *catalogHandlerFunc) getListProduct(c echo.Context) error {
	req := new(ListProductRequest)
	req.Page = parseQueryIntParam(c, "page", 1, 1, 10)
	req.PageSize = parseQueryIntParam(c, "page_size", 5, 5, 10)
	arg := &catalog_service.ListProductRequest{
		Page:     req.Page,
		PageSize: req.PageSize,
	}

	products, err := catalogHandler.catalogManagementOperator.GetListProduct(c.Request().Context(), arg)
	if err != nil {
		return Error(c, err, http.StatusBadRequest)
	}

	return Success(c, products, http.StatusOK, nil)
}
