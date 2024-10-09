package database

import (
	"context"
	"github.com/alisher2605/url-shortener/internal/model"
	"github.com/aws/aws-sdk-go-v2/aws"
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
	UrlByHash(ctx context.Context, url string) (*model.UrlShortener, error)
}

type UrlShortener struct {
	db *dynamodb.Client
}

func NewUrlShortener(db *dynamodb.Client) UrlShortenerRepository {
	return &UrlShortener{
		db: db,
	}
}

func (u *UrlShortener) AddUrl(ctx context.Context, url *model.UrlShortener) error {

	input := &dynamodb.PutItemInput{
		TableName: aws.String(TableName),
		Item: map[string]types.AttributeValue{
			"id":         &types.AttributeValueMemberN{Value: strconv.Itoa(url.Id)},
			"url_hash":   &types.AttributeValueMemberS{Value: url.UrlHash},
			"long_url":   &types.AttributeValueMemberS{Value: url.LongUrl},
			"created_at": &types.AttributeValueMemberS{Value: url.CreatedAt.Format(time.RFC3339)},
			"ttl":        &types.AttributeValueMemberN{Value: strconv.Itoa(int(url.ExpirationTime.Nanoseconds()))}, // TTL attribute in Unix epoch time
		},
	}

	_, err := u.db.PutItem(ctx, input)
	if err != nil {
		return err
	}

	return nil
}

func (u *UrlShortener) UrlByHash(ctx context.Context, url string) (*model.UrlShortener, error) {
	//input := &dynamodb.QueryInput{}

	return nil, nil
}
