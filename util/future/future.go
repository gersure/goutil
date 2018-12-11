package future



import (
	"time"

	"golang.org/x/net/context"
)


type Value interface{}

type Future interface {
	// Blocks on Future awaiting result
	Get() (Value, error)
	// Blocks on Future awaiting result, but returns a ErrTimeout if
	// the timeout Duration is hit before result returns.
	// Note that the execution still continues in Future after timeout.
	GetWithTimeout(timeout time.Duration) (Value, error)
	// Blocks on Future awaiting result as well as provided Context.
	// Execution continues in Future, even if Context hits deadline
	// or is canceled.
	GetWithContext(context.Context) (Value, error)
}

// Returned when a Future has timed out
var ErrTimeout = context.DeadlineExceeded

// Returned when a Future has be canceled
var ErrCanceled = context.Canceled

// Creates a new Future. Func is asynchronously called and it is
// resolved with a Get or GetWithTimeout call on the Future.
func NewFuture(Func func() (Value, error)) Future {
	return newFutureResult(Func)
}

type futureResult struct {
	result chan *result
}

type result struct {
	value Value
	err   error
}

func newFutureResult(Func func() (Value, error)) *futureResult {
	f := &futureResult{
		result: make(chan *result),
	}
	go func() {
		defer close(f.result)
		value, err := Func()
		f.result <- &result{value, err}
	}()
	return f
}

func (f *futureResult) Get() (Value, error) {
	return f.GetWithContext(context.Background())
}

func (f *futureResult) GetWithTimeout(timeout time.Duration) (Value, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	return f.GetWithContext(ctx)
}

func (f *futureResult) GetWithContext(ctx context.Context) (Value, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case result := <-f.result:
		return result.value, result.err
	}
}