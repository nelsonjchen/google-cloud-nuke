package resources

import (
	"context"
	"github.com/nelsonjchen/google-cloud-nuke/v1/pkg/gcputil"
	"github.com/nelsonjchen/google-cloud-nuke/v1/pkg/types"
	"google.golang.org/api/compute/v1"
	"path"
)

func init() {
	register("compute#autoscalers", ListComputeAutoscalers)
}

type ComputeAutoscaler struct {
	service *compute.AutoscalersService
	name    string
	project string
	zone    string
}

func ListComputeAutoscalers(p *gcputil.Project) ([]Resource, error) {
	ctx := context.Background()
	computeService, err := compute.NewService(ctx)
	if err != nil {
		return nil, err
	}
	service := compute.NewAutoscalersService(computeService)

	resources := make([]Resource, 0)

	var pageToken string
	for {
		call := service.AggregatedList(p.ID()).PageToken(pageToken)

		resp, err := call.Do()
		if err != nil {
			return nil, err
		}

		for zone, items := range resp.Items {
			for _, item := range items.Autoscalers {

				resources = append(resources, &ComputeAutoscaler{
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

func (r *ComputeAutoscaler) Remove() error {
	_, err := r.service.Delete(r.project, r.zone, r.name).Do()
	if err != nil {
		return err
	}

	return err
}

func (r *ComputeAutoscaler) Properties() types.Properties {
	properties := types.NewProperties().
		Set("Name", r.name)

	return properties
}

func (r *ComputeAutoscaler) String() string {
	return r.name
}
