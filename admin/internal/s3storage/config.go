package s3storage

import "github.com/aws/aws-sdk-go/aws"

// Config S3 client config
type Config struct {
	AccessID         string `yaml:"accessId"`
	AccessKey        string `yaml:"accessKey"`
	Region           string `yaml:"region"`
	Bucket           string `yaml:"bucket"`
	SessionToken     string `yaml:"sessionToken"`
	ACL              string `yaml:"acl"`
	Endpoint         string `yaml:"endpoint"`
	S3Endpoint       string `yaml:"s3Endpoint"`
	S3ForcePathStyle bool   `yaml:"s3ForcePathStyle"`
	CacheControl     string `yaml:"cacheControl"`

	AWSConfig *aws.Config

	RoleARN string `yaml:"roleArn"`
}
