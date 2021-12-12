package projects

import (
	"context"
	"github.com/crossplane/crossplane-runtime/pkg/reconciler/managed"
	"github.com/crossplane/crossplane-runtime/pkg/resource"
	"github.com/crossplane/provider-harbor/apis/harbor/project/v1alpha1"
	config "github.com/crossplane/provider-harbor/apis/v1alpha1"
	h "github.com/mittwald/goharbor-client/v5/apiv2"
	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// A Connector is expected to produce a DtrClient when its Connect method
// is called.
type Connector struct {
	kube   client.Client
	usage  resource.Tracker
	client func(baseUrl string, username string, password string) (*h.RESTClient, error)
}

// Connect typically produces a Client by:
// 1. Tracking that the managed resource is using a ProviderConfig.
// 2. Getting the managed resource's ProviderConfig.
// 3. Getting the credentials specified by the ProviderConfig.
// 4. Using the credentials to form a client.
func (c *Connector) Connect(ctx context.Context, mg resource.Managed) (managed.ExternalClient, error) {
	cr, ok := mg.(*v1alpha1.Project)
	if !ok {
		return nil, errors.New(errNotProject)
	}

	if err := c.usage.Track(ctx, mg); err != nil {
		return nil, errors.Wrap(err, errTrackPCUsage)
	}

	pc := &config.ProviderConfig{}
	if err := c.kube.Get(ctx, types.NamespacedName{Name: cr.GetProviderConfigReference().Name}, pc); err != nil {
		return nil, errors.Wrap(err, errGetPC)
	}

	un := pc.Spec.Username
	username, err := resource.CommonCredentialExtractor(ctx, un.Source, c.kube, un.CommonCredentialSelectors)
	if err != nil {
		return nil, errors.Wrap(err, errGetCreds)
	}

	pass := pc.Spec.Password
	password, err := resource.CommonCredentialExtractor(ctx, pass.Source, c.kube, pass.CommonCredentialSelectors)
	if err != nil {
		return nil, errors.Wrap(err, errGetCreds)
	}

	svc, _ := c.client(pc.Spec.BaseURL, string(username), string(password))

	return &External{client: svc}, nil
}
