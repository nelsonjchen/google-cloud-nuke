package resources

import (
	"context"
	"github.com/nelsonjchen/google-cloud-nuke/v1/pkg/gcputil"
	"github.com/nelsonjchen/google-cloud-nuke/v1/pkg/types"
	"google.golang.org/api/compute/v1"
)

func init() {
	register("compute#image", ListComputeImages)
}

type ComputeImage struct {
	service *compute.ImagesService
	name    string
	project string
}

func ListComputeImages(p *gcputil.Project) ([]Resource, error) {
	ctx := context.Background()
	computeService, err := compute.NewService(ctx)
	if err != nil {
		return nil, err
	}
	service := compute.NewImagesService(computeService)

	resources := make([]Resource, 0)

	var pageToken string
	for {
		call := service.List(p.ID()).PageToken(pageToken)

		resp, err := call.Do()
		if err != nil {
			return nil, err
		}

		for _, item := range resp.Items {

			resources = append(resources, &ComputeImage{
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

func (r *ComputeImage) Remove() error {
	_, err := r.service.Delete(r.project, r.name).Do()
	if err != nil {
		return err
	}

	return err
}

func (r *ComputeImage) Properties() types.Properties {
	properties := types.NewProperties().
		Set("Name", r.name)

	return properties
}

func (r *ComputeImage) String() string {
	return r.name
}
