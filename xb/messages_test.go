package xb

import (
	"bytes"
	"testing"

	"github.com/francoispqt/gojay"
)

func TestRequest(t *testing.T) {
	var buf bytes.Buffer
	enc := gojay.NewEncoder(&buf)
	foo := &Request{Cmd: ignore}
	if err := enc.Encode(foo); err != nil {
		t.Fatal(err)
	}
	if err := enc.Encode(foo); err != nil {
		t.Fatal(err)
	}

	dec := gojay.NewDecoder(&buf)
	if err := dec.Decode(foo); err != nil {
		t.Fatal(err)
	}
	if err := dec.Decode(foo); err != nil {
		t.Fatal(err)
	}
}

func TestResponse(t *testing.T) {
	var buf bytes.Buffer
	enc := gojay.NewEncoder(&buf)
	foo := &Response{Err: "foo"}
	if err := enc.Encode(foo); err != nil {
		t.Fatal(err)
	}
	if err := enc.Encode(foo); err != nil {
		t.Fatal(err)
	}

	dec := gojay.NewDecoder(&buf)
	if err := dec.Decode(foo); err != nil {
		t.Fatal(err)
	}
	if err := dec.Decode(foo); err != nil {
		t.Fatal(err)
	}
}
