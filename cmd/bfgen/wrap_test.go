package main

import "testing"

func TestWrapAt(t *testing.T) {
	if wrapAt("123", 2) != "12\n3" {
		t.Fail()
	}
	if wrapAt("1234", 2) != "12\n34" {
		t.Fail()
	}
}
