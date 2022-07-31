package resources

import (
	"context"
	"github.com/nelsonjchen/google-cloud-nuke/v1/pkg/gcputil"
	"github.com/nelsonjchen/google-cloud-nuke/v1/pkg/types"
	"google.golang.org/api/compute/v1"
	"google.golang.org/api/googleapi"
)

func init() {
	register("compute#instanceGroupManager", ListComputeInstanceGroupManagers)
}

type ComputeInstanceGroupManager struct {
	service *compute.Service
	name    string
	project string
	region  string
	zone    string
}

func ListComputeInstanceGroupManagers(p *gcputil.Project) ([]Resource, error) {
	ctx := context.Background()
	service, err := compute.NewService(ctx)
	if err != nil {
		return nil, err
	}

	resources := make([]Resource, 0)

	var pageToken string
	for {
		call := service.InstanceGroupManagers.AggregatedList(p.ID()).PageToken(pageToken)

		resp, err := call.Do()
		if err != nil {
			return nil, err
		}

		for _, items := range resp.Items {
			for _, item := range items.InstanceGroupManagers {
				resources = append(resources, &ComputeInstanceGroupManager{
					service: service,
					name:    item.Name,
					project: p.ID(),
					region:  gcputil.Base(item.Region),
					zone:    gcputil.Base(item.Zone),
				})

			}
		}

		if pageToken = resp.NextPageToken; pageToken == "" {
			break
		}
	}

	return resources, nil
}

func (r *ComputeInstanceGroupManager) Remove() error {
	var op *compute.Operation
	var err error
	if r.zone != "" {
		op, err = r.service.InstanceGroupManagers.Delete(r.project, r.zone, r.name).Do()
	} else {
		op, err = r.service.RegionInstanceGroupManagers.Delete(r.project, r.region, r.name).Do()
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

func (r *ComputeInstanceGroupManager) Properties() types.Properties {
	properties := types.NewProperties().
		Set("Name", r.name)

	return properties
}

func (r *ComputeInstanceGroupManager) String() string {
	return r.name
}
