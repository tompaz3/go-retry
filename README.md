go-retry
==============

`go-retry` is a Go package that provides a simple and easy-to-use API for retrying operations in Go.

# Installation

To add `go-retry` to your project add `github.com/tompaz3/go-retry` package to your project.

## Usage

`go-retry` provides a simple API to retry operations in Go. Package supports 2 kinds of retry policies - [FixedDelay](policy.go#L100) and [BackOffPolicy](policy.go#L46).

### Policies

Use `retry.Policy()` function to build retry policy.

Package provides 2 basic functions to retry operations - `retry.Run() error` and `retry.Supply[T any]() (T, error)`.

#### FixedDelay

`FixedDelay` policy will retry the operation with a fixed delay between each retry.

`FixedDelay` policy may be configured with the following options:

* `WithInterval(time.Duration)` - sets the interval between retries.
* `WithMaxAttempts(int64)` - sets the maximum number of retries.
* `WithMaxAttemptsIndefinite()` - sets the maximum number of retries to unlimited.

Additionally, retry functions accept `context.Context` and support context cancellation.

```go
package example

import (
  "time"

  "github.com/tompaz3/go-retry"
)

var retryWithDefaultDelayDefaultNumberOfTimes = retry.Policy().
  FixedDelay().
  Build()

var retryEverySecondMax3Times = retry.Policy().
  FixedDelay().
  WithInterval(time.Second).
  WithMaxAttempts(int64(3)).
  Build()

var retryWithDefaultDelayUnlimitedNumberOfTimes = retry.Policy().
  FixedDelay().
  WithMaxAttemptsIndefinite().
  Build()
```

#### BackOffPolicy

`BackOffPolicy` policy will retry the operation with an increasing delay between each retry.

`BackOffPolicy` policy may be configured with the following options:

* `WithInitialInterval(time.Duration)` - sets the initial interval between retries.
* `WithMaxInterval(time.Duration)` - sets the maximum interval between retries.
* `WithMaxIntervalUnlimited()` - sets the maximum interval between retries to unlimited.
* `WithMaxAttempts(int64)` - sets the maximum number of retries.
* `WithMaxAttemptsIndefinite()` - sets the maximum number of retries to unlimited.
* `WithCoefficient(float64)` - sets the coefficient for the backoff calculation.

Additionally, retry functions accept `context.Context` and support context cancellation.

```go
package example

import (
  "time"

  "github.com/tompaz3/go-retry"
)

var retryWithAllDefaults = retry.Policy().
  BackOff().
  Build()

// retry with initial interval 1 second, max interval 1 minute, max attempts 20 and coefficient 1.5
var customizedRetry = retry.Policy().
  BackOff().
  WithInitialInterval(time.Second).
  WithMaxInterval(time.Minute).
  WithMaxAttempts(int64(20)).
  WithCoefficient(1.5).
  Build()

var retryUnlimitedNumberOfTimes = retry.Policy().
  BackOff().
  WithMaxAttemptsIndefinite().
  Build()

var retryWithUnlimitedMaxInterval = retry.Policy().
  BackOff().
  WithMaxIntervalUnlimited().
  Build()
```

### Retry functions

Operations will be retried until the operation returns no error or the maximum number of retries is reached or the context is canceled.

In case context is canceled, the operation will return [retry.DeadlineExceededError](retry.go#L110) error.

Use one of the 2 functions to trigger retry:

1. `retry.Run(ctx context.Context, slp Sleeper, run RunFunc, p policy) error` - to retry operation that returns error only.
2. `retry.Supply[T any](ctx context.Context, slp Sleeper, supply SupplyFunc[T], p policy) (T, error)` - to retry operation that returns both value and error.

NOTE: `Sleeper` is an interface which provides _sleep_ logic. User must provide their own `Sleeper` implementation to invoke retry functions. See [Sleeper](#sleeper) section for more details.



```go
package example

import (
  "context"
  "time"

  "github.com/tompaz3/go-retry"
)

type EventPublisher interface {
  Publish(ctx context.Context, event Event) error
}

func PublishEventRetry(ctx context.Context, publisher EventPublisher, event Event) error {
  // fixed delay policy - retry at most 3 times with 200ms interval
  policy := retry.Policy().
    FixedDelay().
    WithInterval(200 * time.Millisecond).
    WithMaxAttempts(int64(3)).
    Build()

  systemTimeSleeper := retry.SleeperF(time.Sleep)

  retryFunc := func() error {
    return publisher.Publish(ctx, event)
  }

  return retry.Run(ctx, systemTimeSleeper, retryFunc, policy)
}

type DataRetriever interface {
  Retrieve(ctx context.Context) (Data, error)
}

func RetrieveDataRetry(ctx context.Context, retriever DataRetriever) (Data, error) {
  // exponential back off policy - retry with initial interval 1 second,
  // max interval 1 minute
  // back off coefficient 2.0
  // and unlimited number of retries
  policy := retry.Policy().
    BackOff().
    WithInitialInterval(time.Second).
    WithMaxInterval(time.Minute).
    WithBackOffCoefficient(2.0).
    WithMaxAttemptsIndefinite().
    Build()

  systemTimeSleeper := retry.SleeperF(time.Sleep)

  supplyFunc := func() (Data, error) {
    return retriever.Retrieve(ctx)
  }

  // cancel retry after 30 minutes
  timeoutCtx, cancel := context.WithTimeout(ctx, 30 * time.Minute)
  defer cancel()
  return retry.Supply(timeoutCtx, systemTimeSleeper, supplyFunc, policy)
}
```

#### Sleeper
[Sleeper](retry.go#L35) is an interface that provides _sleep_ logic for retry functions.
User must provide their own `Sleeper` implementation.

For user's convenience `SleeperF` function has been added to create `Sleeper` from a function.

Basic sleeper implementation examples are presented below:

```go
package example

import (
  "time"

  clock "github.com/jonboulle/clockwork"
  "github.com/tompaz3/go-retry"
)

// SystemTimeSleeper - example sleeper implementation using time.Sleep
type SystemTimeSleeper struct {}

func (SystemTimeSleeper) Sleep(d time.Duration) {
  time.Sleep(d)
}

func ExampleSleeper() {
  // system time sleeper implementation
  var systemTimeSleeper retry.Sleeper = SystemTimeSleeper{}

  // simple time.Sleep sleeper implementation using retry.SleeperF
  var systemTimeSleeperFromFunc retry.Sleeper = retry.SleeperF(time.Sleep)

  // fake clock which can be used for testing with time simulated by the user
  fakeClok := clock.NewFakeClockAt(time.Now())
  var fakeClockSleeper retry.Sleeper = retry.SleeperF(fakeClock.Sleep)
}

```

## License

The generator is licensed under the MIT License. License available at [LICENSE](LICENSE).

## Contributing

No contribution policy has been defined yet. It is a tiny, single-contributor project.

The project is considered feature-complete at the moment. Most likely, will be updated for bug fixing and vulnerability patches only.

In case the author cannot maintain the project, a new strategy will be created to keep the project alive.
