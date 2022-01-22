package stdout_mock

type Stdout struct{}

func CreateStdout() Stdout {
	return Stdout{}
}

func (stdout *Stdout) Print(line string) {
}

func (stdout *Stdout) PrintError(symbol string) {
}
