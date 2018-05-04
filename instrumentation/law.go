package instrumentation

import (
	"context"
	"time"

	"github.com/welaw/go-congress/proto"
)

func (mw instrumentatingMiddleware) SendLaw(ctx context.Context, scope *proto.ItemRange) (res *proto.SendLawReply, err error) {
	defer func(begin time.Time) {
		lvs := []string{"method", "send_law", "error", ""}
		mw.requestCount.With(lvs...).Add(1)
		mw.requestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())
	res, err = mw.next.SendLaw(ctx, scope)
	return
}

func (mw instrumentatingMiddleware) Status(ctx context.Context, scope *proto.ItemRange) (res *proto.StatusReply, err error) {
	defer func(begin time.Time) {
		lvs := []string{"method", "status", "error", ""}
		mw.requestCount.With(lvs...).Add(1)
		mw.requestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())
	res, err = mw.next.Status(ctx, scope)
	return
}
