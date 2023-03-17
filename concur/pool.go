package concur

import (
	log "github.com/colt3k/nglog/ng"
	"math"
	"sync"
	"time"
)

// Pool is a worker group that runs a number of tasks at a
// configured concurrency.
type Pool struct {
	Tasks []*Task

	concurrency      int
	maxRunsPerSecond int
	tasksChan        chan *Task
	wg               sync.WaitGroup
}

/*
NewPool initializes a new pool with the given tasks at the given concurrency.

	concurrency			how many jobs to run at once
*/
func NewPool(tasks []*Task, concurrency int) *Pool {
	return &Pool{
		Tasks:            tasks,
		concurrency:      concurrency,
		maxRunsPerSecond: -1,
		tasksChan:        make(chan *Task),
	}
}

/*
NewPoolWithPause initializes a new pool with the given tasks at given concurrency with a pause between every X executions.

	concurrency			how many jobs to run at once
	runCountPause 		how many jobs to be processed before pause
	millisecondsPause	how long to pause in milliseconds
*/
func NewPoolWithPause(tasks []*Task, concurrency, maxRunsPerSecond int) *Pool {
	return &Pool{
		Tasks:            tasks,
		concurrency:      concurrency,
		maxRunsPerSecond: maxRunsPerSecond,
		tasksChan:        make(chan *Task),
	}
}

// Run runs all work within the pool and blocks until it's
// finished.
func (p *Pool) Run() {
	for i := 0; i < p.concurrency; i++ {
		go p.work()
	}

	p.wg.Add(len(p.Tasks))
	for _, task := range p.Tasks {
		p.tasksChan <- task
	}

	// all workers return
	close(p.tasksChan)

	p.wg.Wait()
}

// The work loop for any single goroutine.
func (p *Pool) work() {
	if p.maxRunsPerSecond > -1 {
		lastRunStart := time.Now()
		// divide runs per second by max concurrency then convert to scientific notation; change to time.Duration
		minTimeBetweenEachRun := time.Duration(math.Ceil(1e9 / (float64(p.maxRunsPerSecond) / float64(p.concurrency))))
		for task := range p.tasksChan {
			// subtract wait time from last run; if greater than 0 then wait that long before running again
			timeBeforeNextRun := -(time.Since(lastRunStart) - minTimeBetweenEachRun)
			if timeBeforeNextRun > 0 {
				log.Logf(log.DBGL3, "Worker backing off for %s", timeBeforeNextRun.String())
				time.Sleep(timeBeforeNextRun)
			}
			lastRunStart = time.Now()
			task.Run(&p.wg)
		}
	} else {
		for task := range p.tasksChan {
			task.Run(&p.wg)
		}
	}
}

type TaskResponseReturn interface {
	ProcessResponse(interface{}, error)
}

// Task encapsulates a work item that should go in a work
// pool.
type Task struct {
	// Err holds an error that occurred during a task. Its
	// result is only meaningful after Run has been called
	// for the pool that holds it.
	Err              error
	responseReturner TaskResponseReturn
	response         interface{}
	fnc              func() (interface{}, error)
}

// NewTask initializes a new task based on a given work
// function.
func NewTask(fnc func() (interface{}, error), rtn TaskResponseReturn) *Task {
	return &Task{fnc: fnc, responseReturner: rtn}
}

// Run runs a Task and does appropriate accounting via a
// given sync.WorkGroup.
func (t *Task) Run(wg *sync.WaitGroup) {
	t.response, t.Err = t.fnc()
	if t.responseReturner != nil {
		t.responseReturner.ProcessResponse(t.response, t.Err)
	}
	wg.Done()
}
