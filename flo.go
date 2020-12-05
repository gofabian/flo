package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/gofabian/flo/cli"
)

func main() {
	commands := make(map[string]func(*flag.FlagSet, []string) error)
	commands["convert-pipeline"] = cli.ConvertPipeline

	help := flag.Bool("h", false, "Print help text")
	flag.Usage = func() {
		fmt.Printf("Usage: %s [COMMAND] [OPTIONS]\n", os.Args[0])
		flag.PrintDefaults()
		fmt.Printf("\nCommands:\n")
		for name := range commands {
			fmt.Printf("  %s\n", name)
		}
	}

	if len(os.Args) < 2 {
		fmt.Printf("Missing arguments!\n\n")
		flag.Usage()
		os.Exit(1)
	}

	for name, execute := range commands {
		if os.Args[1] == name {
			flagSet := flag.NewFlagSet(name, flag.ExitOnError)
			flagSet.Usage = func() {
				fmt.Printf("Usage: %s %s [OPTIONS]\n", os.Args[0], os.Args[1])
				flagSet.PrintDefaults()
			}
			flagSet.Bool("h", false, "Print help text")
			err := execute(flagSet, os.Args[2:])
			if err != nil {
				fmt.Println(err.Error())
				os.Exit(1)
			}
			return
		}
	}

	flag.Parse()
	if *help {
		flag.Usage()
		return
	}
	fmt.Printf("Invalid arguments!\n\n")
	flag.Usage()
	os.Exit(1)
}
