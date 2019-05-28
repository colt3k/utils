package profile

//#include <time.h>
import "C"
import (
	"runtime"
	"time"

	log "github.com/colt3k/nglog/ng"
)

/*
Duration Arguments to a defer statement is immediately evaluated and stored.
The deferred function receives the pre-evaluated values when its invoked.
Example defer profile.Duration(time.Now(), "IntFactorial")
*/
func Duration(invocation time.Time, name string) {
	elapsed := time.Since(invocation)

	log.Printf("%s lasted %s", name, elapsed)
}

/* MemUsage print out memory usage
https://golangcode.com/print-the-current-memory-usage/
https://golang.org/src/runtime/mstats.go
 */
func MemUsage() {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	// For info on each, see: https://golang.org/pkg/runtime/#MemStats
	log.Logf(log.DEBUG, "Alloc = %v MiB\tTotalAlloc = %v MiB\tSys = %v MiB\tNumGC = %v", bToMb(m.Alloc), bToMb(m.TotalAlloc), bToMb(m.Sys), m.NumGC)
}

func bToMb(b uint64) uint64 {
	return b / 1024 / 1024
}

var startTime = time.Now()
var startTicks = C.clock()

func CpuUsagePercent() {
	clockSeconds := float64(C.clock()-startTicks) / float64(C.CLOCKS_PER_SEC)
	realSeconds := time.Since(startTime).Seconds()
	usage := clockSeconds / realSeconds * 100
	log.Logf(log.DEBUG,"CPU %f\n", usage)
}