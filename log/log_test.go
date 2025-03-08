package log

import (
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func fakeLog() log {
	log := log{
		IPAddress: fmt.Sprintf("%d.%d.%d.%d", rand.Intn(128), rand.Intn(128), rand.Intn(128), rand.Intn(128)),
		LogLevel: []string{"ERROR", "INFO", "DEBUG"}[rand.Intn(3)],
		Date: time.Now().Add(-(time.Second*time.Duration(rand.Intn(300)))).Truncate(time.Second),
	}

	return log
}

func fakeLogWithMessage() log {
	log := fakeLog()
	log.Message = fmt.Sprintf("Error %d message - oops-%d", rand.Intn(500), rand.Intn(100000))
	return log
}

func logToString(errorLog log) string {
	return fmt.Sprintf(`[%s] %s - IP:%s %s`, errorLog.Date.Format("2006-01-02T15:04:05Z07:00"), errorLog.LogLevel, errorLog.IPAddress, errorLog.Message)
}

func TestNewLog(t *testing.T) {
	t.Run("happy path", func(t *testing.T) {
		expectedLog := fakeLog()
		
		log, err := NewLog(logToString(expectedLog))
		assert.NoError(t, err)
		assert.Equal(t, expectedLog, log)
	})

	t.Run("happy path - with message", func(t *testing.T) {
		expectedLog := fakeLogWithMessage()
		
		log, err := NewLog(logToString(expectedLog))
		assert.NoError(t, err)
		assert.Equal(t, expectedLog, log)
	})
}