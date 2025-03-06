package grpc

import (
	"context"

	"github.com/minhhoanq/ecommerce/user_service/config"
	pbuser "github.com/minhhoanq/ecommerce/user_service/internal/controller/grpc/v1/user_service"
	"github.com/minhhoanq/ecommerce/user_service/internal/token"
	"github.com/minhhoanq/ecommerce/user_service/internal/usecase/rest/repo"
	"github.com/minhhoanq/ecommerce/user_service/internal/worker"
	"github.com/minhhoanq/ecommerce/user_service/pkg/kafka"
)

type GrpcServer struct {
	pbuser.UnimplementedUserServiceServer
	cfg             config.Config
	ctx             context.Context
	taskDistributor worker.TaskDistributor
	q               repo.Querier
	tokenMaker      token.Maker
	kafkaProducer   kafka.Producer
}

func NewGrpcServer(cfg config.Config, ctx context.Context, taskDistributor worker.TaskDistributor, q repo.Querier, tokenMaker token.Maker, kafkaProducer kafka.Producer) (*GrpcServer, error) {
	return &GrpcServer{
		cfg:             cfg,
		ctx:             ctx,
		taskDistributor: taskDistributor,
		q:               q,
		tokenMaker:      tokenMaker,
		kafkaProducer:   kafkaProducer,
	}, nil
}
