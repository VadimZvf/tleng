package stdout

import (
	"github.com/fatih/color"
)

type Stdout struct {
	color color.Attribute
}

func CreateStdout() Stdout {
	return Stdout{}
}

func (stdout *Stdout) PrintLine(line string) {
	color.New(stdout.color).Printf("\n")
	color.New(stdout.color).Printf(line)
	color.New(stdout.color).Printf("\n")
}

func (stdout *Stdout) PrintSymbol(symbol string) {
	color.New(stdout.color).Printf(symbol)
}

func (stdout *Stdout) SetErrorColor() {
	stdout.color = color.FgHiRed
}

func (stdout *Stdout) SetDefaultColor() {
	stdout.color = color.Reset
}
