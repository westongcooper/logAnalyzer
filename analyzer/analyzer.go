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

	"github.com/rivo/tview"
)

type analyzer struct {
	logCount int
	logger logger
	textView *tview.TextView
	textApp *tview.Application
	reportDelay int // report delay in seconds
	peekRate    int
	alerts      alerts
	levelCount  logLevels
	// logMessages logMessages    
	lk sync.Mutex
}

type textSetter interface {
	SetText(string)
}

func NewAnalzyer(logger logger) analyzer {
	app := tview.NewApplication()
	textView := tview.NewTextView()
	return analyzer{
		logCount: 0,
		textView: textView,
		textApp: app,
		reportDelay: 1,
		logger: logger,
		levelCount: logLevels{},
		// logMessages: logMessages{},
		alerts: alerts{},
	}
}

func (a *analyzer) String()string{
	now := time.Now().Format(time.RFC3339)

	return fmt.Sprintf(`
Log Analysis Report (Last Updated: %s)
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
Runtime Stats:
• Entries Processed: %d
• Current Rate: <TODO>
• Adaptive Window: <TODO>

Pattern Analysis:
%s

Dynamic Insights:
<TODO>
• Top Errors:
<TODO>

Self-Evolving Alerts:
%s`, now, a.logCount, a.levelCount, a.alerts)
}

func (a *analyzer) IncludeLog(logLine *string) error {
	a.lk.Lock()
	defer a.lk.Unlock()

	a.logCount++

	log, err := log.NewLog(*logLine)
	if err != nil {
		a.logger.Println("failed to parse log line")
		return nil
	}

	if _, ok := a.levelCount[log.LogLevel]; ok {
		a.levelCount[log.LogLevel] ++
	} else {
		a.levelCount[log.LogLevel] = 0
	}

	// TODO: calc latest log entry rate
	// TODO: trigger any alerts
	// TODO: adjust 'reportDelay' window time if necessary
	// TODO: track error rates
	// TODO: calc any patterns

	return nil
}

func (a *analyzer) PrintUpdates(ctx context.Context) error {
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
