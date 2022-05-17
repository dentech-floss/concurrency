package concurrency

import (
	"errors"
	"sync"
	"time"
)

type Execution struct {
	Request  func() (interface{}, error)
	Response interface{}
	Err      error
}

func ExecuteConcurrently(
	executions []*Execution,
	timout time.Duration,
) error {

	var wg sync.WaitGroup
	waitCh := make(chan struct{})
	wg.Add(len(executions))

	go func() {
		for _, execution := range executions {
			go func(execution *Execution) {
				defer wg.Done()
				response, err := execution.Request()
				execution.Response = response
				execution.Err = err
			}(execution)
		}

		wg.Wait()
		close(waitCh)
	}()

	select { // Block until the wait group is finished or we timeout
	case <-waitCh:
	case <-time.After(timout):
		return errors.New("WaitGroup timed out")
	}

	for _, execution := range executions {
		if execution.Err != nil {
			return execution.Err
		}
	}

	return nil
}
