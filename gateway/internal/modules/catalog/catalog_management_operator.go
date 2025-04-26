package catalog

import (
	"context"

	"github.com/minhhoanq/ecommerce/common/logger"
	"github.com/minhhoanq/ecommerce/gateway/internal/generated/catalog_service"
	"go.uber.org/zap"
)

type CatalogManagementOperator interface {
	CreateProduct(ctx context.Context, arg *catalog_service.CreateProductWithImageRequest) (*catalog_service.CreateProductResponse, error)
	GetListProduct(ctx context.Context, arg *catalog_service.ListProductRequest) (*catalog_service.ListProductResponse, error)
}

type catalogManagementOperator struct {
	catalogServiceClient catalog_service.CatalogServiceClient
	l                    logger.Interface
}

func NewCatalogManagementOperator(
	catalogServiceClient catalog_service.CatalogServiceClient,
	l logger.Interface) CatalogManagementOperator {
	return &catalogManagementOperator{
		catalogServiceClient: catalogServiceClient,
		l:                    l,
	}
}

func (c *catalogManagementOperator) CreateProduct(ctx context.Context, arg *catalog_service.CreateProductWithImageRequest) (*catalog_service.CreateProductResponse, error) {

	product, err := c.catalogServiceClient.CreateProduct(ctx, arg)
	if err != nil {
		c.l.Error("failed to create product", zap.Error(err))
		return nil, err
	}

	return product, nil
}

func (c *catalogManagementOperator) GetListProduct(ctx context.Context, arg *catalog_service.ListProductRequest) (*catalog_service.ListProductResponse, error) {
	products, err := c.catalogServiceClient.ListProduct(ctx, arg)
	if err != nil {
		c.l.Error("module.GetListProduct failed to list product", zap.Error(err))
		return nil, err
	}
	return products, nil
}
