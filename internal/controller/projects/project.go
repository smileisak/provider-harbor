package projects

import (
	"context"
	"fmt"

	xpv1 "github.com/crossplane/crossplane-runtime/apis/common/v1"
	"github.com/crossplane/crossplane-runtime/pkg/meta"
	"github.com/crossplane/crossplane-runtime/pkg/reconciler/managed"
	"github.com/crossplane/crossplane-runtime/pkg/resource"
	"github.com/crossplane/provider-harbor/apis/harbor/project/v1alpha1"
	"github.com/mittwald/goharbor-client/v5/apiv2"
	"github.com/mittwald/goharbor-client/v5/apiv2/model"
	"github.com/pkg/errors"
)

const (
	errNotProject      = "managed resource is not an Project custom resource"
	errTrackPCUsage    = "cannot track ProviderConfig usage"
	errGetPC           = "cannot get ProviderConfig"
	errGetCreds        = "cannot get credentials"
	errProjectNotFound = "project not found on server side"
)

// External observes, then either creates, updates, or deletes an
// External resource to ensure it reflects the managed resource's desired state.
type External struct {
	// A 'client' used to connect to the External resource API
	client *apiv2.RESTClient
}

// Observe runs for reconciliation
func (e *External) Observe(ctx context.Context, mg resource.Managed) (managed.ExternalObservation, error) {

	cr, ok := mg.(*v1alpha1.Project)
	if !ok {
		return managed.ExternalObservation{}, errors.New(errNotProject)
	}

	externalName := meta.GetExternalName(cr)
	if externalName == "" {
		return managed.ExternalObservation{ResourceExists: false}, nil
	}

	_, err := e.client.GetProject(ctx, externalName)
	if err != nil {
		return managed.ExternalObservation{ResourceExists: false}, nil
	}
	cr.Status.AtProvider = v1alpha1.ProjectObservation{}

	return managed.ExternalObservation{
		// Return false when the external resource does not exist. This lets
		// the managed resource reconciler know that it needs to call Create to
		// (re)create the resource, or that it has successfully been deleted.
		ResourceExists: true,

		// Return false when the external resource exists, but it not up to date
		// with the desired managed resource state. This lets the managed
		// resource reconciler know that it needs to call Update.
		ResourceUpToDate: true,

		// Return any details that may be required to connect to the external
		// resource. These will be stored as the connection secret.
		ConnectionDetails: managed.ConnectionDetails{},
	}, nil
}

// Create when CREATE events happens
func (e *External) Create(ctx context.Context, mg resource.Managed) (managed.ExternalCreation, error) {
	cr, ok := mg.(*v1alpha1.Project)
	if !ok {
		return managed.ExternalCreation{}, errors.New(errNotProject)
	}

	cr.SetConditions(xpv1.Creating())

	projectReq := &model.ProjectReq{
		CVEAllowlist: nil,
		ProjectName:  cr.Name,
		Public:       &cr.Spec.ForProvider.Public,
		StorageLimit: &cr.Spec.ForProvider.StorageLimit,
	}

	err := e.client.NewProject(ctx, projectReq)
	if err != nil {
		return managed.ExternalCreation{ExternalNameAssigned: false}, nil
	}

	meta.SetExternalName(cr, projectReq.ProjectName)

	return managed.ExternalCreation{
		// Optionally return any details that may be required to connect to the
		// external resource. These will be stored as the connection secret.
		ConnectionDetails: managed.ConnectionDetails{},
	}, nil
}

// Update when UPDATE events happens
func (e *External) Update(ctx context.Context, mg resource.Managed) (managed.ExternalUpdate, error) {
	cr, ok := mg.(*v1alpha1.Project)
	if !ok {
		return managed.ExternalUpdate{}, errors.New(errNotProject)
	}

	fmt.Printf("Updating: %+v", cr)

	return managed.ExternalUpdate{
		// Optionally return any details that may be required to connect to the
		// external resource. These will be stored as the connection secret.
		ConnectionDetails: managed.ConnectionDetails{},
	}, nil
}

// Delete runs when DELETE events happens
func (e *External) Delete(ctx context.Context, mg resource.Managed) error {
	cr, ok := mg.(*v1alpha1.Project)
	if !ok {
		return errors.New(errNotProject)
	}

	fmt.Printf("Deleting: %+v", cr)

	return nil
}
