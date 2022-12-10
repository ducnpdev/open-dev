package pipeline

import (
	"context"
	"testing"
)

func BenchmarkPrimeNumbers(b *testing.B) {
	ctx := context.Background()
	pline := pipeline()
	for i := 0; i < b.N; i++ {
		checkattempt(ctx, pline)
	}
}
