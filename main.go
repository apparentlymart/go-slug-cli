package main

import (
	"flag"
	"fmt"
	"io"
	"os"

	slug "github.com/hashicorp/go-slug"
)

func main() {
	err := realMain(os.Args[1:])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		os.Exit(1)
	}
}

func realMain(args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("Usage: go-slug-cli <subcommand> [args...]")
	}

	switch args[0] {
	case "pack":
		return cmdPack(args[1:])
	default:
		return fmt.Errorf("unsupported subcommand %q", args[0])
	}
}

func cmdPack(args []string) error {
	opts := flag.NewFlagSet("pack", flag.ExitOnError)
	outFileP := opts.String("out", "", "tar archive output filename")
	err := opts.Parse(args)
	if err != nil {
		return fmt.Errorf("invalid options: %s", err)
	}
	args = opts.Args()

	dir := "."
	if len(args) > 0 {
		dir = args[0]
	}

	var outW io.Writer
	if *outFileP != "" {
		outFile, err := os.Create(*outFileP)
		if err != nil {
			return fmt.Errorf("creating output file: %s", err)
		}
		outW = outFile
	} else {
		outW = io.Discard
	}

	meta, err := slug.Pack(dir, outW, true)
	if err != nil {
		return fmt.Errorf("creating slug archive: %s", err)
	}

	for _, fn := range meta.Files {
		fmt.Println(fn)
	}

	return nil
}
