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

type UrlShortenerRepository interface {
	AddUrl(url *model.UrlShortener) error
	LongUrl(url string) (*model.UrlShortener, error)
}

type UrlShortener struct {
	db *dynamodb.Client
}

func NewUrlShortener(db *dynamodb.Client) UrlShortenerRepository {
	return &UrlShortener{
		db: db,
	}
}

func (u *UrlShortener) AddUrl(url *model.UrlShortener) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	input := &dynamodb.PutItemInput{
		TableName: aws.String("url-shortener"),
		Item: map[string]types.AttributeValue{
			"id":              &types.AttributeValueMemberN{Value: strconv.Itoa(url.Id)},
			"short_url":       &types.AttributeValueMemberS{Value: url.ShortUrl},
			"long_url":        &types.AttributeValueMemberS{Value: url.LongUrl},
			"created_at":      &types.AttributeValueMemberS{Value: url.CreatedAt.Format(time.RFC3339)},
			"expiration_time": &types.AttributeValueMemberS{Value: url.ExpirationTime.Format(time.RFC3339)},
		},
	}

	_, err := u.db.PutItem(ctx, input)
	if err != nil {
		return err
	}

	return nil
}

func (u *UrlShortener) LongUrl(url string) (*model.UrlShortener, error) {
	return nil, nil
}
