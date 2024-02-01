// Package logging provides the operators's recommended logging interface.
//
// The logging interface avoids the complexity of levels and provides a simpler
// api that makes it harder to introduce unnecesasry ambiguity to logs (or
// ascribing value to arbitrary magic numbers).
//
// An Error logging helper exists primarily to facilitate including a stack trace
// when the backing provider supports it.
package logging

import (
	"context"

	"github.com/go-logr/logr"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

// A Logger logs messages. Messages may be supplemented by structured data.
type Logger interface {
	// Info logs a message with optional structured data. Structured data must
	// be supplied as an array that alternates between string keys and values of
	// an arbitrary type. Use Info for messages that users are
	// very likely to be concerned with.
	Info(msg string, keysAndValues ...any)

	// Error logs a message with optional structured data. Structured data must
	// be supplied as an array that alternates between string keys and values of
	// an arbitrary type. Use Error when you want to enrich a message with as much
	// information as a logging provider can.
	Error(err error, msg string, keysAndValues ...any)

	// Debug logs a message with optional structured data. Structured data must
	// be supplied as an array that alternates between string keys and values of
	// an arbitrary type. Use Debug for messages that operators or
	// developers may be concerned with when debugging the operator or spin.
	Debug(msg string, keysAndValues ...any)

	// WithValues returns a Logger that will include the supplied structured
	// data with any subsequent messages it logs. Structured data must
	// be supplied as an array that alternates between string keys and values of
	// an arbitrary type.
	WithValues(keysAndValues ...any) Logger
}

// NewNopLogger returns a Logger that does nothing.
func NewNopLogger() Logger { return nopLogger{} }

type nopLogger struct{}

func (l nopLogger) Info(_ string, _ ...any)           {}
func (l nopLogger) Debug(_ string, _ ...any)          {}
func (l nopLogger) Error(_ error, _ string, _ ...any) {}
func (l nopLogger) WithValues(_ ...any) Logger        { return nopLogger{} }

// NewLogrLogger returns a Logger that is satisfied by the supplied logr.Logger,
// which may be satisfied in turn by various logging implementations.
// Debug messages are logged at V(1) - following the reccomendation of
// controller-runtime.
func NewLogrLogger(l logr.Logger) Logger {
	return logrLogger{log: l}
}

type logrLogger struct {
	log logr.Logger
}

func (l logrLogger) Info(msg string, keysAndValues ...any) {
	l.log.Info(msg, keysAndValues...)
}

func (l logrLogger) Error(err error, msg string, keysAndValues ...any) {
	l.log.Error(err, msg, keysAndValues...)
}

func (l logrLogger) Debug(msg string, keysAndValues ...any) {
	l.log.V(1).Info(msg, keysAndValues...)
}

func (l logrLogger) WithValues(keysAndValues ...any) Logger {
	return logrLogger{log: l.log.WithValues(keysAndValues...)}
}

func FromContext(ctx context.Context) Logger {
	return logrLogger{log: log.FromContext(ctx)}
}
