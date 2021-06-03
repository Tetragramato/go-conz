package internal

import (
	"github.com/go-resty/resty/v2"
	"log"
	"sync"
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
