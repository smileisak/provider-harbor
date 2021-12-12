package common

import (
	h "github.com/mittwald/goharbor-client/v5/apiv2"
	clientconfig "github.com/mittwald/goharbor-client/v5/apiv2/pkg/config"
	"time"
)

func NewHarborClientWithConfig(baseUrl string, username string, password string) (*h.RESTClient, error) {
	opts := clientconfig.Options{
		Timeout:  10 * time.Second,
		PageSize: 10,
	}
	return h.NewRESTClientForHost(baseUrl+"/api", username, password, &opts)
}
