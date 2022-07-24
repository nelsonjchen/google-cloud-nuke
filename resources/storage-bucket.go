package resources

import (
	"context"
	"fmt"
	"time"

	"cloud.google.com/go/storage"
	"github.com/nelsonjchen/google-cloud-nuke/v1/pkg/gcputil"
	"github.com/nelsonjchen/google-cloud-nuke/v1/pkg/types"
	"google.golang.org/api/iterator"
)

func init() {
	register("storage#bucket", ListStorageBuckets)
}

type StorageBucket struct {
	client       *storage.Client
	name         string
	creationDate time.Time
	labels       map[string]string
}

func ListStorageBuckets(p *gcputil.Project) ([]Resource, error) {
	ctx := context.Background()
	client, err := storage.NewClient(ctx)

	bucketAttrs, err := DescribeStorageBuckets(client, p)
	if err != nil {
		return nil, err
	}

	resources := make([]Resource, 0)

	for _, bucket := range bucketAttrs {
		resources = append(resources, &StorageBucket{
			client:       client,
			name:         bucket.Name,
			creationDate: bucket.Created,
			labels:       bucket.Labels,
		})
	}

	return resources, nil
}

func DescribeStorageBuckets(s *storage.Client, p *gcputil.Project) ([]*storage.BucketAttrs, error) {
	ctx := context.Background()
	bucketIt := s.Buckets(ctx, p.ID())
	buckets := make([]*storage.BucketAttrs, 0)

	for {
		bucket, err := bucketIt.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}
		buckets = append(buckets, bucket)
	}

	return buckets, nil
}

func (b *StorageBucket) Remove() error {
	ctx := context.Background()

	bucket := b.client.Bucket(b.name)

	err := b.RemoveAllObjects()
	if err != nil {
		return err
	}

	err = bucket.Delete(ctx)

	return err
}

func (b *StorageBucket) RemoveAllObjects() error {
	ctx := context.Background()
	bucket := b.client.Bucket(b.name)
	its := bucket.Objects(ctx, &storage.Query{Versions: true})
	for {
		object, err := its.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return err
		}
		obj := bucket.Object(object.Name).Generation(object.Generation)
		err = obj.Delete(ctx)
		if err != nil {
			return err
		}
	}
	return nil
}

func (b *StorageBucket) Properties() types.Properties {
	properties := types.NewProperties().
		Set("Name", b.name).
		Set("CreationDate", b.creationDate)

	for key, label := range b.labels {
		properties.SetLabel(&key, &label)
	}

	return properties
}

func (b *StorageBucket) String() string {
	return fmt.Sprintf("gs://%s", b.name)
}
