package util

import (
	"fmt"
	"net"
	"strings"
	"testing"
)

func TestFindFreePort(t *testing.T) {
	port, err := FindFreePort()
	if err != nil {
		t.Fatal(err)
	}

	l, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		if strings.Contains(err.Error(), "address already in use") {
			t.Fatal("unable to listen on a port returnen by FindFreePort: ", err)
		} else {
			t.Fatal(err)
		}

	}
	l.Close()
}
