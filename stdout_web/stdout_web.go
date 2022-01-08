package stdout_web

import (
	"syscall/js"
)

type JSPrintWord func(value string)
type JSPrintErrorWord func(value string)

type StdoutWeb struct {
	JSPrintWord      JSPrintWord
	JSPrintErrorWord JSPrintErrorWord
}

func CreateStdoutWeb() StdoutWeb {
	WindowJS := js.Global().Get("window")

	return StdoutWeb{
		JSPrintWord: func(value string) {
			WindowJS.Call("TlengPrintWord", value)
		},
		JSPrintErrorWord: func(value string) {
			WindowJS.Call("TlengPrintErrorWord", value)
		},
	}
}

func (stdout *StdoutWeb) Print(line string) {
	stdout.JSPrintWord(line)
}

func (stdout *StdoutWeb) PrintError(symbol string) {
	stdout.JSPrintErrorWord(symbol)
}
