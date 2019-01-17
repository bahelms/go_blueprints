package trace

import (
	"fmt"
	"io"
)

// Tracer is the interface describing an object capable of tracing events
type Tracer interface {
	Trace(...interface{})
}

type tracer struct {
	out io.Writer
}

func (t tracer) Trace(a ...interface{}) {
	fmt.Fprint(t.out, a...)
	fmt.Fprintln(t.out)
}

type nilTracer struct{}

func (t nilTracer) Trace(a ...interface{}) {}

// New creates a new tracer with given writer
func New(w io.Writer) Tracer {
	return tracer{out: w}
}

// Off creates a Tracer that ignores calls to Trace
func Off() Tracer {
	return nilTracer{}
}
