package v1

import (
	"context"
	"fmt"
	"github.com/alisher2605/url-shortener/internal/database"
	"github.com/alisher2605/url-shortener/internal/model"
	"github.com/alisher2605/url-shortener/util/snowflake"
	"go.uber.org/zap"
	"os"
	"time"
)

type urlShortenerManager struct {
	urlTtl     int
	repository database.UrlShortenerRepository
}

type UrlShortener interface {
	AddUrl(ctx context.Context, url *model.UrlRequest) (*model.UrlRequest, error)
	LongUrl(ctx context.Context, hash string) string
}

func NewUrlShortener(repository database.UrlShortenerRepository, urlTtl int) UrlShortener {
	return &urlShortenerManager{
		repository: repository,
		urlTtl:     urlTtl,
	}
}

func (u *urlShortenerManager) AddUrl(ctx context.Context, url *model.UrlRequest) (*model.UrlRequest, error) {
	id, err := snowflake.GenerateSnowflakeId()
	if err != nil {
		zap.S().Errorf("[ERROR] couldn't generate snowflake id -  %v", err)

		return nil, err
	}

	hash := snowflake.Base62Conversion(id)

	err = u.repository.AddUrl(ctx, &model.UrlShortener{
		Id:             id,
		UrlHash:        hash,
		LongUrl:        url.Url,
		CreatedAt:      time.Now(),
		ExpirationTime: 10 * time.Second,
	})
	if err != nil {
		zap.S().Errorf("[ERROR] couldn't add url to database -  %v", err)

		return nil, err
	}

	return &model.UrlRequest{Url: fmt.Sprintf("%s/%s", os.Getenv(model.HostKey), hash)}, nil
}

func (u *urlShortenerManager) LongUrl(ctx context.Context, hash string) string {
	//TODO implement me
	panic("implement me")
}
