package cliui

import (
	"fmt"
	"os"
	"strings"
	"time"
)

type Spinner struct {
	chars   []string
	msg     string
	stop    chan bool
	running bool
	pos     int
}

func NewSpinner(msg string) *Spinner {
	return &Spinner{
		chars: []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"},
		msg:   msg,
		stop:  make(chan bool),
	}
}

func (s *Spinner) Start() {
	if s.running {
		return
	}
	s.running = true
	s.pos = 0

	go func() {
		for {
			select {
			case <-s.stop:
				return
			default:
				fmt.Fprintf(os.Stdout, "\r  %s %s ", Cyan, s.chars[s.pos%len(s.chars)])
				fmt.Fprint(os.Stdout, s.msg)
				os.Stdout.Sync()
				s.pos++
				time.Sleep(80 * time.Millisecond)
			}
		}
	}()
}

func (s *Spinner) Stop() {
	if !s.running {
		return
	}
	s.stop <- true
	s.running = false
	fmt.Fprint(os.Stdout, "\r"+strings.Repeat(" ", len(s.msg)+6)+"\r")
	fmt.Fprint(os.Stdout, "  ✅ ")
	fmt.Fprintln(os.Stdout, s.msg)
}

func (s *Spinner) StopError(err error) {
	if !s.running {
		return
	}
	s.stop <- true
	s.running = false
	fmt.Fprint(os.Stdout, "\r"+strings.Repeat(" ", len(s.msg)+6)+"\r")
	fmt.Fprint(os.Stdout, "  ❌ ")
	fmt.Fprintln(os.Stdout, s.msg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "     %s\n", err.Error())
	}
}
