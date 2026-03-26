package main

import "testing"

func TestHelloName(t *testing.T) {
	got := HelloName("Chris", "")
	want := "Hello, Chris"

	if got != want {
		t.Errorf("got %q want %q", got, want)
	}

}
