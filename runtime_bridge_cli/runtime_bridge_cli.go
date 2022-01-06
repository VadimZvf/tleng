package runtime_bridge_cli

import "fmt"

type Bridge struct {
}

func CreateBridge() Bridge {
	return Bridge{}
}

func (bridge *Bridge) Print(value string) {
	fmt.Println(value)
}
