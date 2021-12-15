package projects

import (
    "context"
    "testing"

    xpv1 "github.com/crossplane/crossplane-runtime/apis/common/v1"
    "github.com/crossplane/crossplane-runtime/pkg/meta"
    "github.com/crossplane/crossplane-runtime/pkg/reconciler/managed"
    "github.com/crossplane/crossplane-runtime/pkg/resource"
    "github.com/crossplane/crossplane-runtime/pkg/test"
    "github.com/crossplane/provider-harbor/apis/harbor/project/v1alpha1"
    "github.com/google/go-cmp/cmp"
    "github.com/mittwald/goharbor-client/v5/apiv2"
    h "github.com/mittwald/goharbor-client/v5/apiv2"
    "github.com/mittwald/goharbor-client/v5/apiv2/mocks"
    clienttesting "github.com/mittwald/goharbor-client/v5/apiv2/pkg/testing"
    "github.com/pkg/errors"
)

const (
    baseUrl                 = "https://core.harbor.crossplane.io"
    dummyUsernameOrPassword = "dummy"
)

var (
    unexpectedItem resource.Managed
)

type fields struct {
    client *apiv2.RESTClient
}

type args struct {
    ctx context.Context
    mg  resource.Managed
}

type projectModifier func(project *v1alpha1.Project)

func withClientDefaultValues() projectModifier {
    return func(o *v1alpha1.Project) {
        o.Spec.ForProvider = v1alpha1.ProjectParameters{
            Public:       true,
            StorageLimit: 0,
        }
    }
}

func withExternalName(n string) projectModifier {
    return func(o *v1alpha1.Project) { meta.SetExternalName(o, n) }
}
func withAtProvider(a v1alpha1.ProjectObservation) projectModifier {
    return func(o *v1alpha1.Project) { o.Status.AtProvider = a }
}

func withConditions(c ...xpv1.Condition) projectModifier {
    return func(cr *v1alpha1.Project) { cr.Status.ConditionedStatus.Conditions = c }
}

func user(m ...projectModifier) *v1alpha1.Project {
    cr := &v1alpha1.Project{}
    for _, f := range m {
        f(cr)
    }
    return cr
}

func APIandMockClientsForTests() (*h.RESTClient, *clienttesting.MockClients) {
    desiredMockClients := &clienttesting.MockClients{
        Project: mocks.MockProjectClientService{},
    }

    v2Client := clienttesting.BuildV2ClientWithMocks(desiredMockClients)
    cl := apiv2.NewRESTClient(v2Client, clienttesting.DefaultOpts, clienttesting.AuthInfo)

    return cl, desiredMockClients
}

func TestObserve(t *testing.T) {
    mockedClient, _ := APIandMockClientsForTests()
    type want struct {
        o   managed.ExternalObservation
        err error
    }
    cases := map[string]struct {
        reason string
        fields fields
        args   args
        want   want
    }{
        "InValidInput": {
            reason: "Invalid Input",
            fields: fields{client: mockedClient},
            args: args{
                ctx: context.Background(),
                mg:  unexpectedItem,
            },
            want: want{
                o:   managed.ExternalObservation{},
                err: errors.New(errNotProject),
            },
        },
    }
    for name, tc := range cases {

        t.Run(name, func(t *testing.T) {
            e := External{client: tc.fields.client}
            got, err := e.Observe(tc.args.ctx, tc.args.mg)

            if diff := cmp.Diff(tc.want.err, err, test.EquateErrors()); diff != "" {
                t.Errorf("\n%s\ne.Observe(...): -want error, +got error:\n%s\n", tc.reason, diff)
            }

            if diff := cmp.Diff(tc.want.o, got); diff != "" {
                t.Errorf("\n%s\ne.Observe(...): -want, +got:\n%s\n", tc.reason, diff)
            }
        })
    }
}
