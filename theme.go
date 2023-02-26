package pretty

import (
	"github.com/fatih/color"
	"time"
)

var DefaultTheme = &Theme{
	Nil:        color.New(color.FgGreen),
	Float:      color.New(color.FgMagenta),
	Integer:    color.New(color.FgYellow),
	String:     color.New(color.FgCyan),
	Bool:       color.New(color.FgRed),
	Time:       color.New(color.FgBlue),
	TimeLayout: time.RFC3339,
}

type Theme struct {
	Nil        *color.Color
	Float      *color.Color
	Integer    *color.Color
	String     *color.Color
	Bool       *color.Color
	Time       *color.Color
	TimeLayout string
}
