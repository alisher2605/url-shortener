package model

import "errors"

const (
	HostKey = "HOST"
)

type Response struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

var (
	ErrInvalidRequestBody = errors.New("invalid request body")
)
