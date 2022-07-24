package resources

import (
	"context"
	"github.com/nelsonjchen/google-cloud-nuke/v1/pkg/gcputil"
	"github.com/nelsonjchen/google-cloud-nuke/v1/pkg/types"
	"google.golang.org/api/compute/v1"
	"path"
)

func init() {
	register("compute#instance", ListComputeInstances)
}

type ComputeInstance struct {
	service *compute.InstancesService
	name    string
	project string
	zone    string
	labels  map[string]string
}

func ListComputeInstances(p *gcputil.Project) ([]Resource, error) {
	ctx := context.Background()
	computeService, err := compute.NewService(ctx)
	if err != nil {
		return nil, err
	}
	instanceService := compute.NewInstancesService(computeService)

	call := instanceService.AggregatedList(p.ID())
	resp, err := call.Do()
	if err != nil {
		return nil, err
	}

	resources := make([]Resource, 0)
	for zone, items := range resp.Items {
		for _, item := range items.Instances {
			instance := &ComputeInstance{
				service: instanceService,
				name:    item.Name,
				project: p.ID(),
				zone:    path.Base(zone),
			}
			// Add labels
			for key, value := range item.Labels {
				instance.labels[key] = value
			}

			resources = append(resources, instance)
		}
	}

	return resources, nil
}

func (i *ComputeInstance) Remove() error {
	_, err := i.service.Delete(i.project, i.zone, i.name).Do()
	if err != nil {
		return err
	}

	return err
}

func (i *ComputeInstance) Properties() types.Properties {
	properties := types.NewProperties().
		Set("Name", i.name).
		Set("Zone", i.zone)

	for key, label := range i.labels {
		properties.SetLabel(&key, &label)
	}

	return properties
}

func (i *ComputeInstance) String() string {
	return i.name
}
