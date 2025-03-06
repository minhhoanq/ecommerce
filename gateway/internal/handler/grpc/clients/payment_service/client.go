package orderservice

import (
	"github.com/minhhoanq/ecommerce/common/logger"
	"github.com/minhhoanq/ecommerce/gateway/config"
	"github.com/minhhoanq/ecommerce/gateway/internal/generated/payment_service"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func NewClient(cfg config.Config, l logger.Interface) (payment_service.PaymentServiceClient, error) {
	var opts = []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}

	conn, err := grpc.NewClient(cfg.GRPCOrderAddress, opts...)
	if err != nil {
		l.Error("failed to connect order service grpc client", zap.Error(err))
		return nil, err
	}

	client := payment_service.NewPaymentServiceClient(conn)
	return client, err
}
