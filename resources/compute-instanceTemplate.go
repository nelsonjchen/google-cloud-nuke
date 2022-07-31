package resources

import (
	"context"
	"fmt"
	"github.com/nelsonjchen/google-cloud-nuke/v1/pkg/gcputil"
	"github.com/nelsonjchen/google-cloud-nuke/v1/pkg/types"
	"google.golang.org/api/compute/v1"
)

func init() {
	register("compute#instanceTemplate", ListComputeInstanceTemplates)
}

type ComputeInstanceTemplates struct {
	service   *compute.Service
	name      string
	project   string
	operation *compute.Operation
}

func ListComputeInstanceTemplates(p *gcputil.Project) ([]Resource, error) {
	ctx := context.Background()
	service, err := compute.NewService(ctx)
	if err != nil {
		return nil, err
	}

	resources := make([]Resource, 0)

	var pageToken string
	for {
		call := service.InstanceTemplates.List(p.ID()).PageToken(pageToken)

		resp, err := call.Do()
		if err != nil {
			return nil, err
		}

		for _, item := range resp.Items {

			resources = append(resources, &ComputeInstanceTemplates{
				service: service,
				name:    item.Name,
				project: p.ID(),
			})

		}

		if pageToken = resp.NextPageToken; pageToken == "" {
			break
		}
	}

	return resources, nil
}

func (r *ComputeInstanceTemplates) Remove() error {
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
	op, err := r.service.InstanceTemplates.Delete(r.project, r.name).Do()
	if err != nil {
		// It's already gone, that's great.
		if op != nil {
			if op.HTTPStatusCode == 404 {
				return nil
			}
		}
		return err
	}
	r.operation = op

	if op.Status == "RUNNING" || op.Status == "PENDING" {
		return fmt.Errorf("operation is running")
	}

	return nil
}

func (r *ComputeInstanceTemplates) Properties() types.Properties {
	properties := types.NewProperties().
		Set("Name", r.name)

	return properties
}

func (r *ComputeInstanceTemplates) String() string {
	return r.name
}
