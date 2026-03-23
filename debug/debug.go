// Package debug writes selected pprof snapshots to stderr or a caller-provided
// writer.
package debug

/*
	options:
      goroutine:    stack traces of all current goroutines
      heap:         sampling of all heap allocations
      threadcreate: stack traces that led to the creation of new OS threads
      block:        stack traces that led to blocking on synchronization primitives
*/
import (
	"io"
	"os"
	"runtime/pprof"
)

// PrintStack writes goroutine stack traces to stderr.
func PrintStack() {
	Stack(nil)
}

// Stack writes goroutine stack traces to w, defaulting to stderr.
func Stack(w io.Writer) {
	if w == nil {
		w = os.Stderr
	}
	pprof.Lookup("goroutine").WriteTo(w, 2)
}

// PrintHeap writes the heap profile to stderr.
func PrintHeap() {
	Heap(nil)
}

// Heap writes the heap profile to w, defaulting to stderr.
func Heap(w io.Writer) {
	if w == nil {
		w = os.Stderr
	}
	pprof.Lookup("heap").WriteTo(w, 2)
}

// PrintThreadCreate writes thread creation traces to stderr.
func PrintThreadCreate() {
	ThreadCreate(nil)
}

// ThreadCreate writes thread creation traces to w, defaulting to stderr.
func ThreadCreate(w io.Writer) {
	if w == nil {
		w = os.Stderr
	}
	pprof.Lookup("threadcreate").WriteTo(w, 2)
}

// PrintBlock writes the blocking profile to stderr.
func PrintBlock() {
	Block(nil)
}

// Block writes the blocking profile to w, defaulting to stderr.
func Block(w io.Writer) {
	if w == nil {
		w = os.Stderr
	}
	pprof.Lookup("block").WriteTo(w, 2)
}
