package main

import (
	"io"
	"strings"
	"testing"
)

// An option's key is a to d: the REPL's Topic page takes the other
// letters as its own keys.
func TestNewsOpenOptionKeys(t *testing.T) {
	for _, opts := range [][]string{
		{"--option", "f=mở session"},
		{"--option", "1=một"},
		{"--option", "ab=hai chữ"},
		{"--option", "a=1", "--option", "b=2", "--option", "c=3", "--option", "d=4", "--option", "e=5"},
	} {
		args := append([]string{"p", "--title", "t"}, opts...)
		err := cmdNewsOpen(args, io.Discard, io.Discard)
		if err == nil || !strings.Contains(err.Error(), "option") {
			t.Errorf("%q: err = %v, want an option error", opts, err)
		}
	}
}
