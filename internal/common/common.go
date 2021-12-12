package common

import (
	"time"

	h "github.com/mittwald/goharbor-client/v5/apiv2"
	clientconfig "github.com/mittwald/goharbor-client/v5/apiv2/pkg/config"
	errors2 "github.com/mittwald/goharbor-client/v5/apiv2/pkg/errors"
	"github.com/pkg/errors"
)

// NewHarborClientWithConfig instantiates a Harbor API Client
func NewHarborClientWithConfig(baseURL string, username string, password string) (*h.RESTClient, error) {
	opts := clientconfig.Options{
		Timeout:  10 * time.Second,
		PageSize: 10,
	}
	return h.NewRESTClientForHost(baseURL+"/api", username, password, &opts)
}

// ErrorIsNotFound reports if resource not found
func ErrorIsNotFound(err error) bool {
	return errors.Is(err, errors.New(errors2.ErrProjectNotFoundMsg))
}
