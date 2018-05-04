package logging

import (
	"context"
	"fmt"
	"time"

	"github.com/welaw/go-congress/proto"
)

func (mw loggingMiddleware) SendLaw(ctx context.Context, scope *proto.ItemRange) (resp *proto.SendLawReply, err error) {
	defer func(begin time.Time) {
		mw.logger.Log(
			"method", "send_law",
			"item_range", fmt.Sprintf("%+v", scope),
			"err", err,
			"took", time.Since(begin),
		)
	}(time.Now())
	resp, err = mw.next.SendLaw(ctx, scope)
	return
}

func (mw loggingMiddleware) Status(ctx context.Context, scope *proto.ItemRange) (resp *proto.StatusReply, err error) {
	defer func(begin time.Time) {
		mw.logger.Log(
			"method", "status",
			"scope", fmt.Sprintf("%+v", scope),
			"err", nil,
			"took", time.Since(begin),
		)
	}(time.Now())
	resp, err = mw.next.Status(ctx, scope)
	return
}
