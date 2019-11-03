package parallel

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/morikuni/steps"
)

type seqCounter struct {
	count int
	mu    sync.Mutex
}

func (c *seqCounter) Want(t testing.TB, oneOf ...int) steps.Step {
	t.Helper()

	return steps.StepFunc(func(ctx context.Context) (steps.Result, error) {
		t.Helper()

		t.Log(oneOf, time.Now())

		c.mu.Lock()
		defer c.mu.Unlock()
		c.count++

		for _, n := range oneOf {
			if n == c.count {
				return steps.Success, nil
			}
		}

		t.Fatalf("want %v but %d", oneOf, c.count)

		return steps.Success, nil
	})
}

func TestGroup(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	var (
		g Group
		c seqCounter
	)

	f1 := g.Run(ctx, c.Want(t, 1, 2))
	f2 := g.Run(ctx, c.Want(t, 1, 2))
	f3 := g.Run(ctx, c.Want(t, 3),
		After(f1, steps.Success),
		After(f2, steps.Success),
	)
	f4 := g.Run(ctx, c.Want(t, 4, 5),
		After(f3, steps.Success),
	)
	f5 := g.Run(ctx, c.Want(t, 4, 5),
		After(f3, steps.Success),
	)
	f6 := g.Run(ctx, c.Want(t, 6),
		After(f4, steps.Success),
		After(f5, steps.Fail),
	)

	g.Wait()

	c.Want(t, 6)

	t.Log(f1)
	t.Log(f2)
	t.Log(f3)
	t.Log(f4)
	t.Log(f5)
	t.Log(f6)
}
