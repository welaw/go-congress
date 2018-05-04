package logging

import (
	"github.com/go-kit/kit/log"
	"github.com/welaw/go-congress/services"
)

type loggingMiddleware struct {
	logger log.Logger
	next   services.Service
}

func NewLoggingMiddleware(l log.Logger, s services.Service) services.Service {
	return &loggingMiddleware{
		logger: l,
		next:   s,
	}
}
