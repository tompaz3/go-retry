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

import "time"

const (
	defaultInitialInterval       = time.Second
	defaultMaxInterval           = 30 * time.Second
	defaultMaxAttempts           = int64(3)
	defaultBackOffCoefficient    = float64(2.0)
	fixedDelayBackOffCoefficient = float64(1.0)
	unlimitedMaxInterval         = time.Duration(-1)
	undefinedMaxAttempts         = int64(-1)
)

// policy - normalized retry policy, could represent both BackOffPolicy and FixedDelayPolicy.
type policy interface {
	getInitialInterval() time.Duration
	getMaxInterval() time.Duration
	getMaxAttempts() int64
	getBackOffCoefficient() float64
}

// BackOffPolicy represents the exponential backoff policy for retrying.
type BackOffPolicy struct {
	initialInterval    time.Duration
	maxInterval        time.Duration
	maxAttempts        int64
	backOffCoefficient float64
}

// InitialInterval returns the initial interval between retries.
func (p BackOffPolicy) InitialInterval() time.Duration {
	return p.initialInterval
}

// MaxInterval returns the maximum interval between retries.
func (p BackOffPolicy) MaxInterval() time.Duration {
	return p.maxInterval
}

// MaxAttempts returns the maximum number of attempts.
func (p BackOffPolicy) MaxAttempts() int64 {
	return p.maxAttempts
}

// BackOffCoefficient returns the backoff delay coefficient.
func (p BackOffPolicy) BackOffCoefficient() float64 {
	return p.backOffCoefficient
}

// HasUnlimitedMaxInterval returns true if the policy has an unlimited max interval.
func (p BackOffPolicy) HasUnlimitedMaxInterval() bool {
	return p.maxInterval == unlimitedMaxInterval
}

// IsAttemptingIndefinitely returns true if the policy is attempting indefinitely (no max attempts limit).
func (p BackOffPolicy) IsAttemptingIndefinitely() bool {
	return p.maxAttempts == undefinedMaxAttempts
}

func (p BackOffPolicy) getInitialInterval() time.Duration {
	return p.initialInterval
}

func (p BackOffPolicy) getMaxInterval() time.Duration {
	return p.maxInterval
}

func (p BackOffPolicy) getMaxAttempts() int64 {
	return p.maxAttempts
}

func (p BackOffPolicy) getBackOffCoefficient() float64 {
	return p.backOffCoefficient
}

// FixedDelayPolicy represents the fixed delay policy for retrying.
type FixedDelayPolicy struct {
	interval    time.Duration
	maxAttempts int64
}

// Interval returns the interval between retries.
func (p FixedDelayPolicy) Interval() time.Duration {
	return p.interval
}

// MaxAttempts returns the maximum number of attempts.
func (p FixedDelayPolicy) MaxAttempts() int64 {
	return p.maxAttempts
}

// IsAttemptingIndefinitely returns true if the policy is attempting indefinitely (no max attempts limit).
func (p FixedDelayPolicy) IsAttemptingIndefinitely() bool {
	return p.maxAttempts == undefinedMaxAttempts
}

func (p FixedDelayPolicy) getInitialInterval() time.Duration {
	return p.interval
}

func (p FixedDelayPolicy) getMaxInterval() time.Duration {
	return p.interval
}

func (p FixedDelayPolicy) getMaxAttempts() int64 {
	return p.maxAttempts
}

func (p FixedDelayPolicy) getBackOffCoefficient() float64 {
	return fixedDelayBackOffCoefficient
}
