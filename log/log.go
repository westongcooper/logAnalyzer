package log

import (
	"fmt"
	"net"
	"regexp"
	"time"
)


type log struct {
	IPAddress  string    `gorm:"type:varchar(15)"`
	Date       time.Time
	LogLevel   string
	ErrorCode  string
	Message    string
}

/* log line examples
[2025-03-07T18:52:04Z] ERROR - IP:192.168.82.198 Error 500 - Database connection failed
[2025-03-07T18:52:04Z] INFO - IP:192.168.53.86 
[2025-03-07T18:52:04Z] DEBUG - IP:192.168.38.66 
[2025-03-07T18:52:04Z] ERROR - IP:192.168.164.29 Error 500 - Illegal argument provided
[2025-03-07T18:52:04Z] ERROR - IP:192.168.106.59 Error 500 - Out of memory
[2025-03-07T18:52:04Z] ERROR - IP:192.168.189.144 Error 500 - Access denied
[2025-03-07T18:52:04Z] ERROR - IP:192.168.249.140 Error 500 - Network timeout occurred
[2025-03-07T18:52:04Z] ERROR - IP:192.168.150.11 Error 500 - File not found
[2025-03-07T18:52:04Z] ERROR - IP:192.168.9.208 Error 500 - Null pointer exception
*/

var logRegex = regexp.MustCompile(`^\[(?P<date>.*?)\] (?P<loglevel>.*?) - IP:(?P<ipAddress>(?:\d{1,3}.){3}\d{1,3})\s*(?P<logMessage>.*)$`)

func extractTextFromLine(regexp *regexp.Regexp, text string) map[string]string {
	match := regexp.FindStringSubmatch(text)
    result := make(map[string]string)
    for i, name := range regexp.SubexpNames() {
        if i != 0 && name != "" {
            result[name] = match[i]
        }
    }
	return result
}

func NewLog(logLine string) (log, error) {
	logData := extractTextFromLine(logRegex, logLine)

	log := log{}
	var err error
	if val, ok := logData["date"]; ok {
		log.Date, err = time.Parse(time.RFC3339, val)
		if err != nil {
			return log, err
		}
	}

	if val, ok := logData["loglevel"]; ok {
		log.LogLevel = val
	}

	if val, ok := logData["ipAddress"]; ok {
		ip := net.ParseIP(val)
		if ip == nil {
			return log, fmt.Errorf("invalid IP address format for %s", val)
		}
		log.IPAddress = val
	}

	if val, ok := logData["errorCode"]; ok {
		log.ErrorCode = val
	}

	if val, ok := logData["errorMessage"]; ok {
		log.Message = val
	}

	return log, nil
}