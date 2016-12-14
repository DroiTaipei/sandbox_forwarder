package api

import (
	"fmt"
	"testing"
)

func TestGetFullURI(t *testing.T) {
	URI := getFullURI("123456", []byte("where={}&limit=2"))
	fmt.Println(URI)
	if URI != "123456?where={}&limit=2" {
		t.Fail()
	}
}
