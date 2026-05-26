package cliui

import (
	"fmt"
	"os"
	"strings"
	"time"
)

type ProgressBar struct {
	total   int
	current int
	msg     string
	width   int
	stop    chan bool
	running bool
}

func NewProgressBar(total int, msg string) *ProgressBar {
	return &ProgressBar{
		total: total,
		msg:   msg,
		width: 30,
		stop:  make(chan bool),
	}
}

func (p *ProgressBar) Start() {
	if p.running {
		return
	}
	p.running = true

	go func() {
		for {
			select {
			case <-p.stop:
				return
			default:
				p.draw()
				time.Sleep(100 * time.Millisecond)
			}
		}
	}()
}

func (p *ProgressBar) Inc() {
	if p.current < p.total {
		p.current++
	}
}

func (p *ProgressBar) Stop() {
	if !p.running {
		return
	}
	p.stop <- true
	p.running = false
	p.current = p.total
	p.draw()
	fmt.Fprintln(os.Stdout)
}

func (p *ProgressBar) draw() {
	ratio := float64(p.current) / float64(p.total)
	filled := int(ratio * float64(p.width))

	bar := Green + strings.Repeat("█", filled) + Dim + strings.Repeat("░", p.width-filled) + Reset
	percent := fmt.Sprintf("%3.0f%%", ratio*100)

	fmt.Fprintf(os.Stdout, "\r  %s %s %s", bar, percent, p.msg)
	os.Stdout.Sync()
}
