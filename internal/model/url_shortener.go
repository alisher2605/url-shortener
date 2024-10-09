package model

import "time"

type UrlRequest struct {
	Url string `json:"url"`
}

type UrlShortener struct {
	Id             int
	UrlHash        string
	LongUrl        string
	CreatedAt      time.Time
	ExpirationTime time.Duration
}

type ShortenedUrlResponse struct {
	Response
	*UrlRequest
}
