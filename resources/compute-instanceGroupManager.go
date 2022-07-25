package resources

import (
	"context"
	"github.com/nelsonjchen/google-cloud-nuke/v1/pkg/gcputil"
	"github.com/nelsonjchen/google-cloud-nuke/v1/pkg/types"
	"google.golang.org/api/compute/v1"
	"path"
)

func init() {
	register("compute#instanceGroupManager", ListComputeInstanceGroupManagers)
}

type ComputeInstanceGroupManager struct {
	service *compute.InstanceGroupManagersService
	name    string
	project string
	zone    string
}

func ListComputeInstanceGroupManagers(p *gcputil.Project) ([]Resource, error) {
	ctx := context.Background()
	computeService, err := compute.NewService(ctx)
	if err != nil {
		return nil, err
	}
	service := compute.NewInstanceGroupManagersService(computeService)

	resources := make([]Resource, 0)

	var pageToken string
	for {
		call := service.AggregatedList(p.ID()).PageToken(pageToken)

		resp, err := call.Do()
		if err != nil {
			return nil, err
		}

		for zone, items := range resp.Items {
			for _, item := range items.InstanceGroupManagers {
				resources = append(resources, &ComputeInstanceGroupManager{
					service: service,
					name:    item.Name,
					project: p.ID(),
					zone:    path.Base(zone),
				})

			}
		}

		if pageToken = resp.NextPageToken; pageToken == "" {
			break
		}
	}

	return resources, nil
}

func (r *ComputeInstanceGroupManager) Remove() error {
	_, err := r.service.Delete(r.project, r.zone, r.name).Do()
	if err != nil {
		return err
	}

	return err
}

func (r *ComputeInstanceGroupManager) Properties() types.Properties {
	properties := types.NewProperties().
		Set("Name", r.name)

	return properties
}

func (r *ComputeInstanceGroupManager) String() string {
	return r.name
}
