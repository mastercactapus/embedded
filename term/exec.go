package term

import (
	"context"
	"errors"
	"fmt"

	"github.com/mastercactapus/embedded/term/ansi"
)

func (sh *Shell) Exec(ctx context.Context) error {
	rd := sh.r
	var history []string
	var histIndex int

	var input []rune
	var pos int
	var lastRune rune
	for sh.err == nil {
		sh.p.Println()
		sh.prompt(string(input))
	readInput:
		for {
			r, _, err := rd.ReadRune()
			if err != nil {
				return err
			}
			if r == '\n' && lastRune == '\r' {
				lastRune = r
				continue
			}

			switch r {
			case '\x07', '\t': // ignore
				continue
			case '\x08', '\x7f': // backspace, remove last char
				if pos > 0 {
					input = append(input[:pos-1], input[pos:]...)
					pos--
					sh.p.CurLt(1)
					sh.p.EraseLine(ansi.CurToEnd)
					sh.w.Flush()
				}
				continue
			case '\r', '\n':
				// TODO: consume all newlines
				sh.p.Println()
				break readInput
			case '\x1b':
				esc, err := ansi.ParseEscapeSequence(rd)
				if err != nil {
					return err
				}
				if esc == nil {
					continue
				}
				switch esc.Code {
				case 'A':
					if histIndex == 0 {
						continue
					}
					histIndex--
					input = []rune(history[histIndex])
					pos = len(input)
					sh.prompt(string(input))
				case 'B':
					if histIndex == len(history) {
						continue
					}
					histIndex++
					if histIndex == len(history) {
						input = nil
					} else {
						input = []rune(history[histIndex])
					}
					pos = len(input)
					sh.prompt(string(input))
				case 'D':
					n := esc.Args[0]
					if n == 0 {
						n = 1
					}
					if n > pos {
						n = pos
					}
					if n > 0 {
						sh.p.CurLt(n)
						pos -= n
					}
				case 'C':
					n := esc.Args[0]
					if n == 0 {
						n = 1
					}
					if pos+n > len(input) {
						n = len(input) - pos
					}
					if n > 0 {
						sh.p.CurRt(n)
						pos += n
					}
				}

				sh.w.Flush()
				continue
			default:
				if r < ' ' || r > '~' {
					continue
				}
				sh.p.Print(string(r))
				if pos == len(input) {
					input = append(input, r)
				} else {
					sh.p.SaveCursor()
					sh.p.Print(string(input[pos:]))
					sh.p.RestoreCursor()
					input = append(input[:pos], append([]rune{r}, input[pos:]...)...)
				}
				pos++
			}
			sh.w.Flush()
		}
		sh.w.Flush()

		// remove identical entries
		for i, entry := range history {
			if entry == string(input) {
				history = append(history[:i], history[i+1:]...)
				break
			}
		}
		history = append(history, string(input))
		histIndex = len(history)

		cmdLine, err := ParseCmdLine(string(input))
		if err != nil {
			sh.p.Println(err)
			continue
		}
		input = input[:0]
		pos = 0
		if len(cmdLine.Args) == 0 {
			continue
		}
		if err != nil {
			sh.lastCmdErr = fmt.Errorf("parse args: %w", err)
			continue
		}

		cmd := sh.cmds[cmdLine.Args[0]]
		if cmd == nil {
			sh.p.Println("Unknown command: '" + cmdLine.Args[0] + "' try 'help'.")
			continue
		}

		cmdCtx := context.WithValue(ctx, ctxKeyCmd, &cmdContext{
			sh:      cmd.sh,
			CmdLine: cmdLine,
			desc:    cmd.Desc,
			env:     sh.env,
			fs:      NewFlagSet(cmdLine, sh.env.Get),
		})

		if cmd.Init != nil {
			sh.lastCmdErr = cmd.Init(cmdCtx, cmd.Exec)
		} else {
			sh.lastCmdErr = cmd.Exec(cmdCtx)
		}

		var exit exitErr
		if errors.As(sh.lastCmdErr, &exit) {
			sh.lastCmdErr = nil
			return exit.error
		}

		var usage usageErr
		if errors.As(sh.lastCmdErr, &usage) {
			if usage.err == nil {
				sh.lastCmdErr = nil
			} else {
				sh.p.Fg(ansi.Red)
				sh.p.Println(usage.err.Error())
				sh.p.Reset()
			}
			if cmd.Desc != "" {
				sh.p.Println(cmd.Desc)
				sh.p.Println()
			}
			usage.PrintUsage(sh.p)
		} else if sh.lastCmdErr != nil {
			sh.p.Fg(ansi.Red)
			sh.p.Println(sh.lastCmdErr)
		}
	}

	return sh.err
}
