package reporter

import (
	"errors"
	"fmt"

	"github.com/jhump/protocompile/ast"
)

// ErrInvalidSource is a sentinel error that is returned by compilation and
// stand-alone compilation steps (such as parsing, linking) when one or more
// errors is reported but the configured ErrorReporter always returns nil.
var ErrInvalidSource = errors.New("parse failed: invalid proto source")

// ErrorWithPos is an error about a proto source file that adds information
// about the location in the file that caused the error.
type ErrorWithPos interface {
	error
	// GetPosition returns the source position that caused the underlying error.
	GetPosition() ast.SourcePos
	// Unwrap returns the underlying error.
	Unwrap() error
}

// Error creates a new ErrorWithPos from the given error and source position.
func Error(pos ast.SourcePos, err error) ErrorWithPos {
	return errorWithSourcePos{pos: pos, underlying: err}
}

// Errorf creates a new ErrorWithPos whose underlying error is created using the
// given message format and arguments (via fmt.Errorf).
func Errorf(pos ast.SourcePos, format string, args ...interface{}) ErrorWithPos {
	return errorWithSourcePos{pos: pos, underlying: fmt.Errorf(format, args...)}
}

type errorWithSourcePos struct {
	underlying error
	pos        ast.SourcePos
}

func (e errorWithSourcePos) Error() string {
	sourcePos := e.GetPosition()
	return fmt.Sprintf("%s: %v", sourcePos, e.underlying)
}

func (e errorWithSourcePos) GetPosition() ast.SourcePos {
	return e.pos
}

func (e errorWithSourcePos) Unwrap() error {
	return e.underlying
}

var _ ErrorWithPos = errorWithSourcePos{}
