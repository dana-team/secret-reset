package clients

import (
	"context"
	"fmt"

	"github.com/go-logr/logr"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
)

const (
	errGettingConfig = "unable to get the config"
	errGettingClient = "unable to create a new client"
)

// Initialize initializes a new client.
func Initialize(logger logr.Logger) (client.Client, error) {
	cfg, err := config.GetConfig()
	if err != nil {
		logger.Error(err, errGettingConfig)
		return nil, fmt.Errorf("%s: %w", errGettingConfig, err)
	}

	k8sClient, err := client.New(cfg, client.Options{})
	if err != nil {
		logger.Error(err, errGettingClient)
		return nil, fmt.Errorf("%s: %w", errGettingClient, err)
	}

	return k8sClient, nil
}

// CreateResource creates a new secret.
func CreateResource(K8sClient client.Client, resource client.Object) error {
	if err := K8sClient.Create(context.TODO(), resource); err != nil {
		return err
	}
	return nil
}

// UpdateResource updates an existing secret.
func UpdateResource(K8sClient client.Client, resource client.Object) error {
	if err := K8sClient.Update(context.TODO(), resource); err != nil {
		return err
	}
	return nil
}
