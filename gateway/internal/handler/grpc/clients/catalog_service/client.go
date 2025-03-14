package orderservice

import (
	"github.com/minhhoanq/ecommerce/common/logger"
	"github.com/minhhoanq/ecommerce/gateway/config"
	"github.com/minhhoanq/ecommerce/gateway/internal/generated/catalog_service"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func NewClient(cfg config.Config, l logger.Interface) (catalog_service.CatalogServiceClient, error) {
	var opts = []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}

	conn, err := grpc.NewClient(cfg.GRPCCatalogAddress, opts...)
	if err != nil {
		l.Error("failed to connect catalog service grpc client", zap.Error(err))
		return nil, err
	}

	client := catalog_service.NewCatalogServiceClient(conn)
	l.Info("connectiont to catalog service grpc client", zap.String("Address: ", cfg.GRPCCatalogAddress))
	return client, err
}
