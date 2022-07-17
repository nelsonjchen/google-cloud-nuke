package resources

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go/service/s3"
	"google.golang.org/api/iterator"

	"github.com/nelsonjchen/google-cloud-nuke/v1/pkg/gcputil"
	"time"

	"cloud.google.com/go/storage"
)

func init() {
	register("storage#bucket", ListStorageBuckets)
}

type StorageBucket struct {
	client       *storage.Client
	name         string
	creationDate time.Time
	tags         []*s3.Tag
}

func ListStorageBuckets(s *gcputil.Project) ([]Resource, error) {
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
		resources = append(resources, &StorageBucket{
			client:       client,
			name:         bucket.Name,
			creationDate: bucket.Created,
		})
	}

	return resources, nil
}

func (e *StorageBucket) Remove() error {
	ctx := context.Background()

	err := e.client.Bucket(e.name).Delete(ctx)

	return err
}

func (e *StorageBucket) String() string {
	return fmt.Sprintf("gcs://%s", e.name)
}
