package oos

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	aws_config "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/smithy-go/logging"
	. "github.com/outscale/frieza/internal/common"
	"github.com/outscale/osc-sdk-go/v3/pkg/oos"
	"github.com/outscale/osc-sdk-go/v3/pkg/profile"
	"github.com/teris-io/cli"
)

const (
	Name             = "outscale_oos"
	typeBucketObject = "object"
	typeBucket       = "bucket"
)

type OutscaleOOS struct {
	client *oos.Client
}

func New(config ProviderConfig, debug bool) (*OutscaleOOS, error) {
	profileName := config["profile"]
	profilePath := config["path"]
	profile, err := profile.NewFrom(profileName, profilePath)
	if err != nil {
		return nil, err
	}

	if ak, ok := config["ak"]; ok {
		profile.AccessKey = ak
	}

	if sk, ok := config["sk"]; ok {
		profile.SecretKey = sk
	}

	if region, ok := config["region"]; ok {
		profile.Region = region
	}

	ua := "frieza/" + FullVersion()
	opts := []aws_config.LoadOptionsFunc{aws_config.WithAppID(ua)}
	if debug {
		opts = append(opts,
			aws_config.WithClientLogMode(aws.LogRequest|aws.LogRequestWithBody|aws.LogResponseWithBody),
			aws_config.WithLogger(logging.NewStandardLogger(os.Stderr)),
			oos.WithUseragent(ua),
		)
	}
	// Note: Creating client still needs a context, but this is during initialization
	// In a future refactor, we could pass context to New() as well
	client, err := oos.NewClient(context.Background(), profile, opts...)
	if err != nil {
		return nil, err
	}

	return &OutscaleOOS{
		client: client,
	}, nil
}

func Types() []ObjectType {
	object_types := []ObjectType{
		typeBucketObject,
		typeBucket,
	}
	return object_types
}

func Cli() (string, cli.Command) {
	return Name, cli.NewCommand(Name, "create new Outscale OOS profile").
		WithOption(cli.NewOption("region", "Outscale region (e.g. eu-west-2)")).
		WithOption(cli.NewOption("ak", "access key")).
		WithOption(cli.NewOption("sk", "secret key"))
}

func (provider *OutscaleOOS) Name() string {
	return Name
}

func (provider *OutscaleOOS) Types() []ObjectType {
	return Types()
}

func (provider *OutscaleOOS) AuthTest(ctx context.Context) error {
	_, err := provider.client.ListBuckets(ctx, &s3.ListBucketsInput{})
	if err != nil {
		return errors.New("unable to list buckets")
	}
	return nil
}

func (provider *OutscaleOOS) ReadObjects(ctx context.Context, typeName string) ([]Object, error) {
	switch typeName {
	case typeBucketObject:
		return provider.readBucketObjects(ctx)
	case typeBucket:
		return provider.readBuckets(ctx)
	}
	return []Object{}, nil
}

func (provider *OutscaleOOS) DeleteObjects(ctx context.Context, typeName string, objects []Object) {
	switch typeName {
	case typeBucketObject:
		provider.deleteBucketObjects(ctx, objects)
	case typeBucket:
		provider.deleteBuckets(ctx, objects)
	}
}

func (provider *OutscaleOOS) StringObject(object string, typeName string) string {
	switch typeName {
	case typeBucketObject:
		if bucketName, key, err := decodeBucketobject(&object); err == nil {
			return bucketName + ":" + key
		}
	case typeBucket:
		if bucketName, err := decodeBucket(&object); err == nil {
			return bucketName
		}
	}
	return ""
}

func encodeBucketObject(bucketName *string, key *string) string {
	b64Bucket := base64.StdEncoding.EncodeToString([]byte(*bucketName))
	b64Key := base64.StdEncoding.EncodeToString([]byte(*key))
	encodedObject := b64Bucket + ":" + b64Key
	return encodedObject
}

func decodeBucketobject(encodedObject *string) (string, string, error) {
	bucket, key, ok := strings.Cut(*encodedObject, ":")
	if !ok {
		return "", "", errors.New("cannot decode bucket object")
	}

	binBucket, err := base64.StdEncoding.DecodeString(bucket)
	if err != nil {
		return "", "", err
	}
	bucketName := string(binBucket)

	binkey, err := base64.StdEncoding.DecodeString(key)
	if err != nil {
		return "", "", err
	}
	return bucketName, string(binkey), nil
}

func (provider *OutscaleOOS) readBucketObjects(ctx context.Context) ([]Object, error) {
	var objects []Object
	result, err := provider.client.ListBuckets(ctx, &s3.ListBucketsInput{})
	if err != nil {
		return nil, err
	}
	for _, bucket := range result.Buckets {
		result, err := provider.client.ListObjectsV2(ctx, &s3.ListObjectsV2Input{
			Bucket: bucket.Name,
		})
		if err != nil {
			continue
		}
		for _, object := range result.Contents {
			objects = append(objects, encodeBucketObject(bucket.Name, object.Key))
		}
	}
	return objects, nil
}

func (provider *OutscaleOOS) deleteBucketObjects(ctx context.Context, bucketObjects []Object) {
	for _, encodedBucketObject := range bucketObjects {
		log.Printf(
			"Deleting object: %s ... ",
			provider.StringObject(encodedBucketObject, typeBucketObject),
		)
		bucketName, key, err := decodeBucketobject(&encodedBucketObject)
		if err != nil {
			log.Println("Error while reading object details: ", err.Error())
		}
		_, err = provider.client.DeleteObject(ctx, &s3.DeleteObjectInput{
			Bucket: &bucketName,
			Key:    &key,
		})
		if err != nil {
			log.Println("Error while deleting object: ", err.Error())
		} else {
			log.Println("OK")
		}
	}
}

func encodeBucket(bucketName *string) string {
	return base64.StdEncoding.EncodeToString([]byte(*bucketName))
}

func decodeBucket(b64Bucket *string) (string, error) {
	binBucket, err := base64.StdEncoding.DecodeString(*b64Bucket)
	if err != nil {
		return "", err
	}
	bucketName := string(binBucket)
	return bucketName, nil
}

func (provider *OutscaleOOS) readBuckets(ctx context.Context) ([]Object, error) {
	var buckets []Object
	result, err := provider.client.ListBuckets(ctx, &s3.ListBucketsInput{})
	if err != nil {
		return nil, fmt.Errorf("read buckets: %w", err)
	}
	for _, b := range result.Buckets {
		buckets = append(buckets, encodeBucket(b.Name))
	}
	return buckets, nil
}

func (provider *OutscaleOOS) deleteBuckets(ctx context.Context, buckets []Object) {
	for _, b64Bucket := range buckets {
		bucketName, err := decodeBucket(&b64Bucket)
		if err != nil {
			continue
		}
		log.Printf("Deleting bucket: %s ... ", bucketName)
		_, err = provider.client.DeleteBucket(ctx, &s3.DeleteBucketInput{
			Bucket: &bucketName,
		})
		if err != nil {
			log.Println("Error while deleting bucket: ", err.Error())
		} else {
			log.Println("OK")
		}
	}
}
