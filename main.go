package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/Yves848/wingettui/winget"
)

func main() {
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

		params[0] = "-command Get-WinGetPackage | where-object { $_.Source -eq \"winget\"} | ConvertTo-Json -Depth 1 -WarningAction SilentlyContinue"
		if *listSource != "" {

		}

	case "search":
		searchCmd.Parse(os.Args[2:])
		fmt.Println("searchQuery :  " + *searchQuery)
		params = make([]string, 1)
		params[0] = "Search-WGPackage \"##QUERY##\" -quiet $true | out-json"
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
		items, err := winget.Winget2Json(string(out))
		if err != nil {
			fmt.Println(errors.New("error parsing json"))
			os.Exit(2)
		}
		for _, item := range items.Packages {
			fmt.Printf("%s %s\n", item.Id, item.InstalledVersion)
		}
	}
}
