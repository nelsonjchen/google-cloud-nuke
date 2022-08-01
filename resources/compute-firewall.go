package resources

import (
	"context"
	"github.com/nelsonjchen/google-cloud-nuke/v1/pkg/gcputil"
	"github.com/nelsonjchen/google-cloud-nuke/v1/pkg/types"
	"google.golang.org/api/compute/v1"
	"google.golang.org/api/googleapi"
)

func init() {
	register("compute#firewall", ListComputeFirewalls)
}

type ComputeFirewall struct {
	service *compute.Service
	name    string
	project string
	network string
}

func ListComputeFirewalls(p *gcputil.Project) ([]Resource, error) {
	ctx := context.Background()
	service, err := compute.NewService(ctx)
	if err != nil {
		return nil, err
	}

	resources := make([]Resource, 0)

	var pageToken string
	for {
		call := service.Firewalls.List(p.ID()).PageToken(pageToken)

		resp, err := call.Do()
		if err != nil {
			return nil, err
		}

		for _, item := range resp.Items {
			resources = append(resources, &ComputeFirewall{
				service: service,
				name:    item.Name,
				project: p.ID(),
				network: gcputil.Base(item.Network),
			})
		}

		if pageToken = resp.NextPageToken; pageToken == "" {
			break
		}
	}

	return resources, nil
}

func (r *ComputeFirewall) Remove() error {
	var op *compute.Operation
	var err error

	op, err = r.service.Firewalls.Delete(r.project, r.name).Do()

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

func (r *ComputeFirewall) Properties() types.Properties {
	properties := types.NewProperties().
		Set("Name", r.name).
		Set("Network", r.network)

	return properties
}

func (r *ComputeFirewall) String() string {
	return r.name
}
