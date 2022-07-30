package resources

import (
	"context"
	"fmt"
	"github.com/nelsonjchen/google-cloud-nuke/v1/pkg/gcputil"
	"github.com/nelsonjchen/google-cloud-nuke/v1/pkg/types"
	"google.golang.org/api/compute/v1"
	"google.golang.org/api/googleapi"
	"time"
)

func init() {
	register("compute#instanceTemplate", ListComputeInstanceTemplates)
}

type ComputeInstanceTemplates struct {
	service        *compute.InstanceTemplatesService
	computeService *compute.Service
	name           string
	project        string
}

func ListComputeInstanceTemplates(p *gcputil.Project) ([]Resource, error) {
	ctx := context.Background()
	computeService, err := compute.NewService(ctx)
	if err != nil {
		return nil, err
	}
	service := compute.NewInstanceTemplatesService(computeService)

	resources := make([]Resource, 0)

	var pageToken string
	for {
		call := service.List(p.ID()).PageToken(pageToken)

		resp, err := call.Do()
		if err != nil {
			return nil, err
		}

		for _, item := range resp.Items {

			resources = append(resources, &ComputeInstanceTemplates{
				service:        service,
				computeService: computeService,
				name:           item.Name,
				project:        p.ID(),
			})

		}

		if pageToken = resp.NextPageToken; pageToken == "" {
			break
		}
	}

	return resources, nil
}

func (r *ComputeInstanceTemplates) Remove() error {
	var op *compute.Operation
	op, err := r.service.Delete(r.project, r.name).Do()

	if e, ok := err.(*googleapi.Error); ok && e.Code == 404 {
		// It was already gone, so we're good.
		return nil
	}
	if err != nil {
		return err
	}

	service := compute.NewGlobalOperationsService(r.computeService)

	runningCount := 0
	for {
		op, err = service.Get(r.project, op.Name).Do()
		if op.Status == "DONE" {
			break
		}
		if op.Status == "RUNNING" {
			runningCount++
		}
		if runningCount > 4 {
			// If it is running this long, it's probably fine.
			return fmt.Errorf("operation %s is still running. will try operation again", op.Name)
		}
		if err != nil {
			return err
		}
		time.Sleep(1 * time.Second)
	}

	if op.Error != nil {
		return fmt.Errorf("error removing instance template %s: %s", r.name, op.Error.Errors[0].Message)
	}

	return err
}

func (r *ComputeInstanceTemplates) Properties() types.Properties {
	properties := types.NewProperties().
		Set("Name", r.name)

	return properties
}

func (r *ComputeInstanceTemplates) String() string {
	return r.name
}
