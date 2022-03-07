package main

import (
	"bufio"
	"bytes"
	"errors"
	"os"

	"github.com/mastercactapus/embedded/term"
	"github.com/mastercactapus/embedded/term/ansi"
)

func main() {
	r := bufio.NewReader(os.Stdin)

	w := ansi.NewPrinter(term.NewThrottleWriter(os.Stdout, 1000))

	var b bytes.Buffer
	b.ReadRune()

	p := term.NewPrompt(w, "test> ")
	p.Draw()
	for {
		r, _, err := r.ReadRune()
		if err != nil {
			w.Printf("ERROR: read: %v\n", err)
			break
		}

		cmd, err := p.NextCommand(r)
		if errors.Is(err, term.ErrInterrupt) {
			p.Draw()
			continue
		}
		if err != nil {
			w.Printf("ERROR: parser: %v\n", err)
			p.Draw()
			continue
		}

		if cmd == nil {
			continue
		}

		w.Printf("COMMAND: %#v\n", cmd)
		if cmd.Args[0] == "exit" {
			break
		}

		p.Draw()
	}
}
