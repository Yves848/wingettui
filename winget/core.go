package winget

import (
	"fmt"
	"os/exec"
	"strings"
)

func Invoke(command string) (res []byte, err error) {
	cmds := strings.Split(command, " ")
	fmt.Println("Invoking command: ", command)
	cmd := exec.Command("winget", cmds...)
	out, err := cmd.Output()
	return out, err
}
