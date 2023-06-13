package closer

import (
	"context"
	"fmt"
	"strings"
	"sync"
)

type closer struct {
	mu       sync.Mutex
	handlers []handler
}

func New() *closer {
	return &closer{}
}

func (c *closer) Add(h handler) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.handlers = append(c.handlers, h)
}

func (c *closer) Close(ctx context.Context) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	var (
		errs     = make([]string, 0, len(c.handlers))
		complete = make(chan struct{}, 1)
	)

	go func() {
		for _, h := range c.handlers {
			if err := h(ctx); err != nil {
				errs = append(errs, fmt.Sprintf("[!] %v", err))
			}
		}

		complete <- struct{}{}
	}()

	select {
	case <-complete:
		break
	case <-ctx.Done():
		return fmt.Errorf("shutdown cancelled: %v", ctx.Err())
	}

	if len(errs) != 0 {
		return fmt.Errorf(
			"shutdown finished with errors: \n%s",
			strings.Join(errs, "\n"),
		)
	}

	return nil
}

type handler func(ctx context.Context) error
