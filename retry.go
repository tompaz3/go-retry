// MIT License
//
// Copyright (c) 2024 Tomasz Paździurek
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

package retry

import (
	"context"
	"fmt"
	"time"
)

type (
	RunFunc           func() error
	SupplyFunc[T any] func() (T, error)

	Sleeper interface {
		Sleep(duration time.Duration)
	}
)

type SleeperF func(duration time.Duration)

func (f SleeperF) Sleep(duration time.Duration) {
	f(duration)
}

func Run(ctx context.Context, slp Sleeper, run RunFunc, p policy) error {
	return returnErrOnly(Supply(ctx, slp, runFuncToSupplyFunc(run), p))
}

func Supply[T any](ctx context.Context, slp Sleeper, supply SupplyFunc[T], p policy) (T, error) {
	var res T
	var err error
	nextInterval := p.getInitialInterval()

	for i := int64(0); i == int64(0) || i < p.getMaxAttempts(); {
		select {
		case <-ctx.Done():
			return res, DeadlineExceededError[T]{
				Result: res,
				Err:    err,
			}
		default:
		}

		if res, err = supply(); err == nil {
			return res, nil
		}
		currInterval := nextInterval
		nextInterval = calcNextInterval(nextInterval, p.getMaxInterval(), p.getBackOffCoefficient())
		slp.Sleep(currInterval)
		if p.getMaxAttempts() > 0 {
			i++
		}
	}

	return res, err
}

func calcNextInterval(current, maxInterval time.Duration, backOffCoefficient float64) time.Duration {
	if unlimitedMaxInterval == maxInterval {
		return nextInterval(current, backOffCoefficient)
	}

	if current == maxInterval {
		return current
	}

	next := nextInterval(current, backOffCoefficient)
	if next > maxInterval {
		return maxInterval
	}

	return next
}

func nextInterval(current time.Duration, backOffCoefficient float64) time.Duration {
	return time.Duration(float64(current.Nanoseconds()) * backOffCoefficient)
}

func runFuncToSupplyFunc(run RunFunc) SupplyFunc[any] {
	return func() (any, error) {
		return nil, run()
	}
}

func returnErrOnly[T any](_ T, err error) error {
	return err
}

type DeadlineExceededError[T any] struct {
	Result T
	Err    error
}

func (e DeadlineExceededError[T]) Error() string {
	return fmt.Sprintf("Deadline exceeded %v", e.Err)
}
