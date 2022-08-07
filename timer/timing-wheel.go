/* For license and copyright information please see LEGAL file in repository */

package timer

import (
	"github.com/GeniusesGroup/libgo/protocol"
)

// TimingWheel is not concurrent safe and must call each instance by each CPU core separately,
// or write upper layer to implement needed logic to prevent data race.
type TimingWheel struct {
	wheelSize int
	wheel     [][]*Timer
	interval  protocol.Duration // same as tw.ticker.period
	pos       int
	ticker    Timer
	stop      chan struct{}
}

// if you make one sec interval for a earth day(60*60*24=86400), you need 2MB ram just to hold empty wheel without any timer.
func (tw *TimingWheel) Init(interval protocol.Duration, wheelSize int) {
	tw.wheelSize = wheelSize
	tw.stop = make(chan struct{})
	tw.wheel = make([][]*Timer, wheelSize)
	tw.ticker.Init(nil)
	tw.interval = interval
}

func (tw *TimingWheel) AddTimer(t *Timer) {
	var addedPosition = tw.addedPosition(t)
	if addedPosition > tw.wheelSize {
		panic("timer - wheel: try to add a timer with bad timeout that overflow the current timing wheel")
	}
	if addedPosition == tw.pos {
		t.callback.TimerHandler()
		tw.checkAndAddTimerAgain(t)
		return
	}
	tw.wheel[addedPosition] = append(tw.wheel[addedPosition], t)
}

// call by go keyword if you don't want the current goroutine block.
func (tw *TimingWheel) Start() {
	tw.ticker.Tick(tw.interval, tw.interval, -1)
loop:
	for {
		select {
		case <-tw.ticker.Signal():
			var pos = tw.pos
			tw.incrementTickPosition()
			var timers = tw.wheel[pos]
			tw.wheel[pos] = timers[:0]
			for i := 0; i < len(timers); i++ {
				var timer = timers[i]
				timer.callback.TimerHandler()
				tw.checkAndAddTimerAgain(timer)
			}
		case <-tw.stop:
			close(tw.stop)
			tw.stop = nil
			break loop
		}
	}
	tw.ticker.Stop()
}

// Not concurrent safe.
func (tw *TimingWheel) Stop() (alreadyStopped bool) {
	if tw.stop == nil {
		return true
	}

	select {
	case tw.stop <- struct{}{}:
	default:
		alreadyStopped = true
	}
	tw.ticker.Stop()
	return
}

func (tw *TimingWheel) incrementTickPosition() {
	var wheelLen = len(tw.wheel)
	if wheelLen-1 == tw.pos {
		tw.pos = 0
	} else {
		tw.pos++
	}
}

func (tw *TimingWheel) checkAndAddTimerAgain(t *Timer) {
	if t.periodNumber == 0 {
		t.reset()
	} else {
		var addedPosition = tw.addedPosition(t)
		tw.wheel[addedPosition] = append(tw.wheel[addedPosition], t)
		if t.periodNumber > 0 {
			t.periodNumber--
		}
	}
}

func (tw *TimingWheel) addedPosition(t *Timer) int {
	return int(t.period/tw.interval) + tw.pos
}
