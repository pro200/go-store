package store

import "time"

const defaultTimeout = time.Second

type config struct {
	timeout time.Duration
}

// Option configures a Store created by New.
type Option func(*config)

// WithTimeout sets how long New waits to acquire the database file
// lock before failing. Values <= 0 fall back to the default of 1s.
func WithTimeout(d time.Duration) Option {
	return func(c *config) {
		if d > 0 {
			c.timeout = d
		}
	}
}
