package stdout

import (
	"github.com/fatih/color"
)

type Stdout struct{}

func CreateStdout() Stdout {
	return Stdout{}
}

func (stdout *Stdout) Print(line string) {
	color.New(color.Reset).Printf(line)
}

func (stdout *Stdout) PrintError(symbol string) {
	color.New(color.FgHiRed).Printf(symbol)
}
