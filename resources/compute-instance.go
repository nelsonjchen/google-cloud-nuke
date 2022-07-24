package resources

import (
	"context"
	"path"

	"cloud.google.com/go/compute/apiv1"
	"github.com/nelsonjchen/google-cloud-nuke/v1/pkg/gcputil"
	"github.com/nelsonjchen/google-cloud-nuke/v1/pkg/types"
	"google.golang.org/api/iterator"
	computepb "google.golang.org/genproto/googleapis/cloud/compute/v1"
)

func init() {
	register("compute#instance", ListComputeInstances)
}

type ComputeInstance struct {
	client  *compute.InstancesClient
	name    string
	project string
	zone    string
	labels  map[string]string
}

func ListComputeInstances(p *gcputil.Project) ([]Resource, error) {
	ctx := context.Background()
	client, err := compute.NewInstancesRESTClient(ctx)
	if err != nil {
		return nil, err
	}

	req := &computepb.AggregatedListInstancesRequest{
		Project: p.ID(),
	}

	resources := make([]Resource, 0)

	it := client.AggregatedList(ctx, req)
	for {
		resp, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}

		for _, computeInstance := range resp.Value.Instances {
			resources = append(resources, &ComputeInstance{
				client:  client,
				name:    *computeInstance.Name,
				zone:    path.Base(*computeInstance.Zone),
				project: p.ID(),
				labels:  computeInstance.Labels,
			})
		}
	}

	return resources, nil
}

func (e *ComputeInstance) Remove() error {
	ctx := context.Background()
	_, err := e.client.Delete(ctx, &computepb.DeleteInstanceRequest{
		Instance: e.name,
		Project:  e.project,
		Zone:     e.zone,
	})

	return err
}

func (e *ComputeInstance) Properties() types.Properties {
	properties := types.NewProperties().
		Set("Name", e.name).
		Set("Zone", e.zone)

	for key, label := range e.labels {
		properties.SetLabel(&key, &label)
	}

	return properties
}
