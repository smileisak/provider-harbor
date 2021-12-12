package v1alpha1

import (
	"reflect"

	"k8s.io/apimachinery/pkg/runtime/schema"
	"sigs.k8s.io/controller-runtime/pkg/scheme"
)

const (
	KubernetesGroup = "harbor.crossplane.io"
	Version         = "v1alpha1"
)

var (
	// SchemeGroupVersion is irn version used to register these objects
	SchemeGroupVersion = schema.GroupVersion{Group: KubernetesGroup, Version: Version}

	// SchemeBuilder is used to add go types to the IrnVersionKind scheme
	SchemeBuilder = &scheme.Builder{GroupVersion: SchemeGroupVersion}

	ProjectKind                    = reflect.TypeOf(Project{}).Name()
	ProjectKubernetesGroupVersionKind = SchemeGroupVersion.WithKind(ProjectKind)
)

func init() {
	SchemeBuilder.Register(&Project{}, &ProjectList{})
}
