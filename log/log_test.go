package log

import (
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func fakeLog() log {
	return log{
		IPAddress: fmt.Sprintf("%d.%d.%d.%d", rand.Intn(128), rand.Intn(128), rand.Intn(128), rand.Intn(128)),
		LogLevel: []string{"ERROR", "INFO", "DEBUG"}[rand.Intn(3)],
		Date: time.Now().Add(-(time.Second*time.Duration(rand.Intn(300)))).Truncate(time.Second),
		ErrorCode: fmt.Sprintf("%d", rand.Intn(999)),
		Message: fmt.Sprintf("Error message - %d", rand.Intn(100000)),
	}
}

func errorLogToString(errorLog log) string {
	return fmt.Sprintf(`[%s] %s - IP:%s Error %s - %s`, errorLog.Date.Format("2006-01-02T15:04:05Z07:00"), errorLog.LogLevel, errorLog.IPAddress, errorLog.ErrorCode, errorLog.Message)
}

func TestNewLog(t *testing.T) {
	t.Run("happy path", func(t *testing.T) {
		expectedLog := fakeLog()
		
		log, err := NewLog(errorLogToString(expectedLog))
		assert.NoError(t, err)
		assert.Equal(t, expectedLog, log)
	})
}