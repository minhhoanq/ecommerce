package initial

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/hibiken/asynq"
	"github.com/minhhoanq/ecommerce/common/logger"
	"github.com/minhhoanq/ecommerce/user_service/config"
	"github.com/minhhoanq/ecommerce/user_service/internal/dataaccess/database"
	"github.com/minhhoanq/ecommerce/user_service/internal/dataaccess/kafka/producer"
	"github.com/minhhoanq/ecommerce/user_service/internal/email"
	"github.com/minhhoanq/ecommerce/user_service/internal/handler/grpc"
	"github.com/minhhoanq/ecommerce/user_service/internal/service"
	"github.com/minhhoanq/ecommerce/user_service/internal/token"
	"github.com/minhhoanq/ecommerce/user_service/internal/worker"
	"github.com/minhhoanq/ecommerce/user_service/pkg/constants"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
)

var interruptSignals = []os.Signal{
	os.Interrupt,
	syscall.SIGTERM,
	syscall.SIGINT,
}

func InitialServer(cfg config.Config, l logger.Interface) (grpc.Server, error) {
	db, err := database.New(cfg, l)
	if err != nil {
		return nil, err
	}

	ctx, stop := signal.NotifyContext(context.Background(), interruptSignals...)
	defer stop()

	kafkaProducer, err := producer.NewProducer(cfg, l)
	if err != nil {
		return nil, err
	}

	userCreatedKafkaProducer := producer.NewUserCreatedProducer(kafkaProducer, l)

	mailer := email.NewGmailSender(cfg.EmailSenderName, cfg.EmailSenderAddress, cfg.EmailSenderPassword)

	tokenMaker, err := token.NewJWTMaker(constants.PublicKeyPath, constants.PrivateKeyPath)
	if err != nil {
		l.Error("jwt maker error", zap.String("Error: ", err.Error()))
		return nil, err
	}

	redisOpts := asynq.RedisClientOpt{
		Addr: cfg.RedisAddres,
	}
	taskDistributor := worker.NewRedisTaskDistributor(l, redisOpts)

	// initialize user data access
	userDataAccessor := database.NewUserDataAccessor(db, l)
	// initial user service
	userService := service.NewUserService(tokenMaker, cfg, l, userDataAccessor, taskDistributor, userCreatedKafkaProducer)
	// initialize handler
	handler, err := grpc.NewHander(l, userService)
	if err != nil {
		return nil, err
	}
	// initialize the server
	grpcServer := grpc.NewServer(cfg, l, handler)
	waitGroup, ctx := errgroup.WithContext(ctx)

	runTaskProcessor(ctx, cfg, mailer, redisOpts, waitGroup, userDataAccessor, l)

	return grpcServer, nil
}

func runTaskProcessor(ctx context.Context, cfg config.Config, mailer email.EmailSender, redisOpt asynq.RedisClientOpt, waitGroup *errgroup.Group, userDataAccessor database.UserDataAccessor, l logger.Interface) {
	taskProcessor := worker.NewRedisTaskProcessor(redisOpt, mailer, userDataAccessor, l)
	l.Info("start task processor")
	err := taskProcessor.Start()
	if err != nil {
		l.Error("failed to start task processor", zap.String("ERROR", err.Error()))
	}

	waitGroup.Go(func() error {
		<-ctx.Done()
		l.Info("graceful shudown task processor")
		taskProcessor.Shutdown()
		l.Info("task processor stopped")
		return nil
	})
}
