package main

import (
	"context"
	"fmt"
	log "github.com/colt3k/nglog/ng"
	"github.com/colt3k/utils/concur"
	"math/rand"
	"strconv"
	"strings"
	"sync"
	"time"
)

var (
	// MaxConcurrentPerSession run without time constraints with 10 workers
	MaxConcurrentPerSession = 20

	// MaxConcurrentPerSession2 run with time constraint using 10 works
	MaxConcurrentPerSession2 = 20
	PerSecond                = 10

	jitMaxRange = 3
)

func main() {
	ctx := context.Background()

	parts := make([]*WorkObjectPart, 1000)

	for i := range parts {
		parts[i] = &WorkObjectPart{Name: "someval", UniqueName: "someval" + strconv.Itoa(i)}
	}
	// Create an object that has Work to be done
	w := WorkObject{ValX: "Hello File X", Parts: parts}
	start := time.Now()
	process(ctx, &w)
	// duration
	duration := time.Since(start)
	fmt.Printf("Final Non Time Limited status: %v Seconds: %v\n", w.OverallStatus, duration.Seconds())

	fmt.Println("******************************************************************************************")
	// reinit array
	parts = make([]*WorkObjectPart, 1000)

	for i := range parts {
		parts[i] = &WorkObjectPart{Name: "someval", UniqueName: "someval" + strconv.Itoa(i)}
	}
	// Create an object that has Work to be done which is Time-Limited
	w2 := WorkObject{ValX: "Hello File X", Parts: parts}
	start = time.Now()
	processTimeLimited(ctx, &w2)
	// duration
	duration2 := time.Since(start)
	fmt.Printf("Final Time Limited status: %v Seconds: %v\n", w2.OverallStatus, duration2.Seconds())
}

func process(ctx context.Context, w *WorkObject) {
	//Create Worker Pool ***************************************
	tasks := make([]*concur.Task, len(w.Parts))
	for i, p := range w.Parts {
		p := p
		// If TaskResponseReturn is null that part will be skipped
		task := concur.NewTask(
			func() (interface{}, error) {
				return worker(p)
			},
			NewReturner(w))

		tasks[i] = task
	}
	p := concur.NewPool(ctx, tasks, MaxConcurrentPerSession)
	err := p.Run()
	if err != nil {
		log.Logf(log.ERROR, "run error %v", err)
	}
	// END WORKER POOL ****************************************************

	// do any finalizing
	w.Close()
}

func processTimeLimited(ctx context.Context, w *WorkObject) {
	//Create Worker Pool ***************************************
	tasks := make([]*concur.Task, len(w.Parts))
	for i, p := range w.Parts {
		p := p
		// If TaskResponseReturn is null that part will be skipped
		task := concur.NewTask(
			func() (interface{}, error) {
				return worker(p)
			},
			NewReturner(w))

		tasks[i] = task
	}
	p := concur.NewPoolWithPause(ctx, tasks, MaxConcurrentPerSession2, PerSecond)
	err := p.Run()
	if err != nil {
		log.Logf(log.ERROR, "run error %v", err)
	}
	// END WORKER POOL ****************************************************

	// do any finalizing
	w.Close()
}

// Response setup our Response Object/Interface
type Response interface {
	id() string
	val() string
}

// WorkerResponse our concreate implementation of Response interface
type WorkerResponse struct {
	WorkUniqueID  string
	WorkStatusVal string // success, fail, skip, etc...
}

func (w *WorkerResponse) id() string {
	return w.WorkUniqueID
}
func (w *WorkerResponse) val() string {
	return w.WorkStatusVal
}

func worker(p *WorkObjectPart) (Response, error) {
	//fmt.Printf("start work on %v\n", p.UniqueName)
	// get random val
	jitter := rand.New(rand.NewSource(time.Now().UnixNano())).Int63n(int64(jitMaxRange)) //nolint:gosec
	time.Sleep(time.Duration(jitter) * time.Second)
	if jitter%2 == 0 {
		return &WorkerResponse{WorkUniqueID: p.UniqueName, WorkStatusVal: "pass"}, nil
	}

	return &WorkerResponse{WorkUniqueID: p.UniqueName, WorkStatusVal: "fail"}, fmt.Errorf("some error ABC")
}

// WorkObject store our data to be processed, array of URLs to process, a file to process for some function, etc...
type WorkObject struct {
	ValX          string // bucket
	Parts         []*WorkObjectPart
	PartsMutex    sync.RWMutex
	OverallStatus string
	Error         error
}

type WorkObjectPart struct {
	UniqueName string
	Name       string // object or url
	Status     string
	Error      error
}

// Close here we can ensure everything is completed and close up any tasks or other things, close db, write out data, etc...
func (w *WorkObject) Close() {
	fmt.Println("Closing out our work...")
	failed := false
	for _, j := range w.Parts {
		//fmt.Printf("workpart %v status: %v\n", j.UniqueName, j.Status)
		if strings.Contains(j.Status, "fail") || j.Error != nil {
			w.OverallStatus = "fail"
			failed = true
		}
	}
	if !failed {
		w.OverallStatus = "pass"
	}
}

// NewReturner create a new Returner struct which is used to capture and set a value returned from our worker into our original object
func NewReturner(w *WorkObject) *Returner {
	tmp := &Returner{ToWorkObject: w}
	return tmp
}

// Returner stores a pointer to our original work object
type Returner struct {
	ToWorkObject *WorkObject
}

// ProcessResponse set worker response on our original work object
func (r *Returner) ProcessResponse(i interface{}, e error) {
	completeCount := 0
	failCount := 0
	successCount := 0
	if i != nil {
		response := i.(*WorkerResponse)
		id := response.WorkUniqueID
		val := response.WorkStatusVal
		r.ToWorkObject.PartsMutex.Lock()
		for x, j := range r.ToWorkObject.Parts {
			if j.UniqueName == id {
				r.ToWorkObject.Parts[x].Status = val
				if e != nil {
					r.ToWorkObject.Parts[x].Error = e
				}
			}
			if len(strings.TrimSpace(j.Status)) > 0 {
				completeCount++
				if j.Status == "pass" {
					successCount++
				} else {
					failCount++
				}
			}
		}
		r.ToWorkObject.PartsMutex.Unlock()
		log.Logf(log.INFO, "%v of %v complete (success %v, fail %v)", completeCount, len(r.ToWorkObject.Parts), successCount, failCount)
	}
}
