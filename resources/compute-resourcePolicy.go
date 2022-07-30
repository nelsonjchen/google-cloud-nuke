package resources

import (
	"context"
	"github.com/nelsonjchen/google-cloud-nuke/v1/pkg/gcputil"
	"github.com/nelsonjchen/google-cloud-nuke/v1/pkg/types"
	"google.golang.org/api/compute/v1"
	"path"
)

func init() {
	// AKA Instance Schedules
	register("compute#resourcePolicy", ListComputeResourcePolicies)
}

type ComputeResourcePolicy struct {
	service *compute.Service
	name    string
	project string
	region  string
}

func ListComputeResourcePolicies(p *gcputil.Project) ([]Resource, error) {
	ctx := context.Background()
	service, err := compute.NewService(ctx)
	if err != nil {
		return nil, err
	}

	resources := make([]Resource, 0)

	var pageToken string
	for {
		call := service.ResourcePolicies.AggregatedList(p.ID()).PageToken(pageToken)

		resp, err := call.Do()
		if err != nil {
			return nil, err
		}

		for region, items := range resp.Items {
			for _, item := range items.ResourcePolicies {
				resources = append(resources, &ComputeResourcePolicy{
					service: service,
					name:    item.Name,
					project: p.ID(),
					region:  path.Base(region),
				})
			}
		}

		if pageToken = resp.NextPageToken; pageToken == "" {
			break
		}
	}

	return resources, nil
}

func (r *ComputeResourcePolicy) Remove() error {
	op, err := r.service.ResourcePolicies.Delete(r.project, r.region, r.name).Do()
	if err != nil {
		return err
	}
	op, err = gcputil.ComputeRemoveWaiter(op, r.service, r.project)
	if err != nil {
		return err
	}

	return err
}

func (r *ComputeResourcePolicy) Properties() types.Properties {
	properties := types.NewProperties().
		Set("Name", r.name)

	return properties
}

func (r *ComputeResourcePolicy) String() string {
	return r.name
}
