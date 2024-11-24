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

package retry_test

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/tompaz3/go-retry"

	clock "github.com/jonboulle/clockwork"
)

func Test_Run_ShouldReturnResultAfterRetries(t *testing.T) {
	t.Parallel()

	var startWG sync.WaitGroup
	startWG.Add(1)
	i := 0
	runner := func() error {
		if i == 0 {
			startWG.Done()
		}
		i++
		if i < 3 {
			return assert.AnError
		}
		return nil
	}

	backOffPolicy := retry.Policy().
		BackOff().
		WithInitialInterval(100 * time.Millisecond).
		WithMaxInterval(time.Second).
		WithBackOffCoefficient(2.0).
		WithMaxAttempts(int64(3)).
		Build()

	clk := clock.NewFakeClockAt(time.Now())
	start := clk.Now()

	go func() {
		startWG.Wait()
		for clk.Since(start) <= 5*time.Second {
			clk.Advance(100 * time.Millisecond)
			time.Sleep(50 * time.Millisecond)
		}
	}()

	err := retry.Run(context.Background(), clk, runner, backOffPolicy)

	assert.NoError(t, err)
}

func Test_Run_ShouldReturnErrorWhenMaxAttemptsReached(t *testing.T) {
	t.Parallel()

	var startWG sync.WaitGroup
	startWG.Add(1)
	i := 0
	runner := func() error {
		if i == 0 {
			startWG.Done()
		}
		i++
		if i < 4 {
			return assert.AnError
		}
		return nil
	}
	backOffPolicy := retry.Policy().
		BackOff().
		WithInitialInterval(100 * time.Millisecond).
		WithMaxInterval(time.Second).
		WithBackOffCoefficient(2.0).
		WithMaxAttempts(int64(3)).
		Build()

	clk := clock.NewFakeClockAt(time.Now())
	start := clk.Now()

	go func() {
		startWG.Wait()
		for clk.Since(start) <= 5*time.Second {
			clk.Advance(100 * time.Millisecond)
			time.Sleep(50 * time.Millisecond)
		}
	}()

	err := retry.Run(context.Background(), clk, runner, backOffPolicy)

	assert.Error(t, err)
	assert.Equal(t, assert.AnError, err)
}

func Test_Run_ShouldReturnErrorWhenContextCanceled(t *testing.T) {
	t.Parallel()

	var startWG sync.WaitGroup
	startWG.Add(1)
	i := 0
	runner := func() error {
		if i == 0 {
			startWG.Done()
		}
		i++
		return assert.AnError
	}

	backOffPolicy := retry.Policy().
		BackOff().
		WithInitialInterval(100 * time.Millisecond).
		WithMaxInterval(time.Second).
		WithBackOffCoefficient(2.0).
		WithMaxAttemptsIndefinite().
		Build()

	clk := clock.NewFakeClockAt(time.Now())
	start := clk.Now()

	go func() {
		startWG.Wait()
		for clk.Since(start) <= 5*time.Second {
			clk.Advance(100 * time.Millisecond)
			time.Sleep(50 * time.Millisecond)
		}
	}()

	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(150*time.Millisecond))
	defer cancel()
	err := retry.Run(ctx, clk, runner, backOffPolicy)

	assert.Error(t, err)
	assert.Equal(t, retry.DeadlineExceededError[any]{
		Result: nil,
		Err:    assert.AnError,
	}, err)
}

func Test_Run_ShouldRespectExponentialBackOffPolicy(t *testing.T) {
	t.Parallel()

	clk := clock.NewFakeClockAt(time.Now())

	var startWG sync.WaitGroup
	startWG.Add(1)
	i := 0
	tryTimes := make([]time.Time, 0)
	runner := func() error {
		tryTimes = append(tryTimes, clk.Now())
		if i == 0 {
			startWG.Done()
		}
		i++
		if i < 5 {
			return assert.AnError
		}
		return nil
	}

	backOffPolicy := retry.Policy().
		BackOff().
		WithInitialInterval(100 * time.Millisecond).
		WithMaxInterval(time.Second).
		WithBackOffCoefficient(2.0).
		WithMaxAttempts(int64(5)).
		Build()

	start := clk.Now()

	go func() {
		startWG.Wait()
		for clk.Since(start) <= 5*time.Second {
			clk.Advance(100 * time.Millisecond)
			time.Sleep(50 * time.Millisecond)
		}
	}()

	err := retry.Run(context.Background(), clk, runner, backOffPolicy)

	assert.NoError(t, err)
	assert.Len(t, tryTimes, 5)
	assert.Equal(t, start, tryTimes[0])
	nextTry := start.Add(100 * time.Millisecond)
	assert.Equal(t, nextTry, tryTimes[1])
	nextTry = nextTry.Add(200 * time.Millisecond)
	assert.Equal(t, nextTry, tryTimes[2])
	nextTry = nextTry.Add(400 * time.Millisecond)
	assert.Equal(t, nextTry, tryTimes[3])
	nextTry = nextTry.Add(800 * time.Millisecond)
	assert.Equal(t, nextTry, tryTimes[4])
}

func Test_Run_ShouldRespectFixedDelayPolicy(t *testing.T) {
	t.Parallel()

	clk := clock.NewFakeClockAt(time.Now())

	var startWG sync.WaitGroup
	startWG.Add(1)
	i := 0
	tryTimes := make([]time.Time, 0)
	runner := func() error {
		tryTimes = append(tryTimes, clk.Now())
		if i == 0 {
			startWG.Done()
		}
		i++
		if i < 5 {
			return assert.AnError
		}
		return nil
	}

	backOffPolicy := retry.Policy().
		FixedDelay().
		WithInterval(100 * time.Millisecond).
		WithMaxAttempts(int64(5)).
		Build()

	start := clk.Now()

	go func() {
		startWG.Wait()
		for clk.Since(start) <= 5*time.Second {
			clk.Advance(100 * time.Millisecond)
			time.Sleep(50 * time.Millisecond)
		}
	}()

	err := retry.Run(context.Background(), clk, runner, backOffPolicy)

	assert.NoError(t, err)
	assert.Len(t, tryTimes, 5)
	assert.Equal(t, start, tryTimes[0])
	nextTry := start.Add(100 * time.Millisecond)
	assert.Equal(t, nextTry, tryTimes[1])
	nextTry = nextTry.Add(100 * time.Millisecond)
	assert.Equal(t, nextTry, tryTimes[2])
	nextTry = nextTry.Add(100 * time.Millisecond)
	assert.Equal(t, nextTry, tryTimes[3])
	nextTry = nextTry.Add(100 * time.Millisecond)
	assert.Equal(t, nextTry, tryTimes[4])
}
