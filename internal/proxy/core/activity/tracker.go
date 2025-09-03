package activity

import (
	"context"
	"sync"
	"time"
)

var (
	// if keep changing move to config
	warnIntervals = []time.Duration{30 * time.Second, 20 * time.Second, 10 * time.Second}
)

type Tracker struct {
	ctx  context.Context
	mu   sync.Mutex
	once sync.Once

	lastActive time.Time
	timeout    time.Duration

	stopCh chan struct{}

	onTimeout func()
	onWarning func(int)

	warned map[time.Duration]bool
}

func NewTracker(ctx context.Context, timeout time.Duration, onTimeout func(), onWarning func(int)) *Tracker {
	if onTimeout == nil {
		panic("onTimeout cannot be nil")
	}

	t := &Tracker{
		ctx:        ctx,
		lastActive: time.Now(),
		timeout:    timeout * time.Second,
		stopCh:     make(chan struct{}),
		onTimeout:  onTimeout,
		onWarning:  onWarning,
		warned:     make(map[time.Duration]bool),
	}

	go t.monitor()
	return t
}

// Touch updates the last active time to now
func (t *Tracker) Touch() {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.lastActive = time.Now()

	// reset warn
	if len(t.warned) > 0 {
		t.warned = make(map[time.Duration]bool)
	}
}

func (t *Tracker) Stop() {
	t.once.Do(func() {
		close(t.stopCh)
	})
}

func (t *Tracker) monitor() {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			var (
				doTimeout     bool
				shouldWarn    bool
				secsRemaining int
				onWarnCb      func(int)
				onTimeoutCb   func()
			)

			now := time.Now()

			t.mu.Lock()
			elapsed := now.Sub(t.lastActive)
			inactive := elapsed > t.timeout

			// safety
			onWarnCb = t.onWarning
			onTimeoutCb = t.onTimeout

			// check if should warn
			if onWarnCb != nil {
				for _, interval := range warnIntervals {
					if interval >= t.timeout {
						continue
					}
					warnTime := t.timeout - interval
					if elapsed >= warnTime && !t.warned[interval] {
						t.warned[interval] = true
						shouldWarn = true
						// max 1 warning
						break
					}
				}
			}

			if inactive {
				doTimeout = true
			}

			left := t.timeout - elapsed
			if left < 0 {
				left = 0
			}
			secsRemaining = int(left.Seconds())

			// unlock before callbacks else can deadlock then gg complain again
			t.mu.Unlock()

			// fire callbacks
			if shouldWarn && onWarnCb != nil {
				onWarnCb(secsRemaining)
			}

			if doTimeout {
				if onTimeoutCb != nil {
					onTimeoutCb()
				}
				return
			}

		case <-t.stopCh:
			return
		case <-t.ctx.Done():
			return
		}
	}
}
