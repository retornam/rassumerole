package main

import "testing"

func TestFailedAssumeProfile(t *testing.T) {
	role := "example-profile"
	_, err := assumeProfile(role)
	if err != nil {
		t.Log("Assume profile failed: ")
	}
}
