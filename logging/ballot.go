package logging

import (
	"context"
	"fmt"
	"time"

	"github.com/welaw/go-congress/proto"
)

func (mw loggingMiddleware) SendVote(ctx context.Context, scope *proto.ItemRange) (resp *proto.SendVoteReply, err error) {
	defer func(begin time.Time) {
		mw.logger.Log(
			"method", "send_vote",
			"item_range", fmt.Sprintf("%+v", scope),
			"err", err,
			"took", time.Since(begin),
		)
	}(time.Now())
	resp, err = mw.next.SendVote(ctx, scope)
	return
}
