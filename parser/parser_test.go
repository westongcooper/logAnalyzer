package parser

import (
	"context"
	"fmt"
	"log"
	"math/rand/v2"
	"os"
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseLines(t *testing.T){
	t.Run("happy path - closed channel", func(t *testing.T){
		ctx := context.Background()
		testLines    := make(chan *string)

		concurrencty := 1
		logger       := log.New(os.Stderr, "", log.LstdFlags)
		expectedLog1      := fmt.Sprintf("log-Line-%d", rand.Uint64())
		expectedLog2      := fmt.Sprintf("log-Line-%d", rand.Uint64())

		callback1Called := false
		callback2Called := false
		testParsables := []Parsable{{
			regexp: regexp.MustCompile(regexp.QuoteMeta(expectedLog1)),
			callback: func(line *string) error {
				callback1Called = true
				assert.True(t, &expectedLog1 == line, "callback line pointers are not equal")
				return nil
			},
		},{
			regexp: regexp.MustCompile(regexp.QuoteMeta(expectedLog2)),
			callback: func(line *string) error {
				callback2Called = true
				assert.True(t, &expectedLog2 == line, "callback line pointers are not equal")
				return nil
			},
		}}

		go func() {
			testLines <- &expectedLog1
			testLines <- &expectedLog2
			close(testLines)
		}()

		err := ParseLines(ctx, testLines, concurrencty, logger, testParsables...)
		assert.NoError(t, err)

		assert.True(t, callback1Called, "ParseLines parsable callback was not exicuted")
		assert.True(t, callback2Called, "ParseLines parsable callback was not exicuted")
	})

	t.Run("happy path - closed context", func(t *testing.T){
		ctx, cancel  := context.WithCancel(context.Background())
		testLines    := make(chan *string)
		defer close(testLines)

		concurrencty := 1
		logger       := log.New(os.Stderr, "", log.LstdFlags)
		expectedLog1 := fmt.Sprintf("log-Line-%d", rand.Uint64())
		expectedLog2 := fmt.Sprintf("log-Line-%d", rand.Uint64())

		callback1Called := false
		callback2Called := false
		testParsables := []Parsable{{
			regexp: regexp.MustCompile(regexp.QuoteMeta(expectedLog1)),
			callback: func(line *string) error {
				callback1Called = true
				assert.True(t, &expectedLog1 == line, "callback line pointers are not equal")
				return nil
			},
		},{
			regexp: regexp.MustCompile(regexp.QuoteMeta(expectedLog2)),
			callback: func(line *string) error {
				callback2Called = true
				assert.True(t, &expectedLog2 == line, "callback line pointers are not equal")
				return nil
			},
		}}

		go func() {
			testLines <- &expectedLog1
			testLines <- &expectedLog2
			cancel()
		}()

		err := ParseLines(ctx, testLines, concurrencty, logger, testParsables...)
		assert.NoError(t, err)

		assert.True(t, callback1Called, "ParseLines parsable callback was not exicuted")
		assert.True(t, callback2Called, "ParseLines parsable callback was not exicuted")
	})	
}