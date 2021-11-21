// Package closers provides utilities for working with io.Closer implementations.
package closers

import (
	"io"
	"log"
)

// Close the provided io.Closer implementation, logging if there's an error. This is typically used with a defer
// statement.
func Close(c io.Closer) {
	if err := c.Close(); err != nil {
		log.Printf("failed to close %T: %v\n", c, err)
	}
}

// CloseFunc is used as an adaptor to create io.Closer implementations for functions.
type CloseFunc func() error

// Close the CloseFunc.
func (c CloseFunc) Close() error {
	return c()
}

// Noop is an io.Closer implementation that returns nil on Close.
var Noop = CloseFunc(func() error {
	return nil
})
