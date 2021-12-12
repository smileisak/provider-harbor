package common

import (
	"time"

	h "github.com/mittwald/goharbor-client/v5/apiv2"
	clientconfig "github.com/mittwald/goharbor-client/v5/apiv2/pkg/config"
)

// NewHarborClientWithConfig instantiates a Harbor API Client
func NewHarborClientWithConfig(baseURL string, username string, password string) (*h.RESTClient, error) {
	opts := clientconfig.Options{
		Timeout:  10 * time.Second,
		PageSize: 10,
	}
	return h.NewRESTClientForHost(baseURL+"/api", username, password, &opts)
}
