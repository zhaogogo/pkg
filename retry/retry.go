package main

import (
	"fmt"
	"github.com/zhaoqiang0201/pkg/clock"
	"time"
)

func main() {
	backoff := NewExponentialBackoffManager(time.Second, time.Minute, time.Minute*2, 2, 2.0, 1.0, clock.RealClock{})
	t := time.Now()

	stopCh := make(chan struct{})
	BackoffUtil(func() {
		fmt.Println("-->", time.Now().Sub(t).Seconds())
		t = time.Now()
	}, backoff, true, stopCh)
}
