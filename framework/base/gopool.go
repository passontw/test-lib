package base

import (
	"fmt"
	"sl.framework.com/async"
)

// this is a gogroutine pool
type GoPool struct {
	chTasks      chan Tasker
	chWorkerChan chan chan Tasker
	workerCount  uint32
}

func NewGoPool(poolSize uint32) *GoPool {
	if poolSize <= 0 || poolSize > 9999 {
		panic("poolsize invalid, size[0, 9999]")
	}
	return &GoPool{
		chTasks:      make(chan Tasker, poolSize),
		chWorkerChan: make(chan chan Tasker),
		workerCount:  poolSize}
}

func (g *GoPool) workerReady(w chan Tasker) {
	g.chWorkerChan <- w
}

func (g *GoPool) workerChan() chan Tasker {
	return make(chan Tasker)
}

func (g *GoPool) createWorker(id int, in chan Tasker) {
	fn := func() {
		for {
			g.workerReady(in)
			t := <-in
			t.RunTask()
		}
	}
	async.AsyncRunCoroutine(fn)
}

func (g *GoPool) Run() {
	for i := 0; i < int(g.workerCount); i++ {
		g.createWorker(i, g.workerChan())
	}

	fn := func() {
		var requestQueue []Tasker
		var workerQueue []chan Tasker
		for {
			var activeRequest Tasker
			var activeWorker chan Tasker
			if len(requestQueue) > 0 && len(workerQueue) > 0 {
				activeRequest = requestQueue[0]
				activeWorker = workerQueue[0]
			}
			select {
			case r := <-g.chTasks:
				requestQueue = append(requestQueue, r)
			case w := <-g.chWorkerChan:
				workerQueue = append(workerQueue, w)
			case activeWorker <- activeRequest:
				workerQueue = workerQueue[1:]
				requestQueue = requestQueue[1:]
			}
		}
	}
	async.AsyncRunCoroutine(fn)
}

func (g *GoPool) worker(wid uint32) {
	for t := range g.chTasks {
		t.RunTask()
		fmt.Println("worker:", wid)
	}
}

func (g *GoPool) AsyncRunFn(fn func()) {
	g.chTasks <- &taskEmpty{fn: fn}
}

func (g *GoPool) AsyncRunTask(t Tasker) {
	g.chTasks <- t
}
