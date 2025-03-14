package catalog

import (
	"context"
	"fmt"

	"github.com/minhhoanq/ecommerce/common/logger"
	"github.com/minhhoanq/ecommerce/gateway/internal/generated/catalog_service"
	"go.uber.org/zap"
)

type CatalogManagementOperator interface {
	CreateProduct(ctx context.Context, arg *catalog_service.CreateProductRequest) (*catalog_service.CreateProductResponse, error)
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

func (c *catalogManagementOperator) CreateProduct(ctx context.Context, arg *catalog_service.CreateProductRequest) (*catalog_service.CreateProductResponse, error) {

	fmt.Println("arg: ", arg.Skus[0].Attributes[0].AttributeId)
	fmt.Println("daoooo")
	product, err := c.catalogServiceClient.CreateProduct(ctx, arg)
	if err != nil {
		c.l.Error("failed to create product", zap.Error(err))
		return nil, err
	}

	// fmt.Println("product", product.Product.Product.Id)

	// response := &catalog_service.CreateProductResponse{
	// 	Product: product.Product,
	// }

	return product, nil
}
