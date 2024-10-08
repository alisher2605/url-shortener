package database

import (
	"github.com/alisher2605/url-shortener/config"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

type Datastore interface {
	Connect()
}

type Database struct {
	db     *dynamodb.Client
	config config.Database
}

func NewDatabase(config config.Database) Datastore {
	return &Database{
		config: config,
	}
}

func (d *Database) Connect() {
	d.db = dynamodb.NewFromConfig(aws.Config{
		Region:                      d.config.Region,
		Credentials:                 aws.NewCredentialsCache(credentials.NewStaticCredentialsProvider(d.config.AccessKeyId, d.config.AccessKeySecret, "")),
		RetryMaxAttempts:            5,
		BaseEndpoint:                &d.config.Endpoint,
		RequestMinCompressSizeBytes: 5,
	})

	return
}

func (d *Database) Close() error {
	return nil
}
