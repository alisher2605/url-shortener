package model

import "time"

type ShortenUrlRequest struct {
	Url string `json:"url"`
}

type UrlShortener struct {
	Id             int
	ShortUrl       string
	LongUrl        string
	CreatedAt      time.Time
	ExpirationTime time.Time
}
