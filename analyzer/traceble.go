package analyzer

import (
	"fmt"
	"time"

	"github.com/paulbellamy/ratecounter"
)

type traceable_i interface {
	GetTotal() uint64
	IncrOne()
	Rate() int64
	GetPeekRate() int64
	SnapShot()
}

type rateCounter interface {
	Incr(val int64)
	Rate() int64
}

type traceable struct {
	rateCounter
	total   uint64
	peek    uint64
}

func newTraceable() *traceable{
	return &traceable{
		total: 0,
		rateCounter: ratecounter.NewRateCounter(1 * time.Second),
	}
}

type traceables map[string]*traceable
func (traceables traceables) SnapShot(){
	for _, traceable := range traceables {
		traceable.SnapShot()
	}
}

func (traceables traceables) Rate(label string) int64 {
	if tracable, ok := traceables[label]; ok {
		return tracable.Rate()
	} else {
		return 0
	}
}

func (traceables traceables) IncrOne(label string) {
	if tracable, ok := traceables[label]; ok {
		tracable.Incr(1)
		tracable.total++
	} else {
		tracable := newTraceable()
		tracable.total = 1
		traceables[label] = tracable

	}
}

func (traceables traceables) GetTotal(label string) uint64 {
	if tracable, ok := traceables[label]; ok {
		return tracable.total
	} else {
		return 0
	}
}

func (traceable *traceable) IncrOne() {
	traceable.total++
	traceable.Incr(1)
}

func (traceable *traceable) SnapShot(){
	currentRate := traceable.Rate()
	if currentRate > int64(traceable.peek) {
		traceable.peek = uint64(currentRate)
	}
}

func (traceable traceable) GetPeekRate() int64 {
	return int64(traceable.peek)
}

func (traceable traceable) GetTotal() uint64 {
	return traceable.total
}

func ToLogLevelString(traceables traceables) string {
	message := ""
	allTotal := uint64(0)
	for _, traceable := range traceables {
		allTotal+=traceable.total
	}

	for label, traceable := range traceables {
		precent := float64(traceable.total) / float64(allTotal) * 100
		message = message + fmt.Sprintf(" - %s: %.2f%% (%d entries)\n", label, precent, traceable.total)
	}
	return message
}