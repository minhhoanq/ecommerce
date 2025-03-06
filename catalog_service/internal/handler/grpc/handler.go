package grpc

import (
	"context"
	"fmt"

	pb "github.com/minhhoanq/ecommerce/catalog_service/internal/generated/catalog_service"
	"github.com/minhhoanq/ecommerce/catalog_service/internal/service"
	"github.com/minhhoanq/ecommerce/common/logger"
	"github.com/redis/go-redis/v9"
)

type Handler struct {
	pb.UnimplementedCatalogServiceServer
	catalogService service.CatalogService
	l              logger.Interface
	redisClient    *redis.Client
}

func NewHandler(catalogService service.CatalogService, l logger.Interface, redisClient *redis.Client) (pb.CatalogServiceServer, error) {
	return &Handler{
		catalogService: catalogService,
		l:              l,
		redisClient:    redisClient,
	}, nil
}

func (h *Handler) CreateProduct(ctx context.Context, arg *pb.CreateProductRequest) (*pb.CreateProductResponse, error) {
	h.l.Info("Create product in handler")
	fmt.Println("handle create product")
	product, err := h.catalogService.CreateProduct(ctx, arg)
	if err != nil {
		return nil, err
	}
	return product, nil
}

func (h *Handler) ListProduct(ctx context.Context, arg *pb.ListProductRequest) (*pb.ListProductResponse, error) {
	// if h.redisClient.Get(context.Background(), "products").Val() != "" {
	// 	h.l.Info("Get products from redis")
	// 	// convert string to byte
	// 	bytesProducts, err := h.redisClient.Get(context.Background(), "products").Bytes()
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// 	products := &pb.ListProductResponse{}
	// 	err = json.Unmarshal(bytesProducts, &products)
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// 	return products, nil
	// }
	// get products from database
	products, err := h.catalogService.ListProduct(ctx, arg)
	// jsonProducts, err := json.Marshal(products)
	// if err != nil {
	// 	return nil, err
	// }

	// // set products to redis
	// h.redisClient.Set(context.Background(), "products", jsonProducts, 0)
	// h.l.Info("Get products from database")

	if err != nil {
		return nil, err
	}
	return products, nil
}

func (h *Handler) CreateCart(ctx context.Context, arg *pb.CreateCartRequest) (*pb.CreateCartResponse, error) {
	cart, err := h.catalogService.CreateCart(ctx, arg)
	if err != nil {
		return nil, err
	}
	return cart, nil
}

func (h *Handler) AddToCartItem(ctx context.Context, arg *pb.AddToCartItemRequest) (*pb.AddToCartItemResponse, error) {
	cart, err := h.catalogService.AddToCartItem(ctx, arg)
	if err != nil {
		return nil, err
	}
	return cart, nil
}

func (h *Handler) GetSKU(ctx context.Context, arg *pb.GetSKURequest) (*pb.GetSKUResponse, error) {
	return h.catalogService.GetSKU(ctx, arg)
}

func (h *Handler) GetInventorySKU(ctx context.Context, arg *pb.GetInventorySKURequest) (*pb.GetInventorySKUResponse, error) {
	return h.catalogService.GetInventorySKU(ctx, arg)
}

func (h *Handler) UpdateInventorySKU(ctx context.Context, arg *pb.UpdateInventorySKURequest) (*pb.UpdateInventorySKUResponse, error) {
	return h.catalogService.UpdateInventorySKU(ctx, arg)
}
