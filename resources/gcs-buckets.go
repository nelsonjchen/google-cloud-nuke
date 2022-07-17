package resources

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go/service/s3"
	"google.golang.org/api/iterator"

	"github.com/rebuy-de/aws-nuke/v2/pkg/gcputil"
	"time"

	"cloud.google.com/go/storage"
)

func init() {
	register("storage", ListCSBuckets)
}

type CSBucket struct {
	client       *storage.Client
	name         string
	creationDate time.Time
	tags         []*s3.Tag
}

func ListCSBuckets(s *gcputil.Project) ([]Resource, error) {
	ctx := context.Background()
	client, err := storage.NewClient(ctx)

	bucketIt := client.Buckets(ctx, s.ID())
	if err != nil {
		return nil, err
	}

	resources := make([]Resource, 0)

	for {
		bucket, err := bucketIt.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}
		resources = append(resources, &CSBucket{
			client:       client,
			name:         bucket.Name,
			creationDate: bucket.Created,
		})
	}

	return resources, nil
}

func (e *CSBucket) Remove() error {
	ctx := context.Background()

	err := e.client.Bucket(e.name).Delete(ctx)

	return err
}

func (e *CSBucket) String() string {
	return fmt.Sprintf("gcs://%s", e.name)
}
