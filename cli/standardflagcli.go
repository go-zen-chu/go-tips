package main

import (
	"flag"
	"fmt"
	"log"
	"os"
)

var (
	goFlagSet        = flag.NewFlagSet("go", flag.ExitOnError)
	debugFlag        = goFlagSet.Bool("debug", false, "enable debug")
	helpFlag         = goFlagSet.Bool("help", false, "show help")
	buildFlagSet     = flag.NewFlagSet("build", flag.ExitOnError)
	buildVerboseFlag = buildFlagSet.Bool("n", false, "print the commands but do not run them.")
	output           = "./"
)

func init() {
	buildFlagSet.Var(&PathValue{Path: &output}, "o", "named `output` file")
}

func help() {
	fmt.Println("usage: go [global flags] <command>")
	goFlagSet.PrintDefaults()
	fmt.Println("\ngo build [-o output] [-n]")
	buildFlagSet.PrintDefaults()
}

func main() {
	if len(os.Args) == 1 {
		help()
		os.Exit(2)
	}
	// handling root command
	if err := goFlagSet.Parse(os.Args[1:]); err != nil {
		log.Fatalf("parse root command: %s", err)
	}
	if *helpFlag {
		help()
		os.Exit(0)
	}
	if *debugFlag {
		log.SetPrefix("DEBUG: ")
		log.Println("debug mode is on")
	}
	// handling subcommands
	subCommandArgs := os.Args[1+goFlagSet.NFlag():] // NFlag is number of flags for root command
	if len(subCommandArgs) == 0 {
		help()
		log.Fatalln("specify subcommand")
	}
	switch subCommand := subCommandArgs[0]; subCommand {
	case "build":
		if err := buildFlagSet.Parse(subCommandArgs[1:]); err != nil {
			log.Fatalf("parse subcommand build: %s", err)
		}
		if *buildVerboseFlag {
			log.Println("verbose is on")
		}
		if output != "" {
			log.Printf("output path: %s\n", output)
		} else {
			log.Printf("%v\n", output)
			log.Println("no output path")
		}
	default:
		help()
		os.Exit(2)
	}

	log.Printf("--- Successful ---")
	// DO SOMETHING
	log.Printf("%v\n", os.Args)
	log.Printf("%v\n", subCommandArgs)
}

type PathValue struct {
	Path *string
}

func (v PathValue) String() string {
	if v.Path == nil {
		return ""
	}
	return *v.Path
}

func (v PathValue) Set(s string) error {
	if _, err := os.Stat(s); err != nil {
		return err
	}
	*v.Path = s
	return nil
}
