package instrumentation

import (
	"context"
	"time"

	"github.com/welaw/go-congress/proto"
)

func (mw instrumentatingMiddleware) SendVote(ctx context.Context, scope *proto.ItemRange) (res *proto.SendVoteReply, err error) {
	defer func(begin time.Time) {
		lvs := []string{"method", "send_vote", "error", ""}
		mw.requestCount.With(lvs...).Add(1)
		mw.requestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())
	res, err = mw.next.SendVote(ctx, scope)
	return
}
