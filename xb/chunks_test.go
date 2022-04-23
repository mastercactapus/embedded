package xb

import (
	"bufio"
	"bytes"
	"testing"
)

func TestReadWrite(t *testing.T) {
	var buf bytes.Buffer
	var dat [512]byte
	err := WriteChunk(&buf, 'R', (&Request{Cmd: spiWrite, Data: dat[:]}).encode())
	if err != nil {
		t.Fatal(err)
	}

	tc, data, err := ReadChunk(bufio.NewReader(&buf))
	if err != nil {
		t.Fatal(err)
	}
	if tc != 'R' {
		t.Fatal("wrong type code")
	}
	var req Request
	req.decode(data)
	if req.Cmd != spiWrite {
		t.Fatal("wrong cmd")
	}
	if !bytes.Equal(req.Data, dat[:]) {
		t.Fatal("wrong data")
	}
}
