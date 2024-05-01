package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/Yves848/wingettui/winget"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
	tsize "github.com/kopoli/go-terminal-size"
)

func PadOrTruncateString(input string, width int) string {
	if len(input) >= width {
		// If the input string is longer than or equal to the width,
		// truncate it and append an ellipsis.
		return input[:width-3] + "…"
	} else {
		// If the input string is shorter than the width,
		// pad it with spaces.
		return input + strings.Repeat(" ", width-len(input))
	}
}

func main() {
	var s tsize.Size
	s, err2 := tsize.GetSize()
	if err2 == nil {
		fmt.Println("Current Size is ", s.Width, "x", s.Height)
	}
	colWidths := make(map[string]int)
	colWidths["Name"] = 40
	colWidths["Id"] = 40
	colWidths["Version"] = 10
	listCmd := flag.NewFlagSet("list", flag.ExitOnError)
	listSource := listCmd.String("source", "", "source of the package")

	searchCmd := flag.NewFlagSet("search", flag.ExitOnError)
	searchQuery := searchCmd.String("query", "", "query to search")

	if len(os.Args) < 2 {
		fmt.Println("expected 'list / search / update' subcommand")
		os.Exit(1)
	}

	var params []string = make([]string, 1)

	switch os.Args[1] {
	case "list":
		listCmd.Parse(os.Args[2:])
		fmt.Println("listsource :  " + *listSource)

		params[0] = "-command Get-WinGetPackage | where-object { $_.Source -eq \"winget\"} | ConvertTo-Json -Depth 1 -AsArray -WarningAction SilentlyContinue"
		if *listSource != "" {

		}

	case "search":
		searchCmd.Parse(os.Args[2:])
		fmt.Println("searchQuery :  " + *searchQuery)
		params = make([]string, 1)
		// params[0] = "Search-WGPackage \"##QUERY##\" -quiet $true | out-json"
		params[0] = "-command Find-WingetPackage -query ##QUERY## -source winget | ConvertTo-Json -Depth 1 -AsArray -WarningAction SilentlyContinue"
		if *searchQuery != "" {
			params[0] = strings.Replace(params[0], "##QUERY##", *searchQuery, -1)
		}

	}

	if params[0] != "" {
		out, err := winget.Invoke(strings.Join(params, " "))
		if err != nil {
			fmt.Println(errors.New("error invoking winget"))
		}
		// fmt.Println(string(out))
		items, err := winget.Package2Json(string(out))
		if err != nil {
			fmt.Println(errors.New("error parsing json"))
			os.Exit(2)
		}
		rows := make([][]string, 0)
		for _, item := range items.Packages {
			versions := strings.Split(item.AvailableVersions, " ")
			// fmt.Printf("%s %s %s %s\n", item.Id, item.Name, item.InstalledVersion, versions[0])
			rows = append(rows, []string{item.Name, item.Id, versions[0]})
		}
		HeaderStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
		EvenRowStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
		OddRowStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("245"))
		t := table.New().
			Border(lipgloss.NormalBorder()).
			BorderStyle(lipgloss.NewStyle().Foreground(lipgloss.Color("99"))).
			StyleFunc(func(row, col int) lipgloss.Style {
				switch {
				case row == 0:
					return HeaderStyle
				case row%2 == 0:
					return EvenRowStyle
				default:
					return OddRowStyle
				}
			}).
			Headers("Name", "Id", "Version").
			Width(s.Width).
			Rows(rows...)
		fmt.Println(t)
	}

}
