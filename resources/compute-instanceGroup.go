package resources

import (
	"context"
	"github.com/nelsonjchen/google-cloud-nuke/v1/pkg/gcputil"
	"github.com/nelsonjchen/google-cloud-nuke/v1/pkg/types"
	"google.golang.org/api/compute/v1"
	"path"
)

func init() {
	register("compute#instanceGroup", ListComputeInstanceGroups)
}

type ComputeInstanceGroup struct {
	service *compute.InstanceGroupsService
	name    string
	project string
	zone    string
	region  string
}

func ListComputeInstanceGroups(p *gcputil.Project) ([]Resource, error) {
	ctx := context.Background()
	computeService, err := compute.NewService(ctx)
	if err != nil {
		return nil, err
	}
	service := compute.NewInstanceGroupsService(computeService)

	resources := make([]Resource, 0)

	var pageToken string
	for {
		call := service.AggregatedList(p.ID()).PageToken(pageToken)

		resp, err := call.Do()
		if err != nil {
			return nil, err
		}

		for region, items := range resp.Items {
			for _, item := range items.InstanceGroups {
				var zone string
				if item.Zone != "" {
					zone = path.Base(item.Zone)
				}

				resources = append(resources, &ComputeInstanceGroup{
					service: service,
					name:    item.Name,
					project: p.ID(),
					region:  path.Base(region),
					zone:    zone,
				})

			}
		}

		if pageToken = resp.NextPageToken; pageToken == "" {
			break
		}
	}

	return resources, nil
}

func (r *ComputeInstanceGroup) Remove() error {
	_, err := r.service.Delete(r.project, r.zone, r.name).Do()
	if err != nil {
		return err
	}

	return err
}

func (r *ComputeInstanceGroup) Properties() types.Properties {
	properties := types.NewProperties().
		Set("Name", r.name)

	return properties
}

func (r *ComputeInstanceGroup) String() string {
	return r.name
}
