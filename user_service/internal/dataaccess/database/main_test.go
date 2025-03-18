package database

import (
	"os"
	"testing"

	"github.com/minhhoanq/ecommerce/common/logger"
	"github.com/minhhoanq/ecommerce/user_service/config"
	"go.uber.org/zap"
)

var testQueries UserDataAccessor
var testDB Database

func TestMain(m *testing.M) {
	// config
	config, err := config.LoadConfig("../../../")
	if err != nil {
		panic(err)
	}

	logger.Setup(config.Environment, config.LogLevel)
	l := logger.NewWrapLogger(zap.DebugLevel, false)

	testDB, err = New(config, l)

	testQueries = NewUserDataAccessor(testDB, l)

	os.Exit(m.Run())
}
