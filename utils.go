package safeloop

import (
	"fmt"
	"time"
)

func call[T any](f func(T), v T) {
	if f != nil {
		f(v)
	}
}

func makeWrapper(f Step) Step {
	return func(d time.Duration) (err error) {
		defer func() {
			r := recover()
			if r != nil {
				if err != nil {
					err = fmt.Errorf("recover: %v, %w", r, err)
					return
				}
				err = fmt.Errorf("recover: %v", r)
			}
		}()
		return f(d)
	}
}
