package main

import (
	"context"
	"flag"
	"fmt"
	logger "log"
	"logAnalyzer/analyzer"
	"logAnalyzer/parser"
	"logAnalyzer/tail"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	var (
		filePath string
		tailEnd  bool
	)

	flag.StringVar(&filePath, "file",    "test_logs.log", "log file to read")          // 'test_logs.log' is hard coded in genLogs.sh requirements
	flag.BoolVar(&tailEnd,    "tailEnd", true,            "Read from end of log file") // There may be a file already saved, should we only analyze new data?

	flag.Parse()

	// ======== Initialize logger streamer
	logger := logger.New(os.Stderr, "", logger.LstdFlags)

	// ======== Create main context
	ctx, cancel := context.WithCancel(context.Background())

	// ======== tail file, read lines
	lines, err := tail.TailFile(ctx, filePath, tailEnd, logger)
	if err != nil {
		logger.Panicf("failed to load file (%s)", filePath)
	}

	// ======== setup graceful shutdown
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		<-c
		
		fmt.Println("\n******\nShutdown signal received\n******")
		cancel()
	}()

	analyzer := analyzer.NewAnalzyer(logger)

	analyzer.PrintUpdates(ctx)

	// ======== Set up parser to conditionally include lines
	parsables := []parser.Parsable{
		parser.NewParsable(`^\[.{20}] `, analyzer.IncludeLog),
	}

	if err := parser.ParseLines(ctx, lines, 30, logger, parsables...); err != nil {
		fmt.Println("Process lines stopped: " + err.Error())
	}

	// ======== we're done!
	logger.Println("app is closed")
}
