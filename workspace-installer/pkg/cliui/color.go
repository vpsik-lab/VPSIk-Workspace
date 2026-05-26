package cliui

import (
	"fmt"
	"strings"
)

const (
	Reset     = "\033[0m"
	Bold      = "\033[1m"
	Dim       = "\033[2m"
	Italic    = "\033[3m"
	Underline = "\033[4m"

	Red    = "\033[31m"
	Green  = "\033[32m"
	Yellow = "\033[33m"
	Blue   = "\033[34m"
	Purple = "\033[35m"
	Cyan   = "\033[36m"
	White  = "\033[37m"

	BgRed    = "\033[41m"
	BgGreen  = "\033[42m"
	BgYellow = "\033[43m"
	BgBlue   = "\033[44m"

	RedFg    = "\033[91m"
	GreenFg  = "\033[92m"
	YellowFg = "\033[93m"
	BlueFg   = "\033[94m"
	PurpleFg = "\033[95m"
	CyanFg   = "\033[96m"
)

type Style func(format string, a ...any) string

func Sprintf(code string) Style {
	return func(format string, a ...any) string {
		return code + fmt.Sprintf(format, a...) + Reset
	}
}

var (
	Success   = Sprintf(Green)
	Error     = Sprintf(Red)
	Warning   = Sprintf(Yellow)
	Info      = Sprintf(Cyan)
	Highlight = Sprintf(Purple)
	DimText   = Sprintf(Dim)
	BoldText  = Sprintf(Bold)
)

func Header(text string) string {
	line := strings.Repeat("─", 50)
	return fmt.Sprintf("\n%s\n  %s\n%s\n", BoldText(text), DimText(line))
}

func Summary(items ...string) string {
	var b strings.Builder
	for _, item := range items {
		b.WriteString("  " + item + "\n")
	}
	return b.String()
}

func Label(key, value string) string {
	return fmt.Sprintf("  %s: %s", BoldText(key), value)
}
