package resources

import (
	"context"
	"github.com/nelsonjchen/google-cloud-nuke/v1/pkg/gcputil"
	"github.com/nelsonjchen/google-cloud-nuke/v1/pkg/types"
	"google.golang.org/api/compute/v1"
	"google.golang.org/api/googleapi"
)

func init() {
	register("compute#autoscalers", ListComputeAutoscalers)
}

type ComputeAutoscaler struct {
	service *compute.Service
	name    string
	project string
	zone    string
	region  string
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

		for _, items := range resp.Items {
			for _, item := range items.Autoscalers {
				resources = append(resources, &ComputeAutoscaler{
					service: service,
					name:    item.Name,
					project: p.ID(),
					zone:    gcputil.Base(item.Zone),
					region:  gcputil.Base(item.Region),
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
	var op *compute.Operation
	var err error
	if r.zone != "" {
		op, err = r.service.Autoscalers.Delete(r.project, r.zone, r.name).Do()
	} else {
		op, err = r.service.RegionAutoscalers.Delete(r.project, r.region, r.name).Do()
	}

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

func (r *ComputeAutoscaler) Properties() types.Properties {
	properties := types.NewProperties().
		Set("Name", r.name)

	return properties
}

func (r *ComputeAutoscaler) String() string {
	return r.name
}
