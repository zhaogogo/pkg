package main

import (
	"fmt"
	"math"
	"time"
)

func main() {
	backoff := NewExponentialBackoffManager(time.Second, time.Minute, time.Minute*2, math.MaxInt32, 2.0, 1.0, clock.RealClock{})
	fmt.Println("-->", time.Now())
	stopCh := make(chan struct{})
	BackoffUtil(func() {
		fmt.Println("-->", time.Now())
	}, backoff, true, stopCh)
}
