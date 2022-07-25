package resources

import (
	"context"
	"github.com/nelsonjchen/google-cloud-nuke/v1/pkg/gcputil"
	"github.com/nelsonjchen/google-cloud-nuke/v1/pkg/types"
	"google.golang.org/api/compute/v1"
)

func init() {
	register("compute#instanceTemplate", ListComputeInstanceTemplates)
}

type ComputeInstanceTemplates struct {
	service *compute.InstanceTemplatesService
	name    string
	project string
}

func ListComputeInstanceTemplates(p *gcputil.Project) ([]Resource, error) {
	ctx := context.Background()
	computeService, err := compute.NewService(ctx)
	if err != nil {
		return nil, err
	}
	instanceService := compute.NewInstanceTemplatesService(computeService)

	resources := make([]Resource, 0)

	var pageToken string
	for {
		call := instanceService.List(p.ID()).PageToken(pageToken)

		resp, err := call.Do()
		if err != nil {
			return nil, err
		}

		for _, item := range resp.Items {

			instance := &ComputeInstanceTemplates{
				service: instanceService,
				name:    item.Name,
				project: p.ID(),
			}

			resources = append(resources, instance)

		}

		if pageToken = resp.NextPageToken; pageToken == "" {
			break
		}
	}

	return resources, nil
}

func (t *ComputeInstanceTemplates) Remove() error {
	_, err := t.service.Delete(t.project, t.name).Do()
	if err != nil {
		return err
	}

	return err
}

func (t *ComputeInstanceTemplates) Properties() types.Properties {
	properties := types.NewProperties().
		Set("Name", t.name)

	return properties
}

func (t *ComputeInstanceTemplates) String() string {
	return t.name
}
