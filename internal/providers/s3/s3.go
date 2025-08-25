package s3

import (
	"encoding/base64"
	"errors"
	"log"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	. "github.com/outscale/frieza/internal/common"
	"github.com/teris-io/cli"
)

const (
	Name             = "s3"
	typeBucketObject = "object"
	typeBucket       = "bucket"
)

type S3 struct {
	client *s3.S3
}

func checkConfig(config ProviderConfig) error {
	if len(config["endpoint"]) == 0 {
		return errors.New("endoint is needed")
	}
	if len(config["region"]) == 0 {
		return errors.New("region's name is needed")
	}
	if len(config["ak"]) == 0 {
		return errors.New("access key is needed")
	}
	if len(config["sk"]) == 0 {
		return errors.New("secret key is needed")
	}
	return nil
}

func New(config ProviderConfig, debug bool) (*S3, error) {
	if err := checkConfig(config); err != nil {
		return nil, err
	}
	endpoint := config["endpoint"]
	region := config["region"]
	sessionConfig := aws.Config{
		Endpoint: &endpoint,
		Region:   &region,
	}
	if debug {
		sessionConfig.LogLevel = aws.LogLevel(aws.LogDebugWithRequestErrors |
			aws.LogDebugWithHTTPBody)
		sessionConfig.Logger = aws.NewDefaultLogger()
	}

	session, err := session.NewSession(&sessionConfig)
	if err != nil {
		return nil, errors.New("cannot create s3 session")
	}

	creds := credentials.NewStaticCredentials(config["ak"], config["sk"], "")
	awsConfig := aws.NewConfig().WithCredentials(creds)
	client := s3.New(session, awsConfig)

	return &S3{
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
	return Name, cli.NewCommand(Name, "create new S3 profile").
		WithOption(cli.NewOption("endpoint", "S3 endpoint")).
		WithOption(cli.NewOption("region", "region's name")).
		WithOption(cli.NewOption("ak", "access key")).
		WithOption(cli.NewOption("sk", "secret key"))
}

func (provider *S3) Name() string {
	return Name
}

func (provider *S3) Types() []ObjectType {
	return Types()
}

func (provider *S3) AuthTest() error {
	_, err := provider.client.ListBuckets(nil)
	if err != nil {
		return errors.New("unable to list buckets")
	}
	return nil
}

func (provider *S3) ReadObjects(typeName string) []Object {
	switch typeName {
	case typeBucketObject:
		return provider.readBucketObjects()
	case typeBucket:
		return provider.readBuckets()
	}
	return []Object{}
}

func (provider *S3) DeleteObjects(typeName string, objects []Object) {
	switch typeName {
	case typeBucketObject:
		provider.deleteBucketObjects(objects)
	case typeBucket:
		provider.deleteBuckets(objects)
	}
}

func (provider *S3) StringObject(object string, typeName string) string {
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
	content := strings.Split(*encodedObject, ":")
	if len(content) != 2 {
		return "", "", errors.New("cannot decode bucket object")
	}

	binBucket, err := base64.StdEncoding.DecodeString(content[0])
	if err != nil {
		return "", "", err
	}
	bucketName := string(binBucket)

	binkey, err := base64.StdEncoding.DecodeString(content[1])
	if err != nil {
		return "", "", err
	}
	key := string(binkey)
	return bucketName, key, nil
}

func (provider *S3) readBucketObjects() []Object {
	objects := make([]Object, 0)
	result, err := provider.client.ListBuckets(nil)
	if err != nil {
		return objects
	}
	for _, bucket := range result.Buckets {
		result, err := provider.client.ListObjects(&s3.ListObjectsInput{
			Bucket: bucket.Name,
		})
		if err != nil {
			continue
		}
		for _, object := range result.Contents {
			objects = append(objects, encodeBucketObject(bucket.Name, object.Key))
		}
	}
	return objects
}

func (provider *S3) deleteBucketObjects(bucketObjects []Object) {
	for _, encodedBucketObject := range bucketObjects {
		log.Printf(
			"Deleting object: %s ... ",
			provider.StringObject(encodedBucketObject, typeBucketObject),
		)
		bucketName, key, err := decodeBucketobject(&encodedBucketObject)
		if err != nil {
			log.Println("Error while reading object details: ", err.Error())
		}
		_, err = provider.client.DeleteObject(&s3.DeleteObjectInput{
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

func (provider *S3) readBuckets() []Object {
	buckets := make([]Object, 0)
	result, err := provider.client.ListBuckets(nil)
	if err != nil {
		return buckets
	}
	for _, b := range result.Buckets {
		buckets = append(buckets, encodeBucket(b.Name))
	}
	return buckets
}

func (provider *S3) deleteBuckets(buckets []Object) {
	for _, b64Bucket := range buckets {
		BucketName, err := decodeBucket(&b64Bucket)
		if err != nil {
			continue
		}
		log.Printf("Deleting bucket: %s ... ", BucketName)
		_, err = provider.client.DeleteBucket(&s3.DeleteBucketInput{
			Bucket: &BucketName,
		})
		if err != nil {
			log.Println("Error while deleting bucket: ", err.Error())
		} else {
			log.Println("OK")
		}
	}
}
