package v3

import (
	"context"
	"math/rand"
	"time"
)

// 限速器
// maxDelay表示最大时延，单位为毫秒
func Limiter(ctx context.Context, maxDelay int, limitFunc func()) {
	minDelay := 50
	if maxDelay <= minDelay {
		panic("maxDelay too small, it should be larger than 50")
	}
	for {
		rand.Seed(time.Now().UnixNano())
		delay := minDelay + rand.Intn(maxDelay-minDelay)
		select {
		case <-time.After(time.Millisecond * time.Duration(delay)):
			limitFunc()
		case <-ctx.Done():
			return
		}
	}
}
