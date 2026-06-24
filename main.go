package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/maxdikun/raii/internal"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "Usage: raii <start|stop|check|watch> [flags]")
		os.Exit(1)
	}

	cmd := os.Args[1]
	if cmd == "-h" || cmd == "--help" || cmd == "help" {
		fmt.Println("Usage: raii <start|stop|check|watch> [flags]")
		fmt.Println("")
		fmt.Println("Commands:")
		fmt.Println("  start   Start resources and register owner")
		fmt.Println("  stop    Unregister owner and stop resources if no owners remain")
		fmt.Println("  check   Run the check command")
		fmt.Println("  watch   Watch an owner PID and stop when it dies")
		fmt.Println("")
		fmt.Println("Flags:")
		fmt.Println("  --config string   Path to config file (default \"raii.toml\")")
		fmt.Println("  --owner string    Owner identifier (default: parent PID)")
		os.Exit(0)
	}

	var configPath, owner string
	flags := flag.NewFlagSet(cmd, flag.ExitOnError)
	flags.StringVar(&configPath, "config", "", "path to config file")
	flags.StringVar(&owner, "owner", "", "owner identifier (default: parent pid)")
	flags.Parse(os.Args[2:])

	args := flags.Args()
	if configPath == "" && len(args) > 0 {
		configPath = args[0]
	}
	if configPath == "" {
		configPath = "raii.toml"
	}

	var err error
	switch cmd {
	case "start":
		err = internal.Start(configPath, owner)
	case "stop":
		err = internal.Stop(configPath, owner)
	case "check":
		err = internal.Check(configPath)
	case "watch":
		err = internal.Watch(configPath, owner)
	default:
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n", cmd)
		os.Exit(1)
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
