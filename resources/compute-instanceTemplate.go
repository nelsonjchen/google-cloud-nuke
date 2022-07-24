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

func (t *ComputeInstanceTemplate) Remove() error {
	ctx := context.Background()
	_, err := t.client.Delete(ctx, &computepb.DeleteInstanceTemplateRequest{
		InstanceTemplate: t.name,
		Project:          t.project,
	})

	return err
}

func (t *ComputeInstanceTemplate) Properties() types.Properties {
	properties := types.NewProperties().
		Set("Name", t.name)

	for key, label := range t.labels {
		properties.SetLabel(&key, &label)
	}

	return properties
}

func (t *ComputeInstanceTemplate) String() string {
	return t.name
}
