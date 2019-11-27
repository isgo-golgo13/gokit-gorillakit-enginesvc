package enginesvcpkg

import (
	"context"
	"time"

	"github.com/go-kit/kit/log"
)

// Middleware describes a service (as opposed to endpoint) middleware.
type Middleware func(Service) Service

func LoggingMiddleware(logger log.Logger) Middleware {
	return func(next Service) Service {
		return &loggingMiddleware{
			next:   next,
			logger: logger,
		}
	}
}

type loggingMiddleware struct {
	next   Service
	logger log.Logger
}

func (mw loggingMiddleware) RegisterEngine(ctx context.Context, e Engine) (err error) {
	defer func(begin time.Time) {
		mw.logger.Log("method", "RegisterEngine", "id", e.ID, "took", time.Since(begin), "err", err)
	}(time.Now())
	return mw.next.RegisterEngine(ctx, e)
}

func (mw loggingMiddleware) GetRegisteredEngine(ctx context.Context, id string) (e Engine, err error) {
	defer func(begin time.Time) {
		mw.logger.Log("method", "GetRegisteredEngine", "id", id, "took", time.Since(begin), "err", err)
	}(time.Now())
	return mw.next.GetRegisteredEngine(ctx, id)
}

