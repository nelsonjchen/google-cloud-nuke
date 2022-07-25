package resources

import (
	"context"
	"github.com/nelsonjchen/google-cloud-nuke/v1/pkg/gcputil"
	"github.com/nelsonjchen/google-cloud-nuke/v1/pkg/types"
	"google.golang.org/api/compute/v1"
	"path"
)

func init() {
	register("compute#nodeGroup", ListComputeNodeGroups)
}

type ComputeNodeGroups struct {
	service *compute.NodeGroupsService
	name    string
	project string
	zone    string
}

func ListComputeNodeGroups(p *gcputil.Project) ([]Resource, error) {
	ctx := context.Background()
	computeService, err := compute.NewService(ctx)
	if err != nil {
		return nil, err
	}
	service := compute.NewNodeGroupsService(computeService)

	resources := make([]Resource, 0)

	var pageToken string
	for {
		call := service.AggregatedList(p.ID()).PageToken(pageToken)

		resp, err := call.Do()
		if err != nil {
			return nil, err
		}

		for zone, items := range resp.Items {
			for _, item := range items.NodeGroups {
				resources = append(resources, &ComputeNodeGroups{
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

func (r *ComputeNodeGroups) Remove() error {
	_, err := r.service.Delete(r.project, r.zone, r.name).Do()
	if err != nil {
		return err
	}

	return err
}

func (r *ComputeNodeGroups) Properties() types.Properties {
	properties := types.NewProperties().
		Set("Name", r.name)

	return properties
}

func (r *ComputeNodeGroups) String() string {
	return r.name
}
