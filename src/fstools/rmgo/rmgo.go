package main

import (
    "flag"
    "fmt"
    "os"
    "time"
)

var Config struct {
    workers int
    rootpaths []string
    dryrun bool
    recursive bool
    verbose bool
}

func parse() {
    flag.IntVar(&Config.workers,"w", 10, "Number of go workers. Default is 10")
    flag.BoolVar(&Config.recursive, "r", false, "Remove directories and their contents recursively")
    flag.BoolVar(&Config.verbose, "v", false, "Verbose")

    flag.Parse()

    Config.rootpaths = []string{"."}
    if len(flag.Args()) == 0 {
        flag.PrintDefaults()
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

