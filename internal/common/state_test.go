package common

import (
	"fmt"
	"testing"
)

func Test_State(t *testing.T) {
	state := &State{}
	// unhealthy by default
	if state.IsHealthy() {
		t.Error("should be unhealthy by default")
	}

	// make healthy
	state.SetHealthy()
	if !state.IsHealthy() {
		t.Error("should be healthy")
	}

	// add errors
	state.Error(fmt.Errorf("new error"))
	if state.errors != 1 {
		t.Error("error count should be 1")
	}
	if !state.IsHealthy() {
		t.Error("should be healthy after a single error")
	}
	// add more errors
	state.Error(fmt.Errorf("new error"))
	state.Error(fmt.Errorf("new error"))
	if state.errors != 3 {
		t.Error("should have 3 errors now")
	}
	if state.IsHealthy() {
		t.Error("should not be healthy at thsi point")
	}
}
