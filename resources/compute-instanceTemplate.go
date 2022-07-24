package resources

import (
	"cloud.google.com/go/compute/apiv1"
	"context"
	"github.com/nelsonjchen/google-cloud-nuke/v1/pkg/gcputil"
	"github.com/nelsonjchen/google-cloud-nuke/v1/pkg/types"
	"google.golang.org/api/iterator"
	computepb "google.golang.org/genproto/googleapis/cloud/compute/v1"
)

func init() {
	register("compute#instanceTemplate", ListComputeInstanceTemplates)
}

type ComputeInstanceTemplate struct {
	client  *compute.InstanceTemplatesClient
	name    string
	project string
	labels  map[string]string
}

func ListComputeInstanceTemplates(p *gcputil.Project) ([]Resource, error) {
	ctx := context.Background()
	client, err := compute.NewInstanceTemplatesRESTClient(ctx)

	if err != nil {
		return nil, err
	}

	req := &computepb.ListInstanceTemplatesRequest{
		Project: p.ID(),
	}

	resources := make([]Resource, 0)

	it := client.List(ctx, req)
	for {
		resp, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}

		resources = append(resources, &ComputeInstanceTemplate{
			client:  client,
			name:    *resp.Name,
			project: p.ID(),
		})

	}

	return resources, nil
}

func (e *ComputeInstanceTemplate) Remove() error {
	ctx := context.Background()
	_, err := e.client.Delete(ctx, &computepb.DeleteInstanceTemplateRequest{
		InstanceTemplate: e.name,
		Project:          e.project,
	})

	return err
}

func (e *ComputeInstanceTemplate) Properties() types.Properties {
	properties := types.NewProperties().
		Set("Name", e.name)

	for key, label := range e.labels {
		properties.SetLabel(&key, &label)
	}

	return properties
}
