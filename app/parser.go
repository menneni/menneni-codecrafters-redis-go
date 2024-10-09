package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"
)

var (
	ErrInvalidCommand = errors.New("invalid RESP command")
)

type RESPRequest struct {
	Cmd  string
	Args []interface{}
}

type RESPParser struct {
	reader *bufio.Reader
}

// NewRESPParser creates a new RESP parser
func NewRESPParser(r io.Reader) *RESPParser {
	return &RESPParser{
		reader: bufio.NewReader(r),
	}
}

// Parse parses a RESP request from the connection
func (p *RESPParser) Parse() (*RESPRequest, error) {
	line, err := p.reader.ReadString('\n')
	if err != nil {
		return nil, err
	}
	line = strings.TrimSpace(line)

	if len(line) < 1 {
		return nil, ErrInvalidCommand
	}

	switch line[0] {
	case '*': // Array (Command and Arguments)
		return p.parseArray()
	case '+': // Simple String
		return p.parseSimpleString(line)
	case '$': // Bulk String
		return p.parseBulkString(line)
	case ':': // Integer
		return p.parseInteger(line)
	case '-': // Error
		return nil, fmt.Errorf("RESP error: %s", line[1:])
	default:
		return nil, ErrInvalidCommand
	}
}

func (p *RESPParser) parseArray() (*RESPRequest, error) {
	sizeStr, err := p.reader.ReadString('\n')
	if err != nil {
		return nil, err
	}

	size, err := strconv.Atoi(strings.TrimSpace(sizeStr))
	if err != nil {
		return nil, err
	}

	req := &RESPRequest{}
	for i := 0; i < size; i++ {
		arg, err := p.reader.ReadString('\n')
		if err != nil {
			return nil, err
		}
		arg = strings.TrimSpace(arg)

		if i == 0 {
			req.Cmd = strings.ToUpper(arg)
		} else {
			req.Args = append(req.Args, arg)
		}
	}

	return req, nil
}

func (p *RESPParser) parseSimpleString(line string) (*RESPRequest, error) {
	return &RESPRequest{Cmd: line[1:], Args: nil}, nil
}

func (p *RESPParser) parseBulkString(line string) (*RESPRequest, error) {
	length, err := strconv.Atoi(line[1:])
	if err != nil {
		return nil, err
	}

	bulkString := make([]byte, length+2) // +2 for \r\n
	_, err = io.ReadFull(p.reader, bulkString)
	if err != nil {
		return nil, err
	}

	value := string(bulkString[:length])
	return &RESPRequest{Cmd: "ECHO", Args: []interface{}{value}}, nil
}

func (p *RESPParser) parseInteger(line string) (*RESPRequest, error) {
	intVal, err := strconv.Atoi(line[1:])
	if err != nil {
		return nil, err
	}

	return &RESPRequest{Cmd: "INTEGER", Args: []interface{}{intVal}}, nil
}
