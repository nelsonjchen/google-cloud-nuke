package resources

import (
	"context"
	"github.com/nelsonjchen/google-cloud-nuke/v1/pkg/gcputil"
	"github.com/nelsonjchen/google-cloud-nuke/v1/pkg/types"
	"google.golang.org/api/compute/v1"
	"path"
)

func init() {
	register("compute#nodeTemplates", ListComputeNodeTemplates)
}

type ComputeNodeTemplate struct {
	service *compute.NodeTemplatesService
	name    string
	project string
	region  string
}

func ListComputeNodeTemplates(p *gcputil.Project) ([]Resource, error) {
	ctx := context.Background()
	computeService, err := compute.NewService(ctx)
	if err != nil {
		return nil, err
	}
	instanceService := compute.NewNodeTemplatesService(computeService)

	resources := make([]Resource, 0)

	var pageToken string
	for {
		call := instanceService.AggregatedList(p.ID()).PageToken(pageToken)

		resp, err := call.Do()
		if err != nil {
			return nil, err
		}

		for region, items := range resp.Items {
			for _, item := range items.NodeTemplates {
				resource := &ComputeNodeTemplate{
					service: instanceService,
					name:    item.Name,
					project: p.ID(),
					region:  path.Base(region),
				}

				resources = append(resources, resource)
			}
		}

		if pageToken = resp.NextPageToken; pageToken == "" {
			break
		}
	}

	return resources, nil
}

func (n *ComputeNodeTemplate) Remove() error {
	_, err := n.service.Delete(n.project, n.region, n.name).Do()
	if err != nil {
		return err
	}

	return err
}

func (n *ComputeNodeTemplate) Properties() types.Properties {
	properties := types.NewProperties().
		Set("Name", n.name)

	return properties
}

func (n *ComputeNodeTemplate) String() string {
	return n.name
}
