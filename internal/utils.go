package internal

import (
	"github.com/go-resty/resty/v2"
	"log"
	"sync"
	"time"
)

func Trace(resp *resty.Response, err error) {
	log.Println("Response Info:")
	log.Println("  URL        :", resp.Request.URL)
	log.Println("  Error      :", err)
	log.Println("  Status Code:", resp.StatusCode())
	log.Println("  Status     :", resp.Status())
	log.Println("  Proto      :", resp.Proto())
	log.Println("  Time       :", resp.Time())
	log.Println("  Received At:", resp.ReceivedAt())
	//log.Println("  Body       :\n", resp)
	log.Println()
}

// Parallelize parallelize the function calls
func Parallelize(functions ...func()) {
	var waitGroup sync.WaitGroup
	waitGroup.Add(len(functions))

	defer waitGroup.Wait()

	for _, function := range functions {
		go func(copy func()) {
			defer waitGroup.Done()
			copy()
		}(function)
	}
}

// Contains Check value existence in a slice
func Contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}
	return false
}

type poller struct {
	ticker   *time.Ticker
	function func()
}

// RunNewPoller Run a new Poller around a lambda function, that tick every time.Duration
func RunNewPoller(timeInterval time.Duration, function func()) *poller {
	ticker := time.NewTicker(timeInterval)
	defer ticker.Stop()
	poll := &poller{
		ticker:   ticker,
		function: function,
	}
	poll.Run()
	return poll
}

func (p *poller) Run() {
	for ; true; <-p.ticker.C {
		log.Println("Tick for polling...")
		p.function()
	}
}
