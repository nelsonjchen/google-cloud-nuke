package resources

import (
	"context"
	"github.com/nelsonjchen/google-cloud-nuke/v1/pkg/gcputil"
	"github.com/nelsonjchen/google-cloud-nuke/v1/pkg/types"
	"google.golang.org/api/compute/v1"
)

func init() {
	register("compute#snapshot", ListComputeSnapshots)
}

type ComputeSnapshot struct {
	service *compute.SnapshotsService
	name    string
	project string
	zone    string
}

func ListComputeSnapshots(p *gcputil.Project) ([]Resource, error) {
	ctx := context.Background()
	computeService, err := compute.NewService(ctx)
	if err != nil {
		return nil, err
	}
	service := compute.NewSnapshotsService(computeService)

	resources := make([]Resource, 0)

	var pageToken string
	for {
		call := service.List(p.ID()).PageToken(pageToken)

		resp, err := call.Do()
		if err != nil {
			return nil, err
		}

		for _, item := range resp.Items {

			resources = append(resources, &ComputeSnapshot{
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

func (r *ComputeSnapshot) Remove() error {
	_, err := r.service.Delete(r.project, r.name).Do()
	if err != nil {
		return err
	}

	return err
}

func (r *ComputeSnapshot) Properties() types.Properties {
	properties := types.NewProperties().
		Set("Name", r.name)

	return properties
}

func (r *ComputeSnapshot) String() string {
	return r.name
}
