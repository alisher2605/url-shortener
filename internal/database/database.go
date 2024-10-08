package database

import (
	"context"
	"github.com/alisher2605/url-shortener/config"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"time"
)

type Datastore interface {
	Connect() error
	Close() error
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

func (d *Database) Connect() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	defer cancel()

	conf := &aws.Config{
		Region:                      d.config.Region,
		Credentials:                 aws.NewCredentialsCache(credentials.NewStaticCredentialsProvider(d.config.AccessKeyId, d.config.AccessKeySecret, "")),
		RetryMaxAttempts:            5,
		BaseEndpoint:                &d.config.Endpoint,
		RequestMinCompressSizeBytes: 5,
	}

	aws.C
	//cfg, err := awsConfig.N(ctx,
	//	conf,
	//
	//	//awsConfig.WithEndpointResolver(aws.EndpointResolverFunc(
	//	//	func(service, region string) (aws.Endpoint, error) {
	//	//		if service == dynamodb.ServiceID && region == d.config.Region {
	//	//			return aws.Endpoint{
	//	//				URL: d.config.Endpoint,
	//	//			}, nil
	//	//		}
	//	//		return aws.Endpoint{}, fmt.Errorf("unknown endpoint requested")
	//	//	},
	//	//)),
	//)
	//if err != nil {
	//	log.Fatalf("unable to load SDK config, %v", err)
	//}
	//
	//client := dynamodb.NewFromConfig(cfg)
	//
	//result, err := client.ListTables(ctx, &dynamodb.ListTablesInput{})
	//if err != nil {
	//	log.Fatalf("failed to list tables, %v", err)
	//}
	//
	//fmt.Println("Tables: ", result.TableNames)

	return nil
}

func (d *Database) Close() error {
	return nil
}
