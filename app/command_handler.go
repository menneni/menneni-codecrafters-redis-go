package main

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

const (
	Set  = "SET"
	Get  = "GET"
	Ping = "PING"
	Echo = "ECHO"
)

// HandleRequest processes incoming RESP commands
func HandleRequest(c *CacheWithTTL, req *RESPRequest) (string, error) {
	switch req.Cmd {
	case Set:
		return handleSet(c, req.Args)
	case Get:
		return handleGet(c, req.Args)
	case Ping:
		return handlePing(req.Args)
	case Echo:
		return handleEcho(req.Args)
	default:
		return "", fmt.Errorf("unknown command: %s", req.Cmd)
	}
}

func handleSet(c *CacheWithTTL, args []interface{}) (string, error) {
	if len(args) < 2 {
		return "", fmt.Errorf("invalid number of arguments for SET")
	}

	key := args[0].(string)
	value := args[1].(string)

	// Check if TTL is provided
	if len(args) == 4 && strings.EqualFold(args[2].(string), "px") {
		ttl, err := strconv.Atoi(args[3].(string))
		if err != nil {
			return "", fmt.Errorf("invalid TTL value")
		}
		c.Set(key, value, time.Duration(ttl)*time.Millisecond)
	} else {
		c.Set(key, value, 0)
	}

	return "+OK\r\n", nil
}

func handleGet(c *CacheWithTTL, args []interface{}) (string, error) {
	if len(args) < 1 {
		return "", fmt.Errorf("invalid number of arguments for GET")
	}

	key := args[0].(string)
	value, found := c.Get(key)
	if !found {
		return "$-1\r\n", nil
	}

	valStr, ok := value.(string)
	if !ok {
		return "", fmt.Errorf("invalid value type in cache")
	}

	return fmt.Sprintf("$%d\r\n%s\r\n", len(valStr), valStr), nil
}

func handlePing(args []interface{}) (string, error) {
	if len(args) > 0 {
		return fmt.Sprintf("+%s\r\n", args[0].(string)), nil
	}
	return "+PONG\r\n", nil
}

func handleEcho(args []interface{}) (string, error) {
	if len(args) < 1 {
		return "", fmt.Errorf("ECHO requires an argument")
	}
	return fmt.Sprintf("$%d\r\n%s\r\n", len(args[0].(string)), args[0].(string)), nil
}
