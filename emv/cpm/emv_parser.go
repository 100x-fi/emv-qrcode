package cpm

import (
	"errors"
	"fmt"
	"strconv"
	"unicode/utf8"
)

// ParserError ...
type ParserError struct {
	Func string
	Err  error
}

func (e *ParserError) Error() string {
	return "parser." + e.Func + ": " + e.Err.Error()
}

func notCallError(fn string) *ParserError {
	return &ParserError{
		Func: fn,
		Err:  errors.New("not call Next()"),
	}
}

func outOfRangeError(fn string, current, max, start, end int64) *ParserError {
	return &ParserError{
		Func: fn,
		Err:  fmt.Errorf("bounds out of range. current: %d, max: %d, start: %d, end: %d", current, max, start, end),
	}
}

func syntaxError(fn, str string) *ParserError {
	return &ParserError{
		Func: fn,
		Err:  errors.New("parsing " + strconv.Quote(str) + ": " + strconv.ErrSyntax.Error()),
	}
}

func idRangeError(fn string, id ID) *ParserError {
	return &ParserError{
		Func: fn,
		Err:  errors.New("id range invalid. id: " + id.String()),
	}
}

// const ...
const (
	IDWordCount          = 2
	ValueLengthWordCount = 2
)

// Parser ...
type Parser struct {
	current int64
	max     int64
	source  []rune
	err     error
}

// NewParser ...
func NewParser(payload string) *Parser {
	return &Parser{
		current: -1,
		max:     int64(utf8.RuneCountInString(payload)),
		source:  []rune([]rune(payload)),
		err:     nil,
	}
}

// Next ...
func (p *Parser) Next(idWordCount int) bool {
	if p.err != nil {
		return false
	}
	if p.current < 0 {
		p.current = 0
	} else {
		valueLength := p.ValueLength(idWordCount)
		if p.err != nil {
			return false
		}
		p.current += valueLength + int64(idWordCount) + ValueLengthWordCount
	}
	if p.current >= p.max {
		return false
	}
	return true
}

// ID ...
func (p *Parser) ID(idWordCount int) ID {
	const fnID = "ID"
	start := p.current
	end := start + int64(idWordCount)
	if p.current < 0 {
		p.err = notCallError(fnID)
		return ID("")
	}
	if p.max < end {
		p.err = outOfRangeError(fnID, p.current, p.max, start, end)
		return ID("")
	}
	id := ID(string(p.source[start:end]))
	return id
}

// ValueLength ...
func (p *Parser) ValueLength(idWordCount int) int64 {
	const fnValueLength = "ValueLength"
	start := p.current + int64(idWordCount)
	end := start + ValueLengthWordCount
	if p.current < 0 {
		p.err = notCallError(fnValueLength)
		return 0
	}
	if p.max < end {
		p.err = outOfRangeError(fnValueLength, p.current, p.max, start, end)
		return 0
	}
	strValueLength := string(p.source[start:end])
	len, err := strconv.ParseInt(strValueLength, 16, 64)
	if err != nil {
		p.err = syntaxError(fnValueLength, strValueLength)
		return 0
	}
	return len * 2 // len divided by 2 in generate payload
}

// Value ...
func (p *Parser) Value(idWordCount int) string {
	const fnValue = "Value"
	start := p.current + int64(idWordCount) + ValueLengthWordCount
	end := start + p.ValueLength(idWordCount)
	if p.current < 0 {
		p.err = notCallError(fnValue)
		return ""
	}
	if p.max < end {
		p.err = outOfRangeError(fnValue, p.current, p.max, start, end)
		return ""
	}
	return string(p.source[start:end])
}

// Err ...
func (p *Parser) Err() error {
	return p.err
}
