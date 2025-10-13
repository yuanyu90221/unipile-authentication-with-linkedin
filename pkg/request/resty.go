package request

import (
	"resty.dev/v3"
)

type RequestHandler struct {
	Client *resty.Client
}

func NewRequestHandler() *RequestHandler {
	return &RequestHandler{
		Client: resty.New(),
	}
}
