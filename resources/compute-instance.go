package resources

import (
	"context"
	"github.com/nelsonjchen/google-cloud-nuke/v1/pkg/gcputil"
	"github.com/nelsonjchen/google-cloud-nuke/v1/pkg/types"
	"google.golang.org/api/compute/v1"
	"google.golang.org/api/googleapi"
	"path"
)

func init() {
	register("compute#instance", ListComputeInstances)
}

type ComputeInstance struct {
	service *compute.Service
	name    string
	project string
	zone    string
	labels  map[string]string
}

func ListComputeInstances(p *gcputil.Project) ([]Resource, error) {
	ctx := context.Background()
	service, err := compute.NewService(ctx)
	if err != nil {
		return nil, err
	}

	resources := make([]Resource, 0)

	var pageToken string
	for {
		call := service.Instances.AggregatedList(p.ID()).PageToken(pageToken)

		resp, err := call.Do()
		if err != nil {
			return nil, err
		}

		for _, items := range resp.Items {
			for _, item := range items.Instances {
				resource := &ComputeInstance{
					service: service,
					name:    item.Name,
					project: p.ID(),
					zone:    path.Base(item.Zone),
				}

				for key, value := range item.Labels {
					resource.labels[key] = value
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

func (r *ComputeInstance) Remove() error {
	op, err := r.service.Instances.Delete(r.project, r.zone, r.name).Do()
	if err != nil {
		if err, ok := err.(*googleapi.Error); ok {
			if err.Code == 404 {
				return nil
			}
		}
		return err
	}
	op, err = gcputil.ComputeRemoveWaiter(op, r.service, r.project)
	if err != nil {
		return err
	}

	return err
}

func (r *ComputeInstance) Properties() types.Properties {
	properties := types.NewProperties().
		Set("Name", r.name).
		Set("Zone", r.zone)

	for key, label := range r.labels {
		properties.SetLabel(&key, &label)
	}

	return properties
}

func (r *ComputeInstance) String() string {
	return r.name
}
