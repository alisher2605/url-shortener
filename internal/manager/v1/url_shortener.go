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
	LongUrl(ctx context.Context, hash string) (string, error)
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
		UrlHash:        hash,
		LongUrl:        url.Url,
		CreatedAt:      time.Now(),
		ExpirationTime: time.Now().Add(time.Duration(u.urlTtl) * time.Hour),
	})
	if err != nil {
		zap.S().Errorf("[ERROR] couldn't add url to database -  %v", err)

		return nil, err
	}

	return &model.UrlRequest{Url: fmt.Sprintf("%s/api/v1/%s", os.Getenv(model.HostKey), hash)}, nil
}

func (u *urlShortenerManager) LongUrl(ctx context.Context, hash string) (string, error) {
	url, err := u.repository.UrlByHash(ctx, hash)
	if err != nil {
		zap.S().Errorf("[ERROR] couldn't get url from database -  %v", err)

		return "", err
	}

	return url.LongUrl, nil
}
