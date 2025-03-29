package service_test

import (
	"testing"

	"github.com/minhhoanq/ecommerce/user_service/config"
	"github.com/minhhoanq/ecommerce/user_service/internal/dataaccess/database"
	"github.com/minhhoanq/ecommerce/user_service/internal/dataaccess/kafka/producer"
	"github.com/minhhoanq/ecommerce/user_service/internal/service"
	"github.com/minhhoanq/ecommerce/user_service/internal/worker"
	"github.com/stretchr/testify/require"
)

func newTestServer(t *testing.T, store database.UserDataAccessor, taskDistributor worker.TaskDistributor, userCreatedKafkaProducer producer.UserCreatedProducer) service.UserService {
	config, err := config.LoadConfig("../../")
	require.NoError(t, err)

	service := service.NewUserService(nil, config, nil, store, taskDistributor, userCreatedKafkaProducer)
	require.NoError(t, err)

	return service
}
