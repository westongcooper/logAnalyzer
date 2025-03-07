package analyzer

import (
	"fmt"
)

type logLevels map[string]int

func (logLevels logLevels) String() string {
	message := ""
	total := 0
	for _, count := range logLevels {
		total+=count
	}

	for level, count := range logLevels {
		precent := float64(count) / float64(total) * 100
		message = message + fmt.Sprintf(" - %s: %.2f%% (%d entries)\n", level, precent, count)
	}
	return message
}

// type logMessages map[string]int

// func (logMessages logMessages) String() string {
// 	// for level, count := range logMessages { // TODO: include to log messages
// 	// 	precent := float64(count) / float64(total) * 100
// 	// 	message = message + fmt.Sprintf(" - %s: %.2f%% (%d entries)\n", level, precent, count)
// 	// }
// 	// return message
// 	return fmt.Sprintf("%v", logMessages)
// }