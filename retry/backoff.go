package main

import (
	"fmt"
	"github.com/zhaoqiang0201/pkg/clock"
	"math"
	"time"
)

type DelayFunc func() time.Duration

type Backoff struct {
	// 初始化持续时间
	Duration time.Duration
	Factor   float64
	Jitter   float64
	Steps    int
	Cap      time.Duration
}

func (b *Backoff) Step() time.Duration {
	if b == nil {
		return 0
	}

	var nextDuration time.Duration
	nextDuration, b.Duration, b.Steps = delay(b.Steps, b.Duration, b.Cap, b.Factor, b.Jitter)
	return nextDuration
}

func (b Backoff) DelayFunc() DelayFunc {
	steps := b.Steps
	duration := b.Duration
	cap := b.Cap
	factor := b.Factor
	jitter := b.Jitter

	return func() time.Duration {
		var nextDuration time.Duration
		// jitter is applied per step and is not cumulative over multiple steps
		nextDuration, duration, steps = delay(steps, duration, cap, factor, jitter)
		return nextDuration
	}
}

// 延迟算法
func delay(steps int, duration, cap time.Duration, factor, jitter float64) (_ time.Duration, next time.Duration, nextStep int) {
	//当steps为非正数时，不更改基本持续时间
	if steps < 1 {
		if jitter > 0 {
			return Jitter(duration, jitter), duration, 0
		}
		return duration, duration, 0
	}
	steps--
	// calculate the next step's interval
	if factor != 0 {
		next = time.Duration(float64(duration) * factor)
		if cap > 0 && next > cap {
			next = cap
			steps = 0
		}
	} else {
		next = duration
	}

	// add jitter for this step
	if jitter > 0 {
		duration = Jitter(duration, jitter)
	}

	return duration, next, steps
}

// Jitter 返回持续时间, [duration, duration+maxFactor*duration]
// 如果maxFactor为0.0，则会选择建议的默认值。
func Jitter(duration time.Duration, maxFactor float64) time.Duration {
	if maxFactor <= 0.0 {
		maxFactor = 1.0
	}
	wait := duration + time.Duration(
		//rand.Float64()*maxFactor*float64(duration),
		maxFactor*float64(duration),
	)

	return wait
}

type BackoffManager interface {
	// Backoff returns a shared clock.Timer that is Reset on every invocation. This method is not
	// safe for use from multiple threads. It returns a timer for backoff, and caller shall backoff
	// until Timer.C() drains. If the second Backoff() is called before the timer from the first
	// Backoff() call finishes, the first timer will NOT be drained and result in undetermined
	// behavior.
	Backoff() (clock.Timer, int)
}

type exponentialBackoffManagerImpl struct {
	backoff              *Backoff
	backoffTimer         clock.Timer
	lastBackoffStart     time.Time
	initialBackoff       time.Duration
	backoffResetDuration time.Duration
	clock                clock.Clock
}

func (b *exponentialBackoffManagerImpl) getNextBackoff() time.Duration {
	if b.clock.Now().Sub(b.lastBackoffStart) > b.backoffResetDuration {
		b.backoff.Steps = math.MaxInt32
		b.backoff.Duration = b.initialBackoff
	}
	b.lastBackoffStart = b.clock.Now()
	return b.backoff.Step()
}

func (b *exponentialBackoffManagerImpl) Backoff() (clock.Timer, int) {
	if b.backoffTimer == nil {
		b.backoffTimer = b.clock.NewTimer(b.getNextBackoff())
	} else {
		b.backoffTimer.Reset(b.getNextBackoff())
	}
	fmt.Println("*****", b.backoff.Steps)
	return b.backoffTimer, b.backoff.Steps
}

func NewExponentialBackoffManager(initBackoff, maxbackoff, resetDuration time.Duration, step int, backoffFactor, jitter float64, c clock.Clock) BackoffManager {
	return &exponentialBackoffManagerImpl{
		backoff: &Backoff{
			Duration: initBackoff,
			Factor:   backoffFactor,
			Jitter:   jitter,
			Steps:    step,
			Cap:      maxbackoff,
		},
		backoffTimer:         nil,
		lastBackoffStart:     c.Now(),
		initialBackoff:       initBackoff,
		backoffResetDuration: resetDuration,
		clock:                c,
	}
}

func BackoffUtil(f func(), backoff BackoffManager, sliding bool, stopCh <-chan struct{}) {
	var t clock.Timer
	var s int
	for {
		select {
		case <-stopCh:
			return
		default:
		}

		if !sliding {
			t, s = backoff.Backoff()
		}

		func() {
			defer func() {
				if r := recover(); r != nil {
					fmt.Println(r)
				}
			}()
			f()
		}()

		if sliding {
			t, s = backoff.Backoff()
		}

		if s <= 0 {
			return
		}

		select {
		case <-stopCh:
			if !t.Stop() {
				<-t.C()
			}
			return
		case <-t.C():
		}
	}
}
