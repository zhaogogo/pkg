package retry

import (
	"github.com/zhaoqiang0201/pkg/clock"
	"math"
	"testing"
	"time"
)

func TestBackoff(t *testing.T) {
	backoff := NewExponentialBackoffManager(time.Second, time.Minute, time.Minute*2, math.MaxInt32, 2.0, 1.0, clock.RealClock{})
	t.Log(">>>>", time.Now())
	stopCh := make(chan struct{})
	BackoffUtil(func() {
		t.Log("-->", time.Now())
	}, backoff, true, stopCh)
}
