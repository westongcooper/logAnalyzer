package analyzer

import (
	"context"
	"fmt"
	"logAnalyzer/log"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"sync"
	"time"
)

var originalSlidingWindow = 5 // in seconds
var errorRateThreshold    = 5  // how many per second

type analyzer struct {
	logger         logger
	reportDelay    int // report delay in seconds
	alerts         alerts
	totalCounter   traceable_i
	levelCounter   rateCounters
	messageCounter rateCounters
	slidingWindow  int
	lk sync.Mutex
}

type rateCounters interface {
	IncrOne(label string)
	Rate(label string) int64
	GetTotal(label string) uint64
	SnapShot()
}

type textSetter interface {
	SetText(string)
}

func NewAnalzyer(logger logger) analyzer {
	return analyzer{
		reportDelay: originalSlidingWindow,
		logger: logger,
		levelCounter: traceables{},
		messageCounter: traceables{},
		totalCounter: newTraceable(),
		slidingWindow: originalSlidingWindow,
		alerts: alerts{},
	}
}

func (a *analyzer) String()string{
	return fmt.Sprintf(`
Log Analysis Report (Last Updated: %s)
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
Runtime Stats:
• Entries Processed: %d
• Current Rate: %d entries/sec (Peak: %d entries/sec)
• Adaptive Window: %d sec%s

Pattern Analysis:
%s

Dynamic Insights:
• Error Rate: %d errors/sec
<TODO: Emerging patterns>
• Top Errors:
<TODO>

Self-Evolving Alerts:
%s`, 
	time.Now().Format(time.RFC3339), 
	a.totalCounter.GetTotal(), 
	a.totalCounter.Rate(),
	a.totalCounter.GetPeekRate(),
	a.reportDelay,
	"", // TODO: update window size and show 
	ToLogLevelString(a.levelCounter.(traceables)), 
	a.levelCounter.Rate("ERROR"),
	a.alerts)
}

func (a *analyzer) analyzeRates() {
	a.totalCounter.SnapShot()
	a.levelCounter.SnapShot()
	a.messageCounter.SnapShot()

	// TODO: trigger any alerts
	// TODO: adjust 'reportDelay' window time if necessary
}

func (a *analyzer) IncludeLog(logLine *string) error {
	a.lk.Lock()
	defer a.lk.Unlock()

	a.totalCounter.IncrOne()

	log, err := log.NewLog(*logLine)
	if err != nil {
		a.logger.Println("failed to parse log line")
		return nil
	}

	a.levelCounter.IncrOne(log.LogLevel)
	a.messageCounter.IncrOne(log.Message)

	return nil
}

func (a *analyzer) Analyze(ctx context.Context) error {
	go func() {
		for ctx.Err() == nil {
			time.Sleep(time.Second)
			a.analyzeRates()
		}
	}()

	go func() {
		for ctx.Err() == nil {
			clearScreen()
			fmt.Printf("%s\n%s\nPress Ctrl+C to exit", a, strings.Repeat("-", 80))

			time.Sleep(time.Second * time.Duration(a.reportDelay))
		}
	}()

	return nil
}

func clearScreen() {
	command := "clear"
	args := []string{}
    if runtime.GOOS == "windows" {
		command = "cmd"
        args = []string{"/c", "cls"}
    }

     cmd := exec.Command(command, args...)
     if cmd.Err != nil {
        return
     }

     cmd.Stdout = os.Stdout
     err :=  cmd.Run()
     if err != nil {
        return
     }
}

type logger interface {
	Println(v ...interface{})
}
