package pipeline

import (
	"context"
	"fmt"
	"open-dev/redis/redisPkg"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
)

func Main() {
	ctx := context.Background()
	pline := pipeline()
	excutePipe(ctx, pline)
}

func countExec(r []redis.Cmder) int {
	var (
		countSlice int
	)
	for _, item := range r {
		switch v := item.(type) {
		case *redis.ZSliceCmd:
			for range v.Val() {
				countSlice++
			}
		}
	}
	return countSlice
}

func excutePipe(ctx context.Context, pline redis.Pipeliner) {

	now := time.Now()

	var (
		key     = "zaddkey1"
		ttl int = 3
	)

	timeRemove := time.Now().Add((-1) * time.Second * time.Duration(ttl))
	formatTime := strconv.FormatInt(timeRemove.UnixNano(), 10)
	pline.ZRemRangeByScore(ctx, key, "0", formatTime)

	rcmd := pline.ZAdd(ctx, key, &redis.Z{
		Score:  float64(now.UnixNano()),
		Member: now,
	})
	if err := rcmd.Err(); err != nil {
		panic(err)
	}
	pline.Expire(ctx, key, time.Second*time.Duration(ttl))

	pline.ZRangeWithScores(ctx, key, 0, 100)
	outputResult, err := pline.Exec(ctx)

	if err != nil {
		panic(err)
	}
	count := countExec(outputResult)
	fmt.Println("count:", count)
}

func pipeline() redis.Pipeliner {
	redisclient := redisPkg.InitRedis()
	return redisclient.Pipeline()
}
