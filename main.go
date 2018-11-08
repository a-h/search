package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/a-h/search/contains"

	"github.com/a-h/search/searcher"
)

var flagIncludeNames = flag.String("names", "", "path names to search for")
var flagExcludeNames = flag.String("exclude-names", ".git", "path names to exclude")
var flagIncludeDirectories = flag.Bool("directories", true, "set to true to include directories")
var flagIncludeFiles = flag.Bool("files", true, "set to true to include files")
var flagIncludeText = flag.String("text", "", "text to search for")
var flagExcludeText = flag.String("exclude-text", "", "text to search for")
var flagPrintSettings = flag.Bool("print-settings", true, "prints out the settings before searching")
var flagPrintSummary = flag.Bool("print-summary", true, "prints out a summary after searching")

func main() {
	flag.Parse()

	tail := flag.Args()
	if len(os.Args) == 1 || len(tail) == 0 || len(tail) > 1 {
		fmt.Fprintln(os.Stdout, "usage: search [<args>] directory")
		flag.PrintDefaults()
		return
	}

	var wd string
	wd, err := os.Getwd()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to get working directory: %v\n", err)
		os.Exit(1)
	}
	wd, err = filepath.Abs(tail[0])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to get directory '%v': %v\n", tail[0], err)
		os.Exit(1)
	}

	ss := searcher.Settings{
		IncludeNames:       spaceSeparated(*flagIncludeNames),
		ExcludeNames:       spaceSeparated(*flagExcludeNames),
		IncludeDirectories: *flagIncludeDirectories,
		IncludeFiles:       *flagIncludeFiles,
		IncludeText:        *flagIncludeText,
		ExcludeText:        *flagExcludeText,
		TextSearch:         contains.TextInFile,
	}

	if *flagPrintSettings {
		fmt.Fprintf(os.Stdout, "%s: %v\n", "Directory", wd)
		fmt.Fprintf(os.Stdout, "%s: %v\n", "Include Names", ss.IncludeNames)
		fmt.Fprintf(os.Stdout, "%s: %v\n", "Exclude Names", ss.ExcludeNames)
		fmt.Fprintf(os.Stdout, "%s: %v\n", "Include Directories", ss.IncludeDirectories)
		fmt.Fprintf(os.Stdout, "%s: %v\n", "Include Files", ss.IncludeFiles)
		fmt.Fprintf(os.Stdout, "%s: %v\n", "Include Text", ss.IncludeText)
		fmt.Fprintf(os.Stdout, "%s: %v\n", "Exclude Text", ss.ExcludeText)
	}

	// Cancel if a signal if received.
	ctx, cancel := context.WithCancel(context.Background())
	sigs := make(chan os.Signal)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigs
		fmt.Fprintf(os.Stderr, "\nShutdown received.\n")
		cancel()
	}()
	defer cancel()

	s := searcher.New(ss)

	var summary searcher.Summary

	paths := make(chan string)
	errors := make(chan error)
	done := make(chan bool)
	go func() {
		var err error
		summary, err = s.Walk(ctx, wd, paths, errors)
		if err != nil && err != searcher.ErrCancelled {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		}
		close(paths)
		close(errors)
	}()

	go func() {
		for path := range paths {
			fmt.Fprintln(os.Stdout, path)
		}
		done <- true
	}()

	go func() {
		for err := range errors {
			fmt.Fprintf(os.Stderr, "%v\n", err)
		}
		done <- true
	}()
	for i := 0; i < 2; i++ {
		<-done
	}
	if *flagPrintSummary {
		fmt.Fprintln(os.Stdout, summary.String())
	}
}

func spaceSeparated(s string) []string {
	if s == "" {
		return []string{}
	}
	return strings.Split(strings.TrimSpace(s), " ")
}
