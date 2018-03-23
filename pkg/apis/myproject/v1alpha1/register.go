package v1alpha1

import (
	"reflect"

	opkit "github.com/rook/operator-kit"
	apiextensionsv1beta1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

var (
	SchemeBuilder      = runtime.NewSchemeBuilder(addKnownTypes)
	localSchemeBuilder = &SchemeBuilder
	AddToScheme        = SchemeBuilder.AddToScheme
)

// schemeGroupVersion is group version used to register these objects
var SchemeGroupVersion = schema.GroupVersion{Group: "myproject.io", Version: "v1alpha1"}

var MySqlResource = opkit.CustomResource{
	Name:    "mysql",
	Plural:  "mysqls",
	Group:   "myproject.io",
	Version: "v1alpha1",
	Scope:   apiextensionsv1beta1.NamespaceScoped,
	Kind:    reflect.TypeOf(MySql{}).Name(),
}

func init() {
	localSchemeBuilder.Register(addKnownTypes)
}

// Resource takes an unqualified resource and returns back a Group qualified GroupResource
func Resource(resource string) schema.GroupResource {
	return SchemeGroupVersion.WithResource(resource).GroupResource()
}

// Adds the list of known types to api.Scheme.
func addKnownTypes(scheme *runtime.Scheme) error {
	scheme.AddKnownTypes(SchemeGroupVersion,
		&MySql{},
		&MySqlList{})
	metav1.AddToGroupVersion(scheme, SchemeGroupVersion)
	return nil
}
