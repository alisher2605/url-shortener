package database

import (
	"context"
	"github.com/alisher2605/url-shortener/config"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"go.uber.org/zap"
	"time"
)

type Datastore interface {
	Connect()
	HealthCheck() error
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
		Region:           d.config.Region,
		Credentials:      aws.NewCredentialsCache(credentials.NewStaticCredentialsProvider(d.config.AccessKeyId, d.config.AccessKeySecret, "")),
		RetryMaxAttempts: d.config.RetryAttempts,
		BaseEndpoint:     &d.config.Endpoint,
	})

	return
}

func (d *Database) HealthCheck() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := d.db.ListTables(ctx, &dynamodb.ListTablesInput{})
	if err != nil {
		zap.S().Errorf("can't connect to the database: %v", err)
		return err
	}

	return nil
}
