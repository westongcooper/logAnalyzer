package parser

import (
	"context"
	"regexp"

	"golang.org/x/sync/errgroup"
)

type Parsable struct {
	regexp   *regexp.Regexp
	callback func(line *string) error
}

func NewParsable(regex string, callback func(*string)error) Parsable {
	return Parsable{
		regexp: regexp.MustCompile(regex),
		callback:  callback,
	}
}

type logger interface {
	Printf(format string, v ...interface{})
	Println(v ...interface{})
}

type parser struct {
	lines       chan *string
	concurrency int
	logger      logger
	parsables 	[]Parsable
}

func ParseLines(ctx context.Context, lines chan *string, concurrency int, logger logger, parsables ...Parsable) error {
	return newParser(lines, concurrency, logger, parsables...).
		parseLines(ctx)
}

func newParser(lines chan *string, concurrency int, logger logger, parsables ...Parsable) parser {
	return parser{
		lines: lines,
		concurrency: concurrency,
		parsables: parsables,
		logger: logger,
	}
}

func (parser parser) parseLines(ctx context.Context) error {
	eg := errgroup.Group{}
	for range parser.concurrency {
		func () {
			eg.Go(func() error {
				for {
					select {
					case line, ok := <-parser.lines:
						if !ok {
							return nil
						}
						if err := parser.parseLine(line); err != nil {
							return err
						}
					case <- ctx.Done():
						return nil
					}
				}
			})
		}()
	}

	return eg.Wait()
}

func (parser parser) parseLine(line *string) error {
	for _, parsable := range parser.parsables {
		if parsable.regexp.MatchString(*line) {
			if err := parsable.callback(line); err != nil {
				return err
			}
		}
	}
	return nil
}
