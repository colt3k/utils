package main

import (
	"fmt"
	"github.com/colt3k/utils/concur"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

var (
	// MaxConcurrentPerSession run without time constraints with 10 workers
	MaxConcurrentPerSession = 10

	// MaxConcurrentPerSession2 run with time constraint using 10 works
	MaxConcurrentPerSession2 = 10
	PerSecond                = 5

	jitMaxRange = 3
)

func main() {
	//ca := log.NewConsoleAppender("*")
	//log.Modify(log.LogLevel(log.DEBUG), log.Appenders(ca))
	parts := make([]*WorkObjectPart, 100)

	for i := range parts {
		parts[i] = &WorkObjectPart{Name: "someval", UniqueName: "someval" + strconv.Itoa(i)}
	}
	// Create an object that has Work to be done
	w := WorkObject{ValX: "Hello File X", Parts: parts}
	start := time.Now()
	process(&w)
	// duration
	duration := time.Since(start)
	fmt.Printf("Final Non Time Limited status: %v Seconds: %v\n", w.OverallStatus, duration.Seconds())

	fmt.Println("******************************************************************************************")
	// Create an object that has Work to be done which is Time Limited
	w2 := WorkObject{ValX: "Hello File X", Parts: parts}
	start = time.Now()
	processTimeLimited(&w2)
	// duration
	duration2 := time.Since(start)
	fmt.Printf("Final Time Limited status: %v Seconds: %v\n", w2.OverallStatus, duration2.Seconds())
}

func process(w *WorkObject) {

	//Create Worker Pool ***************************************

	var tasks []*concur.Task
	for _, p := range w.Parts {

		p := p
		// If TaskResponseReturn is null that part will be skipped
		task := concur.NewTask(
			func() (interface{}, error) {
				return worker(p)
			},
			NewReturner(w))

		tasks = append(tasks, task)
	}
	p := concur.NewPool(tasks, MaxConcurrentPerSession)
	p.Run()
	// END WORKER POOL ****************************************************

	// do any finalizing
	w.Close()
}

func processTimeLimited(w *WorkObject) {

	//Create Worker Pool ***************************************

	var tasks []*concur.Task
	for _, p := range w.Parts {

		p := p
		// If TaskResponseReturn is null that part will be skipped
		task := concur.NewTask(
			func() (interface{}, error) {
				return worker(p)
			},
			NewReturner(w))

		tasks = append(tasks, task)
	}
	p := concur.NewPoolWithPause(tasks, MaxConcurrentPerSession2, PerSecond)
	p.Run()
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
	WorkUniqueId  string
	WorkStatusVal string // success, fail, skip, etc...
}

func (w *WorkerResponse) id() string {
	return w.WorkUniqueId
}
func (w *WorkerResponse) val() string {
	return w.WorkStatusVal
}

func worker(p *WorkObjectPart) (Response, error) {
	//fmt.Printf("start work on %v\n", p.UniqueName)
	// get random val
	jitter := rand.New(rand.NewSource(time.Now().UnixNano())).Int63n(int64(jitMaxRange))
	time.Sleep(time.Duration(jitter) * time.Second)
	if jitter%2 == 0 {
		return &WorkerResponse{WorkUniqueId: p.UniqueName, WorkStatusVal: "pass"}, nil
	}

	return &WorkerResponse{WorkUniqueId: p.UniqueName, WorkStatusVal: "fail"}, fmt.Errorf("some error ABC")
}

// WorkObject store our data to be processed, array of URLs to process, a file to process for some function, etc...
type WorkObject struct {
	ValX          string // bucket
	Parts         []*WorkObjectPart
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
		if strings.Index(j.Status, "fail") > -1 || j.Error != nil {
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
	if i != nil {
		response := i.(*WorkerResponse)
		id := response.WorkUniqueId
		val := response.WorkStatusVal
		for i, j := range r.ToWorkObject.Parts {
			if j.UniqueName == id {
				r.ToWorkObject.Parts[i].Status = val
				if e != nil {
					r.ToWorkObject.Parts[i].Error = e
				}
			}
		}
	}
}
