package analyzer

import (
	"fmt"
	"time"
)

type alerts []alert

type alert struct {
	date    time.Time
	message string
}

func (alert alert)String() string {
	return fmt.Sprintf("[%s] ⚠️ %s\n", alert.date.Format("15:04:05"), alert.message)
}

func (alerts alerts)String() string {
	if len(alerts) == 0 {
		return "< no alerts >"
	}
	message := ""
	for _, alert := range alerts {
		message+= alert.String()
	}
	return message
}