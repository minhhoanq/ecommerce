package routes

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
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

func (catalogHandler *catalogHandlerFunc) createProduct(c echo.Context) error {

	metadataJson := c.FormValue("metadata")
	var req CreateProductParams
	if err := json.Unmarshal([]byte(metadataJson), &req); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	image, err := c.FormFile("image")
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	file, err := image.Open()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}
	defer file.Close()

	buf := new(bytes.Buffer)
	if _, err := io.Copy(buf, file); err != nil {
		return fmt.Errorf("failed to copy buffer file image")
	}

	catalogHandler.l.Info("order params: ", zap.String("order: ", req.Name))
	catalogHandler.l.Info("order params: ", zap.String("iamge: ", image.Filename))

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
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusCreated, product)
}
