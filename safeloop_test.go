package safeloop_test

import (
	"context"
	"fmt"
	"github.com/wirekang/safeloop"
	"testing"
	"time"
)

func makeStep(f func() error) safeloop.Step {
	return func(_ time.Duration) error {
		return f()
	}
}

var emptyStep = makeStep(func() error {
	return nil
})

func TestLoopLimit(t *testing.T) {
	i := 0
	const limit = 3
	f := makeStep(func() error {
		if i == limit {
			t.Fatal("limit")
		}

		i += 1
		return nil
	})

	safeloop.Loop(context.Background(), safeloop.LoopOption{
		Step:  f,
		Limit: limit,

		OnError: func(err error) {
			t.Fatal(err)
		},
		OnFinish: func(err error) {
			if err != nil {
				t.Fatal(err)
			}
		},
	})
}

func TestLoopDelayBefore(t *testing.T) {
	d := time.Millisecond * 500
	start := time.Now()
	safeloop.Loop(context.Background(), safeloop.LoopOption{
		DelayBefore: d,
		Step:        emptyStep,
		Limit:       1,
	})
	if time.Since(start) < d {
		t.Fatal("delayBefore")
	}
}

func TestLoopOnce1(t *testing.T) {
	f := makeStep(func() error {
		return nil
	})

	onFinished := false
	safeloop.Loop(context.Background(), safeloop.LoopOption{
		Step: f,
		OnError: func(err error) {
			t.Fatal(err)
		},
		OnFinish: func(err error) {
			if onFinished {
				t.Fatal("onFinish twice")
			}

			onFinished = true
			if err != nil {
				t.Fatal(err)
			}
		},
		Once: true,
	})
}

func TestLoopOnce2(t *testing.T) {
	rErr := fmt.Errorf("returning error")
	i := 0
	const count = 3
	f := makeStep(func() error {
		if i == count {
			return nil
		}

		i += 1
		return rErr
	})

	onErrorCount := 0

	safeloop.Loop(context.Background(), safeloop.LoopOption{
		Step: f,
		OnError: func(err error) {
			onErrorCount += 1
		},
		OnFinish: func(err error) {
			if err != nil {
				t.Fatal(err)
			}
		},
		Once: true,
	})

	if onErrorCount != count {
		t.Fatalf("count mismatch %d %d", onErrorCount, count)
	}
}

func TestLoopPanic(t *testing.T) {
	var f func() = nil
	i := 0
	safeloop.Loop(context.Background(), safeloop.LoopOption{
		Step: func(_ time.Duration) error {
			f() // panic
			return nil
		},
		OnError: func(err error) {
			i += 1
			if err == nil {
				t.Fatal("no error")
			}
		},
		Limit: 1,
	})
	if i != 1 {
		t.Fatalf("OnError count is %d", i)
	}
}
