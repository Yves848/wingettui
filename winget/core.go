package winget

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"
)

type WingetPackage struct {
	AvailableVersions string `json:"AvailableVersions"`
	Id                string `json:"Id"`
	Source            string `json:"Source"`
	InstalledVersion  string `json:"InstalledVersion"`
	IsUpdateAvailable bool   `json:"IsUpdateAvailable"`
	Name              string `json:"Name"`
}

type WingetPackages struct {
	Packages []WingetPackage `json:"packages"`
}

func Invoke(command string) (res []byte, err error) {
	cmds := strings.Split(command, " ")
	fmt.Println("Invoking command: ", command)
	cmd := exec.Command("pwsh.exe", cmds...)
	out, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Println("Error creating StdoutPipe: ", err)
		return
	}

	err = cmd.Start()
	if err != nil {
		fmt.Println("Error starting command: ", err)
		return
	}
	var result []byte = make([]byte, 0)
	scanner := bufio.NewScanner(out)
	for scanner.Scan() {
		line := scanner.Text()
		result = append(result, []byte(line)...)
	}

	err = cmd.Wait()
	if err != nil {
		fmt.Println("Error:", err)
	}
	return result, err
}

func Package2Json(out string) (items WingetPackages, err error) {
	// fmt.Println(out)
	js := "{\"packages\":" + out + "}"
	err = json.Unmarshal([]byte(js), &items)
	if err != nil {
		return items, err
	}
	return items, nil
}
