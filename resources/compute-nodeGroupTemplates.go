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
	register("compute#nodeTemplates", ListComputeNodeTemplates)
}

type ComputeNodeTemplate struct {
	service   *compute.Service
	name      string
	project   string
	region    string
	operation *compute.Operation
}

func ListComputeNodeTemplates(p *gcputil.Project) ([]Resource, error) {
	ctx := context.Background()
	service, err := compute.NewService(ctx)
	if err != nil {
		return nil, err
	}

	resources := make([]Resource, 0)

	var pageToken string
	for {
		call := service.NodeTemplates.AggregatedList(p.ID()).PageToken(pageToken)

		resp, err := call.Do()
		if err != nil {
			return nil, err
		}

		for region, items := range resp.Items {
			for _, item := range items.NodeTemplates {
				resources = append(resources, &ComputeNodeTemplate{
					service: service,
					name:    item.Name,
					project: p.ID(),
					region:  path.Base(region),
				})
			}
		}

		if pageToken = resp.NextPageToken; pageToken == "" {
			break
		}
	}

	return resources, nil
}

func (r *ComputeNodeTemplate) Remove() error {
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
	op, err := r.service.NodeTemplates.Delete(r.project, r.region, r.name).Do()
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

func (r *ComputeNodeTemplate) Properties() types.Properties {
	properties := types.NewProperties().
		Set("Name", r.name)

	return properties
}

func (r *ComputeNodeTemplate) String() string {
	return r.name
}
