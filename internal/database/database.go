package database

import (
	"context"
	"errors"
	"github.com/alisher2605/url-shortener/config"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"go.uber.org/zap"
	"time"
)

type Datastore interface {
	Connect()
	HealthCheck() error
	SetupDatabase() error
	UrlShortenerRepository() UrlShortenerRepository
}

type Database struct {
	urlShortener UrlShortenerRepository
	db           *dynamodb.Client
	config       config.Database
}

func (d *Database) UrlShortenerRepository() UrlShortenerRepository {
	if d.urlShortener == nil {
		d.urlShortener = NewUrlShortener(d.db)
	}

	return d.urlShortener
}

func NewDatabase(config config.Database) Datastore {
	return &Database{
		config: config,
	}
}

func (d *Database) SetupDatabase() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err := d.createTable(ctx)
	if err != nil {
		return err
	}

	return d.prepareTtl(ctx)
}

func (d *Database) createTable(ctx context.Context) error {
	exist, err := d.checkIfTableExists(ctx)
	if err != nil {
		return err
	}

	if exist {
		return nil
	}

	zap.S().Infof("Table %s doesn't exist. Startiong table creation process at %s", TableName, time.Now())

	_, err = d.db.CreateTable(ctx, &dynamodb.CreateTableInput{
		AttributeDefinitions: []types.AttributeDefinition{
			{
				AttributeName: aws.String("id"),
				AttributeType: types.ScalarAttributeTypeN,
			},
		},
		KeySchema: []types.KeySchemaElement{
			{
				AttributeName: aws.String("id"),
				KeyType:       types.KeyTypeHash,
			},
		},
		TableName: aws.String(TableName),
		ProvisionedThroughput: &types.ProvisionedThroughput{
			ReadCapacityUnits:  aws.Int64(10),
			WriteCapacityUnits: aws.Int64(10),
		},
	})
	if err != nil {
		zap.S().Errorf("failed to create table, %v", err)
		return err
	}

	waiter := dynamodb.NewTableExistsWaiter(d.db)

	err = waiter.Wait(ctx, &dynamodb.DescribeTableInput{TableName: aws.String(TableName)}, 10*time.Second)
	if err != nil {
		zap.S().Errorf("failed to create table, %v", err)
		return err
	}

	zap.S().Infof("Successfully created table %s at %s", TableName, time.Now())

	return nil
}

func (d *Database) checkIfTableExists(ctx context.Context) (bool, error) {
	description, err := d.db.DescribeTable(ctx, &dynamodb.DescribeTableInput{
		TableName: aws.String(TableName),
	})
	if err == nil {
		zap.S().Infof("Table with name %s exists with status %s", TableName, description.Table.TableStatus)

		return true, nil
	}

	var notFoundEx *types.ResourceNotFoundException

	if !errors.As(err, &notFoundEx) {
		zap.S().Errorf("failed to describe table, %v", err)

		return false, err
	}

	return false, nil
}

func (d *Database) prepareTtl(ctx context.Context) error {
	result, err := d.db.DescribeTimeToLive(ctx, &dynamodb.DescribeTimeToLiveInput{TableName: aws.String(TableName)})
	if err != nil {
		zap.S().Errorf("failed to describe TTL, %v", err)
		return err
	}
	if result.TimeToLiveDescription.TimeToLiveStatus == StatusEnabled {
		return nil
	}

	_, err = d.db.UpdateTimeToLive(ctx, &dynamodb.UpdateTimeToLiveInput{
		TableName: aws.String(TableName),
		TimeToLiveSpecification: &types.TimeToLiveSpecification{
			AttributeName: aws.String("ttl"),
			Enabled:       aws.Bool(true),
		},
	})
	if err != nil {
		zap.S().Errorf("couldn't enable ttl: %v", err)
		return err
	}

	return nil
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
