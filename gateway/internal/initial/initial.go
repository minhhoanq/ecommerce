package initial

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/labstack/echo/v4"
	"github.com/minhhoanq/ecommerce/common/logger"
	"github.com/minhhoanq/ecommerce/gateway/config"
	"github.com/minhhoanq/ecommerce/gateway/internal/gateway"
	"github.com/minhhoanq/ecommerce/gateway/internal/gateway/middleware"
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

	// wait group
	waitGroup, ctx := errgroup.WithContext(ctx)
	// start, shutdown server
	gateway.NewServer(handler, l, waitGroup, ctx)

	err := waitGroup.Wait()
	if err != nil {
		l.Error("error from wait group", zap.Error(err))
	}
}
