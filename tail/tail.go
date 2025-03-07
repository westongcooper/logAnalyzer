package tail

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"time"
)

type tail struct {
	line     chan *string
	file     *os.File
	reader   *bufio.Reader
	logger   logger
}

func (tail *tail) getFileSize() (int64, error) {
	info, err := tail.file.Stat()
	if err != nil {
		return 0, fmt.Errorf("failed to get file stat: %s", err.Error())
	}

	return info.Size(), nil
}

func newTail(filePath string, tailEndOfFile bool, logger logger) (*tail, error) {
	openFile, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("os.Open failed: %s", err)
	}

	if tailEndOfFile { // tail end of file by default
		_, err := openFile.Seek(0, io.SeekEnd)
		logger.Printf("Moving to end of log file: %s", filePath)
		if err != nil {
			return nil, fmt.Errorf("failed to move to end of file: %s", err)
		}
	}

	return &tail{
		file:   openFile,
		reader: bufio.NewReader(openFile),
		line:   make(chan *string),
		logger: logger,
	}, nil
}

// ======== This is a recursive function,  If ReadLine() returns `isPrefix == true`, then the line exceeded the readers buffer size and we should try again until `isPrefix` returns false.
func readLine(reader *bufio.Reader, _lineBuilder ...*bytes.Buffer) (string, error) {
	var lineBuilder *bytes.Buffer
	if len(_lineBuilder) == 0 {
		lineBuilder = &bytes.Buffer{}
	} else {
		lineBuilder = _lineBuilder[0]
	}

	line, isPrefix, err := reader.ReadLine();

	if err != nil {
		if err == io.EOF { // we reached the end of the file, there was nothing else to read.
			return lineBuilder.String(), err // return EOF error with built string
		}
		return lineBuilder.String(), fmt.Errorf("failed to read line: %s", err.Error())
	}

	if isPrefix { // there's more to collect for this line
		lineBuilder.Write(line)

		return readLine(reader, lineBuilder)
	} else if lineBuilder.Len() == 0 { // this is end of line, and there was no recursion
		return string(line), nil // return the line without building a new string.
	} else { // end of the line, add to previous data	
		lineBuilder.Write(line)
	}

	return lineBuilder.String(), nil
}

func TailFile(ctx context.Context, filePath string, tailEndOfFile bool, logger logger) (chan *string, error) {
	tail, err := newTail(filePath, tailEndOfFile, logger)
	if err != nil {
		return nil, err
	}

	logger.Printf("tailing file: %s", tail.file.Name())

	oldSize, err := tail.getFileSize()
	if err != nil {
		close(tail.line)
		return nil, fmt.Errorf("failed to get file size: %s", err.Error())
	}

	go func() {
		defer close(tail.line)

		tailFile:for {
			readLineInFile: for {
				line, err := readLine(tail.reader)

				if line != "" {
					tail.line <- &line
				}

				if err != nil {
					if err == io.EOF {
						break readLineInFile
					}
					logger.Panicf("failed to read line: %s", err.Error())
				}

				if ctx.Err() != nil {
					break tailFile
				}
			}

			waitForChanges: for {
				if err := sleepWithContext(ctx, 100*time.Millisecond); err != nil {
					logger.Printf("stopped waiting for log updates: %s", err.Error())
					break tailFile
				}

				newSize, err := tail.getFileSize()
				if err != nil {
					logger.Panicf("failed to get file size: %s", err.Error())
					break tailFile
				}

				if newSize != oldSize {
					if newSize < oldSize { // file may be truncated
						tail.file.Seek(0, io.SeekStart)// reset position to beginning
						tail.reader = bufio.NewReader(tail.file)
					}
					
					oldSize = newSize

					break waitForChanges // new data to check
				}	
			}
		}
	}()

	return tail.line, nil
}

func sleepWithContext(ctx context.Context, d time.Duration) error {
    timer := time.NewTimer(d)
	defer timer.Stop()

    select {
    case <-timer.C:	
		return nil
    case <-ctx.Done():
        return ctx.Err()
	}
}

type logger interface {
	Panicf(format string, v ...interface{})
	Printf(format string, v ...interface{})
	Println(v ...interface{})
}
