package main

import (
	"strings"
	"testing"
)

func TestParseSetWithPX(t *testing.T) {
	input := "*5\r\n$3\r\nSET\r\n$4\r\npear\r\n$6\r\norange\r\n$2\r\npx\r\n$3\r\n100\r\n" // SET pear pear PX 100
	reader := strings.NewReader(input)
	parser := NewRESPParser(reader)

	req, err := parser.Parse()
	if err != nil {
		t.Fatalf("expected no error, got %v", err.Error())
	}

	expectedArgs := []string{"SET", "pear", "orange", "px", "100"}
	args, ok := req.Args.([]interface{})
	if !ok {
		t.Fatalf("type mismatch, expected array of arguments, got %v", req.Args)
	}
	if len(args) != len(expectedArgs) {
		t.Fatalf("expected %d arguments, got %d", len(expectedArgs), len(args))
	}

	for i, arg := range args {
		if arg != expectedArgs[i] {
			t.Errorf("expected arg %d to be %q, got %q", i, expectedArgs[i], arg)
		}
	}
}
