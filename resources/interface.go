package resources

import (
	"fmt"
	"github.com/nelsonjchen/google-cloud-nuke/v1/pkg/config"
	"github.com/nelsonjchen/google-cloud-nuke/v1/pkg/gcputil"
	"github.com/nelsonjchen/google-cloud-nuke/v1/pkg/types"
)

type ResourceListers map[string]ResourceLister

type ResourceLister func(s *gcputil.Project) ([]Resource, error)

type Resource interface {
	Remove() error
}

type Filter interface {
	Resource
	Filter() error
}

type LegacyStringer interface {
	Resource
	String() string
}

type ResourcePropertyGetter interface {
	Resource
	Properties() types.Properties
}

type FeatureFlagGetter interface {
	Resource
	FeatureFlags(config.FeatureFlags)
}

var resourceListers = make(ResourceListers)

func register(name string, lister ResourceLister, opts ...registerOption) {
	_, exists := resourceListers[name]
	if exists {
		panic(fmt.Sprintf("a resource with the name %s already exists", name))
	}

	resourceListers[name] = lister

	for _, opt := range opts {
		opt(name, lister)
	}
}

var cloudControlMapping = map[string]string{}

func GetCloudControlMapping() map[string]string {
	return cloudControlMapping
}

type registerOption func(name string, lister ResourceLister)

func GetLister(name string) ResourceLister {
	return resourceListers[name]
}

func GetListerNames() []string {
	names := []string{}
	for resourceType := range resourceListers {
		names = append(names, resourceType)
	}

	return names
}
