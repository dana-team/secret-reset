package token

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/dana-team/secretreset/internal/util"
	"net/http"
	"os"
	"strings"

	"github.com/dana-team/secretreset/internal/clients"
	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"

	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	errParsing            = "Failed to parse JSON"
	errNotSet             = "Error: Please set "
	errUpdating           = "Failed to update resource"
	errCreating           = "Failed to create a new resource"
	errCreatingSecret     = "Failed to create a new resource"
	errUpdatingSecret     = "Failed to update secret"
	errGettingResource    = "Failed getting the resource"
	errSendingHTTPRequest = "Failed sending the http request"
	errExtractingToken    = "Failed to extract token"
)

const (
	authUsername     = "AUTH_USERNAME"
	authClientSecret = "AUTH_CLIENT_SECRET"
	authUrl          = "AUTH_URL"
	secretName       = "SECRET_NAME"
	secretNamespace  = "SECRET_NAMESPACE"
	token            = "token"
	grantType        = "?grant_type=client_credentials"
)

type Manager struct {
	Logger     logr.Logger
	K8sClient  client.Client
	HTTPClient *http.Client
}

// extractAccessToken parses the access token from a json response body.
func (t *Manager) extractAccessToken(body []byte) (string, error) {
	var token AccessToken
	if err := json.Unmarshal(body, &token); err != nil {
		t.Logger.Error(err, errParsing)
		return "", fmt.Errorf("%s: %w", errParsing, err)
	}
	return token.Token, nil
}

// updateSecret updates an existing secret with the given access token.
func (t *Manager) updateSecret(accessToken string, secret *corev1.Secret, secretName string, secretNamespace string) error {
	t.Logger.Info("Secret %q exists in namespace %q", secretName, secretNamespace)

	encodedToken := base64.StdEncoding.EncodeToString([]byte(accessToken))

	if secret.Data == nil {
		secret.Data = make(map[string][]byte)
	}

	secret.Data[token] = []byte(encodedToken)

	if err := clients.UpdateResource(t.K8sClient, secret); err != nil {
		return fmt.Errorf("%s: %w", errUpdatingSecret, err)
	}

	t.Logger.Info("Secret %q exists in namespace %q, updating...", secretName, secretNamespace)
	return nil
}

// prepareSecret creates a new corev1 secret.
func prepareSecret(secretName string, secretNamespace string, accessToken string) *corev1.Secret {
	newSecret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      secretName,
			Namespace: secretNamespace,
		},
		StringData: map[string]string{
			token: accessToken,
		},
	}

	return newSecret
}

// createSecret creates a new secret that stores the given access token.
func (t *Manager) createSecret(secretName string, secretNamespace string, accessToken string) error {
	t.Logger.Info("Secret %q does not exist in namespace %q. Creating...", secretName, secretNamespace)
	newSecret := prepareSecret(secretName, secretNamespace, accessToken)

	if err := clients.CreateResource(t.K8sClient, newSecret); err != nil {
		return fmt.Errorf("%s: %w", errCreatingSecret, err)
	}

	t.Logger.Info("Secret %q created successfully in namespace %q", secretName, secretNamespace)
	return nil
}

// buildAuthParams creates map of all the authentication related environment variables
func buildAuthParams() map[string]string {
	return map[string]string{
		authUsername:     os.Getenv(authUsername),
		authClientSecret: os.Getenv(authClientSecret),
		authUrl:          fmt.Sprintf("%s%s", os.Getenv(authUrl), grantType),
	}
}

// buildSecretParams creates map of all the secret related environment variables
func buildSecretParams() map[string]string {
	return map[string]string{
		secretName:      os.Getenv(secretName),
		secretNamespace: os.Getenv(secretNamespace),
	}
}

// CreateOrUpdate handles storing an access token in a secret (by updating the secret or creating a new one).
func (t *Manager) CreateOrUpdate() error {
	requiredVariables := []string{authUsername, authClientSecret, authUrl, secretName, secretNamespace}
	missingVariables := util.CheckRequiredVariables(requiredVariables)
	if len(missingVariables) > 0 {
		t.Logger.Error(nil, "%s %q", errNotSet, strings.Join(missingVariables, ","))
		return fmt.Errorf("%s %s", errNotSet, strings.Join(missingVariables, ","))
	}

	authParams := buildAuthParams()
	secretParams := buildSecretParams()
	authHeader := util.EncodeResource(authParams[authUsername], authParams[authClientSecret])

	response, err := clients.SendRequest(authParams[authUrl], authHeader, t.Logger, t.HTTPClient)
	if err != nil {
		return fmt.Errorf(errSendingHTTPRequest)
	}

	accessToken, err := t.extractAccessToken(response)
	if err != nil {
		return fmt.Errorf("%s: %w", errExtractingToken, err)
	}

	secret := &corev1.Secret{}
	if err = t.K8sClient.Get(context.TODO(), types.NamespacedName{Name: secretParams[secretName], Namespace: secretParams[secretNamespace]}, secret); err == nil {
		if err = t.updateSecret(accessToken, secret, secretParams[secretName], secretParams[secretNamespace]); err != nil {
			t.Logger.Error(err, errUpdating)
			return fmt.Errorf("%s: %w", errUpdating, err)
		}
	} else if apierrors.IsNotFound(err) {
		if err = t.createSecret(secretParams[secretName], secretParams[secretNamespace], accessToken); err != nil {
			t.Logger.Error(err, errCreating)
			return fmt.Errorf("%s: %w", errCreating, err)
		}
	} else {
		t.Logger.Error(err, errGettingResource)
		return fmt.Errorf("%s: %w", errGettingResource, err)
	}
	return nil
}
