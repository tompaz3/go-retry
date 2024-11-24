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
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/tompaz3/go-retry"
)

func Test_Policy_FixedDelay_Interval_WhenDefault(t *testing.T) {
	t.Parallel()
	p := retry.Policy().
		FixedDelay().
		Build()
	got := p.Interval()
	assert.Equal(t, time.Second, got)
}

func Test_Policy_FixedDelay_Interval_WhenLessThanZero(t *testing.T) {
	t.Parallel()
	p := retry.Policy().
		FixedDelay().
		WithInterval(time.Duration(-5)).
		Build()
	got := p.Interval()
	assert.Equal(t, time.Second, got)
}

func Test_Policy_FixedDelay_Interval_WhenZero(t *testing.T) {
	t.Parallel()
	p := retry.Policy().
		FixedDelay().
		WithInterval(time.Duration(0)).
		Build()
	got := p.Interval()
	assert.Equal(t, time.Second, got)
}

func Test_Policy_FixedDelay_Interval_WhenGreaterThanZero(t *testing.T) {
	t.Parallel()
	p := retry.Policy().
		FixedDelay().
		WithInterval(5 * time.Second).
		Build()
	got := p.Interval()
	assert.Equal(t, 5*time.Second, got)
}

func Test_Policy_FixedDelay_MaxAttempts_WhenDefault(t *testing.T) {
	t.Parallel()
	p := retry.Policy().
		FixedDelay().
		Build()
	got := p.MaxAttempts()
	assert.Equal(t, int64(3), got)
}

func Test_Policy_FixedDelay_MaxAttempts_WhenLessThanZero(t *testing.T) {
	t.Parallel()
	p := retry.Policy().
		FixedDelay().
		WithMaxAttempts(int64(-5)).
		Build()
	got := p.MaxAttempts()
	assert.Equal(t, int64(-1), got)
}

func Test_Policy_FixedDelay_MaxAttempts_WhenZero(t *testing.T) {
	t.Parallel()
	p := retry.Policy().
		FixedDelay().
		WithMaxAttempts(int64(0)).
		Build()
	got := p.MaxAttempts()
	assert.Equal(t, int64(3), got)
}

func Test_Policy_FixedDelay_MaxAttempts_WhenGreaterThanZero(t *testing.T) {
	t.Parallel()
	p := retry.Policy().
		FixedDelay().
		WithMaxAttempts(int64(5)).
		Build()
	got := p.MaxAttempts()
	assert.Equal(t, int64(5), got)
}
