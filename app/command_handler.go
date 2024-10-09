package main

import (
	"errors"
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
	args, ok := req.Args.([]interface{})
	if !ok {
		return "", errors.New("invalid args")
	}
	switch args[0] {
	case Set:
		return handleSet(c, args)
	case Get:
		return handleGet(c, args)
	case Ping:
		return handlePing(args)
	case Echo:
		return handleEcho(args)
	default:
		return "", fmt.Errorf("unknown command: %s", args[0])
	}
}

func handleSet(c *CacheWithTTL, args []interface{}) (string, error) {
	if len(args) < 3 {
		return "", fmt.Errorf("invalid number of arguments for SET")
	}

	key := args[1].(string)
	value := args[2].(string)

	// Check if TTL is provided
	if len(args) == 5 && strings.EqualFold(args[3].(string), "px") {
		ttl, err := strconv.Atoi(args[4].(string))
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
	if len(args) < 2 {
		return "", fmt.Errorf("invalid number of arguments for GET")
	}

	key := args[1].(string)
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
	if len(args) > 1 {
		return fmt.Sprintf("+%s\r\n", args[1]), nil
	}
	return "+PONG\r\n", nil
}

func handleEcho(args []interface{}) (string, error) {
	if len(args) < 2 {
		return "", fmt.Errorf("ECHO requires an argument")
	}
	return fmt.Sprintf("$%d\r\n%s\r\n", len(args[1].(string)), args[1].(string)), nil
}
