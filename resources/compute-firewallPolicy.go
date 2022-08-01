package resources

import (
	"context"
	"github.com/nelsonjchen/google-cloud-nuke/v1/pkg/gcputil"
	"github.com/nelsonjchen/google-cloud-nuke/v1/pkg/types"
	"google.golang.org/api/compute/v1"
	"google.golang.org/api/googleapi"
)

func init() {
	register("compute#firewallPolicy", ListComputeFirewallPolicies)
}

type ComputeFirewallPolicy struct {
	service *compute.Service
	name    string
	project string
	region  string
}

func ListComputeFirewallPolicies(p *gcputil.Project) ([]Resource, error) {
	ctx := context.Background()
	service, err := compute.NewService(ctx)
	if err != nil {
		return nil, err
	}

	resources := make([]Resource, 0)

	var pageToken string
	for {
		call := service.NetworkFirewallPolicies.List(p.ID()).PageToken(pageToken)

		resp, err := call.Do()
		if err != nil {
			return nil, err
		}

		for _, item := range resp.Items {
			resources = append(resources, &ComputeFirewallPolicy{
				service: service,
				name:    item.Name,
				project: p.ID(),
				region:  item.Region,
			})
		}

		if pageToken = resp.NextPageToken; pageToken == "" {
			break
		}
	}

	// Go through each region
	// TODO: Make this run in parallel somehow since there's no aggregatedList equivalent
	regions, err := service.Regions.List(p.ID()).Do()
	if err != nil {
		return nil, err
	}
	for _, region := range regions.Items {
		pageToken = ""
		for {
			call := service.RegionNetworkFirewallPolicies.List(p.ID(), region.Name).PageToken(pageToken)

			resp, err := call.Do()
			if err != nil {
				return nil, err
			}

			for _, item := range resp.Items {
				resources = append(resources, &ComputeFirewallPolicy{
					service: service,
					name:    item.Name,
					project: p.ID(),
					region:  region.Name,
				})
			}

			if pageToken = resp.NextPageToken; pageToken == "" {
				break
			}
		}
	}

	return resources, nil
}

func (r *ComputeFirewallPolicy) Remove() error {
	var op *compute.Operation
	var err error

	if r.region != "" {
		op, err = r.service.NetworkFirewallPolicies.Delete(r.project, r.name).Do()
	} else {
		op, err = r.service.RegionNetworkFirewallPolicies.Delete(r.project, r.region, r.name).Do()
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

func (r *ComputeFirewallPolicy) Properties() types.Properties {
	properties := types.NewProperties().
		Set("Name", r.name).
		Set("Region", r.region)

	return properties
}

func (r *ComputeFirewallPolicy) String() string {
	return r.name
}
