package main

import (
	"bufio"
	"errors"
	"io"
	"strconv"
	"strings"
)

// token types for RESP Protocol
const (
	SimpleString = '+'
	ErrorString  = '-'
	Integer      = ':'
	BulkString   = '$'
	Array        = '*'
	Set          = "SET"
	Get          = "GET"
)

var err_invalid_token = errors.New("invalid token")

type RESPParser struct {
	reader *bufio.Reader
}

// NewRESPParser returns a new instance of RESPParser
func NewRESPParser(r io.Reader) *RESPParser {
	return &RESPParser{reader: bufio.NewReader(r)}
}

func (p *RESPParser) parse() (string, interface{}, error) {
	firstByte, err := p.reader.ReadByte()
	if err != nil {
		return "", nil, err
	}
	switch firstByte {
	case SimpleString:
		res, err := p.parseSimpleString()
		return "SimpleString", res, err
	case ErrorString:
		res, err := p.parseErrorString()
		return "ErrorString", res, err
	case Integer:
		res, err := p.parseInteger()
		return "Integer", res, err
	case BulkString:
		res, err := p.parseBulkString()
		return "BulkString", res, err
	case Array:
		cmd, res, err := p.parseArray()
		return cmd, res, err
	default:
		return "", nil, err_invalid_token
	}
}

// Helper function to read a line and trim CRLF
func (p *RESPParser) readLine() (string, error) {
	line, err := p.reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	return strings.TrimSuffix(line, "\r\n"), nil
}

// Helper function to read a line and trim CRLF
func (p *RESPParser) readWord() (string, error) {
	line, err := p.reader.ReadString(' ')
	if err != nil {
		return "", err
	}
	return strings.TrimSuffix(line, " "), nil
}

// Helper function to read a line and trim CRLF
func (p *RESPParser) readEOL() (string, error) {
	line, err := p.reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	return strings.TrimSuffix(line, "\n"), nil
}

func (p *RESPParser) parseSimpleString() (string, error) {
	line, err := p.readLine()
	if err != nil {
		return "", err
	}
	return line, nil
}

func (p *RESPParser) parseErrorString() (string, error) {
	line, err := p.readLine()
	if err != nil {
		return "", err
	}
	return "", errors.New("RESP Error: " + line)
}

func (p *RESPParser) parseInteger() (int64, error) {
	line, err := p.readLine()
	if err != nil {
		return 0, err
	}
	return strconv.ParseInt(line, 10, 64)
}

func (p *RESPParser) parseBulkString() (string, error) {
	line, err := p.readLine()
	if err != nil {
		return "", err
	}

	length, err := strconv.Atoi(line)
	if err != nil {
		return "", err
	}

	if length == -1 {
		return "", nil // Represents a NULL Bulk String in RESP
	}

	buf := make([]byte, length+2) // +2 to account for CRLF
	_, err = io.ReadFull(p.reader, buf)
	if err != nil {
		return "", err
	}
	return string(buf[:length]), nil
}

func (p *RESPParser) parseArray() (string, interface{}, error) {
	line, err := p.readLine()
	if err != nil {
		return "", nil, err
	}

	length, err := strconv.Atoi(line)
	if err != nil {
		return "", nil, err
	}

	if length == 0 {
		return "", nil, nil
	}

	result := make([]interface{}, length)
	for i := 0; i < length; i++ {
		_, elem, err := p.parse()
		if err != nil {
			return "", nil, err
		}
		result[i] = elem
	}

	// check for commands
	firstEle := strings.ToUpper(result[0].(string))
	switch firstEle {
	case "ECHO":
		res, err := p.handleEcho(result)
		return "ECHO", res, err
	case "PING":
		res, err := p.handlePing(result)
		return "PING", res, err
	case "SET":
		res, err := p.handleSet(result)
		return "SET", res, err
	case "GET":
		res, err := p.handleGet(result)
		return "GET", res, err
	}

	return "", result, nil
}

func (p *RESPParser) handlePing(result []interface{}) (interface{}, error) {
	if len(result) == 2 {
		if arg, ok := result[1].(string); ok {
			arg = strings.Join([]string{"+", arg, "\r\n"}, "")
			return arg, nil
		}
	}
	return "+PONG\r\n", nil
}

func (p *RESPParser) handleEcho(result []interface{}) (interface{}, error) {
	if len(result) != 2 {
		return nil, errors.New("ECHO requires exactly one argument " + result[0].(string))
	}
	if arg, ok := result[1].(string); ok {
		arg = strings.Join([]string{"+", arg, "\r\n"}, "")
		return arg, nil
	}
	return nil, errors.New("invalid argument for ECHO command" + result[1].(string))
}

func (p *RESPParser) handleSet(result []interface{}) (interface{}, error) {
	if len(result) <= 2 {
		return nil, errors.New("SET requires at least two arguments " + result[0].(string))
	}
	key := result[1].(string)
	value := result[2].(string)
	elements := make([]string, 0)
	elements = append(elements, key)
	elements = append(elements, value)
	return elements, nil

}

func (p *RESPParser) handleGet(result []interface{}) (interface{}, error) {
	if len(result) != 2 {
		return nil, errors.New("GET requires exactly one argument " + result[0].(string))
	}
	key := result[1].(string)
	elements := make([]string, 0)
	elements = append(elements, key)
	return elements, nil
}
