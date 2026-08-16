package main

import (
	"fmt"
	"io"
	"os"

	"log-shipper/internal/bundle"
	"log-shipper/internal/collect"
	"log-shipper/internal/index"
	"log-shipper/internal/report"
)

func usage(w io.Writer) {
	fmt.Fprintln(w, "usage:")
	fmt.Fprintln(w, "  log-shipper collect --src DIR --out BUNDLE")
	fmt.Fprintln(w, "  log-shipper search --bundle BUNDLE --keyword K")
}

func main() {
	if err := run(os.Args[1:], os.Stdout, os.Stderr); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}

func run(args []string, stdout, stderr io.Writer) error {
	if len(args) < 1 {
		usage(stderr)
		return fmt.Errorf("missing command")
	}
	switch args[0] {
	case "collect":
		return cmdCollect(args[1:], stdout, stderr)
	case "search":
		return cmdSearch(args[1:], stdout, stderr)
	default:
		usage(stderr)
		return fmt.Errorf("unknown command %q", args[0])
	}
}

func cmdCollect(args []string, stdout, stderr io.Writer) error {
	var src, out string
	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "--src":
			if i+1 < len(args) {
				src = args[i+1]
				i++
			}
		case "--out":
			if i+1 < len(args) {
				out = args[i+1]
				i++
			}
		}
	}
	if src == "" || out == "" {
		usage(stderr)
		return fmt.Errorf("collect requires --src and --out")
	}
	recs, err := collect.Collect(src)
	if err != nil {
		return err
	}
	if err := bundle.Write(out, recs); err != nil {
		return err
	}
	fmt.Fprintf(stdout, "collected %d log files into %s\n", len(recs), out)
	return nil
}

func cmdSearch(args []string, stdout, stderr io.Writer) error {
	var bundlePath, keyword string
	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "--bundle":
			if i+1 < len(args) {
				bundlePath = args[i+1]
				i++
			}
		case "--keyword":
			if i+1 < len(args) {
				keyword = args[i+1]
				i++
			}
		}
	}
	if bundlePath == "" || keyword == "" {
		usage(stderr)
		return fmt.Errorf("search requires --bundle and --keyword")
	}
	recs, err := bundle.Read(bundlePath)
	if err != nil {
		return err
	}
	idx := index.Build(recs)
	hits := idx.Search(keyword)
	fmt.Fprintf(stdout, "matches for %q: %d\n", keyword, len(hits))
	report.Render(hits, stdout)
	return nil
}
