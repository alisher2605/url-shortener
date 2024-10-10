package model

import "time"

type UrlRequest struct {
	Url string `json:"url"`
}

type UrlShortener struct {
	UrlHash        string    `dynamodbav:"id"`
	LongUrl        string    `dynamodbav:"long_url"`
	CreatedAt      time.Time `dynamodbav:"id"`
	ExpirationTime time.Time `dynamodbav:"-"`
}

type ShortenedUrlResponse struct {
	Response
	*UrlRequest
}
