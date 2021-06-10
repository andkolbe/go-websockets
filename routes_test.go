package main

import (
	"fmt"
	"testing"

	"github.com/bmizerany/pat"
)


func TestRoutes(t *testing.T) { 

	mux := routes()

	switch mux.(type) {
	case *pat.PatternServeMux:
		// do nothing. test passed
	default:
		t.Error(fmt.Sprintln("type is not *pat.PatternServeMux"))
	}
}