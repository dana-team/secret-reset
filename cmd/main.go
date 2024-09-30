package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/dana-team/secretreset/internal/clients"
	"github.com/dana-team/secretreset/internal/token"
	"github.com/dana-team/secretreset/internal/util"
	"github.com/go-logr/logr"
	"github.com/go-logr/zapr"
	"go.uber.org/zap"
)

const (
	errInitializingLogger  = "failed initializing logger"
	errInitializingClient  = "failed initializing client"
	errCreateOrUpdateToken = "failed creating or updating token"
)

func main() {
	logger, err := initializeLogger()
	if err != nil {
		panic(fmt.Errorf("%s: %v", errInitializingLogger, err))
	}

	k8sClient, err := clients.Initialize(logger)
	if err != nil {
		logger.Error(err, errInitializingClient)
		os.Exit(1)
	}

	tokenManager := token.Manager{
		Logger:     logger,
		K8sClient:  k8sClient,
		HTTPClient: &http.Client{Transport: util.GetTransportSettings()},
	}

	err = tokenManager.CreateOrUpdate()
	if err != nil {
		logger.Error(err, errCreateOrUpdateToken)
		os.Exit(1)
	}

}

// initializeLogger initializes a new logger.
func initializeLogger() (logr.Logger, error) {
	zapLogger, err := zap.NewProduction()
	if err != nil {
		return logr.Logger{}, err
	}

	logger := zapr.NewLogger(zapLogger)
	return logger, nil
}
