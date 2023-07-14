package main

import (
	"fmt"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/utils/clock"
	"testing"
	"time"
)

func TestBackoff(t *testing.T) {
	c := clock.RealClock{}
	backoff := wait.NewExponentialBackoffManager(time.Second, time.Minute, time.Minute*2, 2.0, 1.0, c)

	stopCh := make(chan struct{})
	wait.BackoffUntil(func() {
		fmt.Println("-->", time.Now())
	}, backoff, false, stopCh)
}
