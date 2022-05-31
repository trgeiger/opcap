package operator

import (
	"context"
	"encoding/json"
	"log"

	"strings"

	"fmt"

	operatorv1alpha1 "github.com/operator-framework/api/pkg/operators/v1alpha1"
	registryApi "github.com/operator-framework/operator-registry/pkg/api"
	registryClient "github.com/operator-framework/operator-registry/pkg/client"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	runtimeclient "sigs.k8s.io/controller-runtime/pkg/client"
)

type SubscriptionData struct {
	Name                   string
	Channel                string
	CatalogSource          string
	CatalogSourceNamespace string
	Package                string
}

// SubscriptionList represent the set of operators
// to be installed and tested
// It's a unique list of package/channels for operator install
func Subscriptions(catalogSource string, catalogSourceNamespace string) []SubscriptionData {

	SubscriptionList := []SubscriptionData{}

	for _, b := range bundleList() {

		s := SubscriptionData{
			Name:                   strings.Join([]string{b.PackageName, b.ChannelName, "subscription"}, "-"),
			Channel:                b.ChannelName,
			CatalogSource:          catalogSource,
			CatalogSourceNamespace: catalogSourceNamespace,
			Package:                b.PackageName,
		}
		SubscriptionList = append(SubscriptionList, s)
	}
	return uniqueElementsOf(SubscriptionList)
}

func uniqueElementsOf(s []SubscriptionData) []SubscriptionData {
	unique := make(map[SubscriptionData]bool, len(s))
	uniqueSubscriptionData := make([]SubscriptionData, len(unique))
	for _, elem := range s {
		if !unique[elem] {
			uniqueSubscriptionData = append(uniqueSubscriptionData, elem)
			unique[elem] = true
		}
	}
	return uniqueSubscriptionData
}

func InstallModesForSubscription(s SubscriptionData) []operatorv1alpha1.InstallMode {

	c, err := registryClient.NewClient("localhost:50051")
	if err != nil {
		fmt.Println(err)
	}

	bundleInChannelRequest := registryApi.GetBundleInChannelRequest{
		PkgName:     s.Package,
		ChannelName: s.Channel,
	}

	bundleInChannel, err := c.Registry.GetBundleForChannel(context.Background(), &bundleInChannelRequest)
	if err != nil {
		log.Fatal(err)
	}

	csv := operatorv1alpha1.ClusterServiceVersion{}
	err = json.Unmarshal([]byte(bundleInChannel.CsvJson), &csv)
	if err != nil {
		log.Fatal(err)
	}

	return csv.Spec.InstallModes
}

func (c operatorClient) CreateSubscription(ctx context.Context, data SubscriptionData, namespace string) (*operatorv1alpha1.Subscription, error) {
	subscription := &operatorv1alpha1.Subscription{
		ObjectMeta: metav1.ObjectMeta{
			Name:      data.Name,
			Namespace: namespace,
		},
		Spec: &operatorv1alpha1.SubscriptionSpec{
			CatalogSource:          data.CatalogSource,
			CatalogSourceNamespace: data.CatalogSourceNamespace,
			Channel:                data.Channel,
			Package:                data.Package,
		},
	}
	err := c.Client.Create(ctx, subscription)
	return subscription, err
}

func (c operatorClient) DeleteSubscription(ctx context.Context, name string, namespace string) error {
	subscription := &operatorv1alpha1.Subscription{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
	}
	return c.Client.Delete(ctx, subscription)
}

func (c operatorClient) GetSubscription(ctx context.Context, name string, namespace string) (*operatorv1alpha1.Subscription, error) {
	subscription := &operatorv1alpha1.Subscription{}
	err := c.Client.Get(ctx, runtimeclient.ObjectKey{
		Name:      name,
		Namespace: namespace,
	}, subscription)

	return subscription, err
}
