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

func Test_Policy_BackOff_InitialInterval_WhenDefault(t *testing.T) {
	t.Parallel()
	p := retry.Policy().
		BackOff().
		Build()
	got := p.InitialInterval()
	assert.Equal(t, time.Second, got)
}

func Test_Policy_BackOff_InitialInterval_WhenLessThanZero(t *testing.T) {
	t.Parallel()
	p := retry.Policy().
		BackOff().
		WithInitialInterval(time.Duration(-5)).
		Build()
	got := p.InitialInterval()
	assert.Equal(t, time.Second, got)
}

func Test_Policy_BackOff_InitialInterval_WhenZero(t *testing.T) {
	t.Parallel()
	p := retry.Policy().
		BackOff().
		WithInitialInterval(time.Duration(0)).
		Build()
	got := p.InitialInterval()
	assert.Equal(t, time.Second, got)
}

func Test_Policy_BackOff_InitialInterval_WhenGreaterThanZero(t *testing.T) {
	t.Parallel()
	p := retry.Policy().
		BackOff().
		WithInitialInterval(5 * time.Second).
		Build()
	got := p.InitialInterval()
	assert.Equal(t, 5*time.Second, got)
}

func Test_Policy_BackOff_MaxInterval_WhenDefault(t *testing.T) {
	t.Parallel()
	p := retry.Policy().
		BackOff().
		Build()
	got := p.MaxInterval()
	assert.Equal(t, 30*time.Second, got)
}

func Test_Policy_BackOff_MaxInterval_WhenLessThanZero(t *testing.T) {
	t.Parallel()
	p := retry.Policy().
		BackOff().
		WithMaxInterval(time.Duration(-5)).
		Build()
	got := p.MaxInterval()
	assert.Equal(t, time.Duration(-1), got)
}

func Test_Policy_BackOff_MaxInterval_WhenZero(t *testing.T) {
	t.Parallel()
	p := retry.Policy().
		BackOff().
		WithMaxInterval(time.Duration(0)).
		Build()
	got := p.MaxInterval()
	assert.Equal(t, 30*time.Second, got)
}

func Test_Policy_BackOff_MaxInterval_WhenGreaterThanZero(t *testing.T) {
	t.Parallel()
	p := retry.Policy().
		BackOff().
		WithMaxInterval(5 * time.Second).
		Build()
	got := p.MaxInterval()
	assert.Equal(t, 5*time.Second, got)
}

func Test_Policy_BackOff_MaxAttempts_WhenDefault(t *testing.T) {
	t.Parallel()
	p := retry.Policy().
		BackOff().
		Build()
	got := p.MaxAttempts()
	assert.Equal(t, int64(3), got)
}

func Test_Policy_BackOff_MaxAttempts_WhenLessThanZero(t *testing.T) {
	t.Parallel()
	p := retry.Policy().
		BackOff().
		WithMaxAttempts(int64(-5)).
		Build()
	got := p.MaxAttempts()
	assert.Equal(t, int64(-1), got)
}

func Test_Policy_BackOff_MaxAttempts_WhenZero(t *testing.T) {
	t.Parallel()
	p := retry.Policy().
		BackOff().
		WithMaxAttempts(int64(0)).
		Build()
	got := p.MaxAttempts()
	assert.Equal(t, int64(3), got)
}

func Test_Policy_BackOff_MaxAttempts_WhenGreaterThanZero(t *testing.T) {
	t.Parallel()
	p := retry.Policy().
		BackOff().
		WithMaxAttempts(int64(5)).
		Build()
	got := p.MaxAttempts()
	assert.Equal(t, int64(5), got)
}

func Test_Policy_BackOff_BackOffCoefficient_WhenDefault(t *testing.T) {
	t.Parallel()
	p := retry.Policy().
		BackOff().
		Build()
	got := p.BackOffCoefficient()
	assert.Equal(t, 2.0, got)
}

func Test_Policy_BackOff_BackOffCoefficient_WhenLessThanZero(t *testing.T) {
	t.Parallel()
	p := retry.Policy().
		BackOff().
		WithBackOffCoefficient(-1.0).
		Build()
	got := p.BackOffCoefficient()
	assert.Equal(t, 2.0, got)
}

func Test_Policy_BackOff_BackOffCoefficient_WhenZero(t *testing.T) {
	t.Parallel()
	p := retry.Policy().
		BackOff().
		WithBackOffCoefficient(0.0).
		Build()
	got := p.BackOffCoefficient()
	assert.Equal(t, 2.0, got)
}

func Test_Policy_BackOff_BackOffCoefficient_WhenGreaterThanZero(t *testing.T) {
	t.Parallel()
	p := retry.Policy().
		BackOff().
		WithBackOffCoefficient(0.1).
		Build()
	got := p.BackOffCoefficient()
	assert.Equal(t, 0.1, got)
}
