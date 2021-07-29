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

func PrintStack() {
	Stack(nil)
}
func Stack(w io.Writer) {
	if w == nil {
		w = os.Stderr
	}
	pprof.Lookup("goroutine").WriteTo(w, 2)
}
func PrintHeap() {
	Heap(nil)
}
func Heap(w io.Writer) {
	if w == nil {
		w = os.Stderr
	}
	pprof.Lookup("heap").WriteTo(w, 2)
}
func PrintThreadCreate() {
	ThreadCreate(nil)
}
func ThreadCreate(w io.Writer) {
	if w == nil {
		w = os.Stderr
	}
	pprof.Lookup("threadcreate").WriteTo(w, 2)
}
func PrintBlock() {
	Block(nil)
}
func Block(w io.Writer) {
	if w == nil {
		w = os.Stderr
	}
	pprof.Lookup("block").WriteTo(w, 2)
}