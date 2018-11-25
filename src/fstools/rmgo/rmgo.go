package main

import (
    "flag"
    "fmt"
    "os"
    "time"
)

var Config struct {
    workers uint
    rootpaths []string
    recursive bool
    verbose bool
}

func parse() {
    flag.UintVar(&Config.workers,"w", 10, "Number of go workers.")
    flag.BoolVar(&Config.recursive, "r", false, "Remove directories and their contents recursively")
    flag.BoolVar(&Config.verbose, "v", false, "Verbose")

    flag.Usage = func() {
        fmt.Fprintf(os.Stderr, "Usage of %s [OPTION]... [FILE]...\n", os.Args[0])
        fmt.Fprintf(os.Stderr, "Remove (unlink) the FILE(s).\n\n")
        flag.PrintDefaults()
    }

    flag.Parse()

    Config.rootpaths = []string{"."}
    if len(flag.Args()) == 0 {

        flag.Usage()
        os.Exit(1)
    }
    Config.rootpaths = flag.Args()
}

func main() {
    parse()

    start := time.Now()

    RemoveDir()

    elapsed := time.Since(start)
    fmt.Printf("\nrmgo took %v\n", elapsed)
}

