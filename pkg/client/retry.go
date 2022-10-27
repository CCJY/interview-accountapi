package client

import (
	"fmt"
	"time"
)

type Retry struct {
	RetryInterval time.Duration
	RetryMax      int
}

func Req(timeout int) bool {

	time.Sleep(time.Duration(timeout) * time.Millisecond)

	return true
}

func (r *Retry) RetryRequest() (bool, bool) {
	var result bool
	ch := make(chan bool, r.RetryMax)
	for i := 0; i < r.RetryMax; i++ {
		fmt.Printf("retrying:%d\n", i)
		fmt.Printf("Req timeout:%d\n", int(r.RetryInterval)-(int(r.RetryInterval)*i)/r.RetryMax)
		go func() {
			ch <- Req(int(r.RetryInterval) - (int(r.RetryInterval)*i)/r.RetryMax)
		}()
		select {
		case result = <-ch:
			return result, true
		case <-time.After(time.Duration(r.RetryInterval) * time.Millisecond):
		}

	}
	return false, false
}
