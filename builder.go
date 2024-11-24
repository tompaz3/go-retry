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

type Builder struct{}

func Policy() *Builder {
	return &Builder{}
}

func (b Builder) BackOff() *BackOffPolicyBuilder {
	return &BackOffPolicyBuilder{}
}

func (b Builder) FixedDelay() *FixedDelayPolicyBuilder {
	return &FixedDelayPolicyBuilder{}
}

type BackOffPolicyBuilder struct {
	initialInterval    time.Duration
	maxInterval        time.Duration
	maxAttempts        int64
	backOffCoefficient float64
}

func (b BackOffPolicyBuilder) WithInitialInterval(initialInterval time.Duration) BackOffPolicyBuilder {
	return BackOffPolicyBuilder{
		initialInterval:    initialInterval,
		maxInterval:        b.maxInterval,
		maxAttempts:        b.maxAttempts,
		backOffCoefficient: b.backOffCoefficient,
	}
}

func (b BackOffPolicyBuilder) WithMaxInterval(maxInterval time.Duration) BackOffPolicyBuilder {
	return BackOffPolicyBuilder{
		initialInterval:    b.initialInterval,
		maxInterval:        maxInterval,
		maxAttempts:        b.maxAttempts,
		backOffCoefficient: b.backOffCoefficient,
	}
}

func (b BackOffPolicyBuilder) WithMaxIntervalUnlimited() BackOffPolicyBuilder {
	return BackOffPolicyBuilder{
		initialInterval:    b.initialInterval,
		maxInterval:        unlimitedMaxInterval,
		maxAttempts:        b.maxAttempts,
		backOffCoefficient: b.backOffCoefficient,
	}
}

func (b BackOffPolicyBuilder) WithMaxAttempts(maxAttempts int64) BackOffPolicyBuilder {
	return BackOffPolicyBuilder{
		initialInterval:    b.initialInterval,
		maxInterval:        b.maxInterval,
		maxAttempts:        maxAttempts,
		backOffCoefficient: b.backOffCoefficient,
	}
}

func (b BackOffPolicyBuilder) WithMaxAttemptsIndefinite() BackOffPolicyBuilder {
	return BackOffPolicyBuilder{
		initialInterval:    b.initialInterval,
		maxInterval:        b.maxInterval,
		maxAttempts:        undefinedMaxAttempts,
		backOffCoefficient: b.backOffCoefficient,
	}
}

func (b BackOffPolicyBuilder) WithBackOffCoefficient(backOffCoefficient float64) BackOffPolicyBuilder {
	return BackOffPolicyBuilder{
		initialInterval:    b.initialInterval,
		maxInterval:        b.maxInterval,
		maxAttempts:        b.maxAttempts,
		backOffCoefficient: backOffCoefficient,
	}
}

func (b BackOffPolicyBuilder) Build() BackOffPolicy {
	return BackOffPolicy{
		initialInterval:    b.resolveInitialInterval(),
		maxInterval:        b.resolveMaxInterval(),
		maxAttempts:        b.resolveMaxAttempts(),
		backOffCoefficient: b.resolveBackOffCoefficient(),
	}
}

func (b BackOffPolicyBuilder) resolveInitialInterval() time.Duration {
	if b.initialInterval <= 0 {
		return defaultInitialInterval
	}
	return b.initialInterval
}

func (b BackOffPolicyBuilder) resolveMaxInterval() time.Duration {
	if b.maxInterval < 0 {
		return unlimitedMaxInterval
	}
	if b.maxInterval == 0 {
		return defaultMaxInterval
	}
	return b.maxInterval
}

func (b BackOffPolicyBuilder) resolveMaxAttempts() int64 {
	if b.maxAttempts < 0 {
		return undefinedMaxAttempts
	}
	if b.maxAttempts == 0 {
		return defaultMaxAttempts
	}
	return b.maxAttempts
}

func (b BackOffPolicyBuilder) resolveBackOffCoefficient() float64 {
	if b.backOffCoefficient <= 0 {
		return defaultBackOffCoefficient
	}
	return b.backOffCoefficient
}

type FixedDelayPolicyBuilder struct {
	interval    time.Duration
	maxAttempts int64
}

func (b FixedDelayPolicyBuilder) WithInterval(interval time.Duration) FixedDelayPolicyBuilder {
	return FixedDelayPolicyBuilder{
		interval:    interval,
		maxAttempts: b.maxAttempts,
	}
}

func (b FixedDelayPolicyBuilder) WithMaxAttempts(maxAttempts int64) FixedDelayPolicyBuilder {
	return FixedDelayPolicyBuilder{
		interval:    b.interval,
		maxAttempts: maxAttempts,
	}
}

func (b FixedDelayPolicyBuilder) WithMaxAttemptsIndefinite() FixedDelayPolicyBuilder {
	return FixedDelayPolicyBuilder{
		interval:    b.interval,
		maxAttempts: undefinedMaxAttempts,
	}
}

func (b FixedDelayPolicyBuilder) Build() FixedDelayPolicy {
	return FixedDelayPolicy{
		interval:    b.resolveInterval(),
		maxAttempts: b.resolveMaxAttempts(),
	}
}

func (b FixedDelayPolicyBuilder) resolveInterval() time.Duration {
	if b.interval <= 0 {
		return defaultInitialInterval
	}
	return b.interval
}

func (b FixedDelayPolicyBuilder) resolveMaxAttempts() int64 {
	if b.maxAttempts < 0 {
		return undefinedMaxAttempts
	}
	if b.maxAttempts == 0 {
		return defaultMaxAttempts
	}
	return b.maxAttempts
}
