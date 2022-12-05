package safeloop

import (
	"context"
	"time"
)

type Step func(delta time.Duration) error

type LoopOption struct {
	//DelayBefore is delay before loop start.
	DelayBefore time.Duration

	// DelayBetween is delay between loop step.
	DelayBetween time.Duration

	// Step must not nil.
	Step Step

	// OnError called when Step returns error.
	OnError func(error)

	// OnFinish called when the loop finished. The loop could be finished by Limit or Once.
	OnFinish func(error)

	// Limit is max loop count regardless of success. Default is 0 which means infinity.
	Limit uint64

	// If Once is true, the loop finished when Step returns nil.
	Once bool
}

func Loop(ctx context.Context, opt LoopOption) {
	if opt.Step == nil {
		panic("Step is nil")
	}

	stepWrapper := makeWrapper(opt.Step)
	time.Sleep(opt.DelayBefore)
	startTime := time.Now()
	var i uint64 = 1
	for ; ; i += 1 {
		err := ctx.Err()
		if err != nil {
			call(opt.OnError, err)
			call(opt.OnFinish, err)
			return
		}

		err = stepWrapper(time.Since(startTime))
		if err != nil {
			call(opt.OnError, err)
		}

		if (opt.Once && err == nil) || (opt.Limit > 0 && i == opt.Limit) {
			call(opt.OnFinish, err)
			return
		}

		time.Sleep(opt.DelayBetween)
		startTime = time.Now()
	}
}
