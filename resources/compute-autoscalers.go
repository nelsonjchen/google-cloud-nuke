package resources

import (
	"context"
	"fmt"
	"github.com/nelsonjchen/google-cloud-nuke/v1/pkg/gcputil"
	"github.com/nelsonjchen/google-cloud-nuke/v1/pkg/types"
	"google.golang.org/api/compute/v1"
	"path"
)

func init() {
	register("compute#autoscalers", ListComputeAutoscalers)
}

type ComputeAutoscaler struct {
	service   *compute.Service
	name      string
	project   string
	zone      string
	operation *compute.Operation
}

func ListComputeAutoscalers(p *gcputil.Project) ([]Resource, error) {
	ctx := context.Background()
	service, err := compute.NewService(ctx)
	if err != nil {
		return nil, err
	}

	resources := make([]Resource, 0)

	var pageToken string
	for {
		call := service.Autoscalers.AggregatedList(p.ID()).PageToken(pageToken)

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
	if r.operation != nil {
		_, err := gcputil.ComputeRemoveWaiter(r.operation, r.service, r.project)
		if err != nil {
			// Try deleting again on next poll
			r.operation = nil
			return err
		}
		// Operation is done (resource already gone), pending or running
		return nil
	}
	op, err := r.service.Autoscalers.Delete(r.project, r.zone, r.name).Do()
	if err != nil {
		// It's already gone, that's great.
		if op.HTTPStatusCode == 404 {
			return nil
		}
		return err
	}
	r.operation = op

	if op.Status == "RUNNING" || op.Status == "PENDING" {
		return fmt.Errorf("operation is running")
	}

	return nil
}

func (r *ComputeAutoscaler) Properties() types.Properties {
	properties := types.NewProperties().
		Set("Name", r.name)

	return properties
}

func (r *ComputeAutoscaler) String() string {
	return r.name
}
