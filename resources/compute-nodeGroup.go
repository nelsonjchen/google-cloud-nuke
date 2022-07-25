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
	instanceService := compute.NewNodeGroupsService(computeService)

	resources := make([]Resource, 0)

	var pageToken string
	for {
		call := instanceService.AggregatedList(p.ID()).PageToken(pageToken)

		resp, err := call.Do()
		if err != nil {
			return nil, err
		}

		for zone, items := range resp.Items {
			for _, item := range items.NodeGroups {
				instance := &ComputeNodeGroups{
					service: instanceService,
					name:    item.Name,
					project: p.ID(),
					zone:    path.Base(zone),
				}

				resources = append(resources, instance)
			}
		}

		if pageToken = resp.NextPageToken; pageToken == "" {
			break
		}
	}

	return resources, nil
}

func (n *ComputeNodeGroups) Remove() error {
	_, err := n.service.Delete(n.project, n.zone, n.name).Do()
	if err != nil {
		return err
	}

	return err
}

func (n *ComputeNodeGroups) Properties() types.Properties {
	properties := types.NewProperties().
		Set("Name", n.name)

	return properties
}

func (n *ComputeNodeGroups) String() string {
	return n.name
}
