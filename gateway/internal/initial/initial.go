package initial

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/labstack/echo/v4"
	echomiddleware "github.com/labstack/echo/v4/middleware"
	"github.com/minhhoanq/ecommerce/common/logger"
	"github.com/minhhoanq/ecommerce/gateway/config"
	"github.com/minhhoanq/ecommerce/gateway/internal/app/gateway"
	"github.com/minhhoanq/ecommerce/gateway/internal/app/gateway/middleware"
	catalogservice "github.com/minhhoanq/ecommerce/gateway/internal/handler/grpc/clients/catalog_service"
	orderservice "github.com/minhhoanq/ecommerce/gateway/internal/handler/grpc/clients/order_service"
	"github.com/minhhoanq/ecommerce/gateway/internal/modules/catalog"
	"github.com/minhhoanq/ecommerce/gateway/internal/modules/order"
	"github.com/minhhoanq/ecommerce/gateway/internal/routes"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
)

var interruptSignals = []os.Signal{
	os.Interrupt,
	syscall.SIGTERM,
	syscall.SIGINT,
}

func Initial(cfg config.Config, l logger.Interface) {

	// signal notify
	ctx, stop := signal.NotifyContext(context.Background(), interruptSignals...)
	defer stop()

	handler := echo.New()
	handler.Use(middleware.CORS)
	handler.Use(echomiddleware.Logger())
	handler.Use(echomiddleware.Recover())

	orderServiceGRPCClient, err := orderservice.NewClient(cfg, l)
	if err != nil {
		l.Error("failed to get instance order service grpc client", zap.Error(err))
	}

	orderOperator := order.NewOrderManagementOperator(orderServiceGRPCClient, l)

	catalogServiceGRPCClient, err := catalogservice.NewClient(cfg, l)
	if err != nil {
		l.Error("failed to get instance catalog service grpc client", zap.Error(err))
	}

	catalogOperator := catalog.NewCatalogManagementOperator(catalogServiceGRPCClient, l)

	routes.NewRouter(handler, l, orderOperator, catalogOperator)

	// wait group
	waitGroup, ctx := errgroup.WithContext(ctx)
	// start, shutdown server
	gateway.NewServer(handler, l, waitGroup, ctx)

	err = waitGroup.Wait()
	if err != nil {
		l.Error("error from wait group", zap.Error(err))
	}
}
