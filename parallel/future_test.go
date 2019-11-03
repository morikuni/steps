package parallel

import (
	"context"
	"testing"
	"time"

	"github.com/morikuni/steps"
	"github.com/morikuni/steps/internal/assert"
)

func TestFuture(t *testing.T) {
	ctx := context.Background()

	t.Run("wait timeout", func(t *testing.T) {
		f := NewFuture()

		ctx, cancel := context.WithTimeout(ctx, time.Millisecond)
		defer cancel()

		r, err := f.Wait(ctx)
		assert.Equal(t, steps.Fail, r)
		assert.Equal(t, context.DeadlineExceeded, err)
	})

	t.Run("wait success", func(t *testing.T) {
		f := NewFuture()
		f.Report(steps.Success, nil)

		ctx, cancel := context.WithTimeout(ctx, time.Millisecond)
		defer cancel()

		r, err := f.Wait(ctx)
		assert.Equal(t, steps.Success, r)
		assert.Nil(t, err)
	})

	t.Run("on called", func(t *testing.T) {
		f := NewFuture()

		var (
			done  = make(chan struct{})
			onR   steps.Result
			onErr error
		)
		f.On(steps.Success, func(r steps.Result, err error) {
			onR = r
			onErr = err
			close(done)
		})

		toCtx, cancel := context.WithTimeout(ctx, time.Millisecond)
		defer cancel()
		select {
		case <-done:
			t.Fatal("unexpected done")
		case <-toCtx.Done():
		}

		f.Report(steps.Success, nil)

		toCtx, cancel = context.WithTimeout(ctx, time.Millisecond)
		defer cancel()
		select {
		case <-done:
		case <-toCtx.Done():
			t.Fatal("unexpected done")
		}

		assert.Equal(t, steps.Success, onR)
		assert.Equal(t, nil, onErr)
	})

	t.Run("match", func(t *testing.T) {
		f := NewFuture()

		assert.Equal(t, false, f.Match(steps.Success))
		assert.Equal(t, false, f.Match(steps.Fail))

		f.Report(steps.Success, nil)

		assert.Equal(t, true, f.Match(steps.Success))
	})
}
