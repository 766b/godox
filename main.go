package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {
	flagPrintHelp := flag.Bool("h", false, "Print help")
	flag.Usage = func() {
		_, _ = fmt.Fprintf(os.Stderr, "Usage: godox [path] [path] ...\n")
		flag.PrintDefaults()
	}
	flag.Parse()

	if *flagPrintHelp {
		flag.Usage()
		return
	}

	if len(os.Args) == 1 {
		newGodox(".", os.Stdout).parse()
		return
	}

	for _, arg := range os.Args {
		fs, err := os.Stat(arg)
		if err != nil {
			panic(err)
		}
		if fs.IsDir() {
			newGodox(arg, os.Stdout).parse()
		}
	}
}
