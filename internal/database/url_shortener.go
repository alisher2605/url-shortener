package database

import (
	"context"
	"github.com/alisher2605/url-shortener/internal/model"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"strconv"
	"time"
)

const (
	StatusEnabled = "ENABLED"
	TableName     = "url_shortener"
)

type UrlShortenerRepository interface {
	AddUrl(ctx context.Context, url *model.UrlShortener) error
	UrlByHash(ctx context.Context, hash string) (*model.UrlShortener, error)
}

type UrlShortener struct {
	client *dynamodb.Client
}

func NewUrlShortener(client *dynamodb.Client) UrlShortenerRepository {
	return &UrlShortener{
		client: client,
	}
}

func (u *UrlShortener) AddUrl(ctx context.Context, url *model.UrlShortener) error {
	input := &dynamodb.PutItemInput{
		TableName: aws.String(TableName),
		Item: map[string]types.AttributeValue{
			"id":         &types.AttributeValueMemberS{Value: url.UrlHash},
			"long_url":   &types.AttributeValueMemberS{Value: url.LongUrl},
			"created_at": &types.AttributeValueMemberS{Value: url.CreatedAt.Format(time.RFC3339)},
			"ttl":        &types.AttributeValueMemberN{Value: strconv.Itoa(int(url.ExpirationTime.Unix()))},
		},
	}

	_, err := u.client.PutItem(ctx, input)
	if err != nil {
		return err
	}

	return nil
}

func (u *UrlShortener) UrlByHash(ctx context.Context, hash string) (*model.UrlShortener, error) {
	input := &dynamodb.GetItemInput{
		TableName: aws.String(TableName),
		Key: map[string]types.AttributeValue{
			"id": &types.AttributeValueMemberS{Value: hash},
		},
	}

	result, err := u.client.GetItem(ctx, input)
	if err != nil {
		return nil, err
	}

	if result.Item == nil {
		return nil, nil
	}

	url := new(model.UrlShortener)

	err = attributevalue.UnmarshalMap(result.Item, &url)
	if err != nil {
		return nil, err
	}

	return url, nil
}
