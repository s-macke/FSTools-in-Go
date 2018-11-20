package main

import (
    "flag"
    "fmt"
    "math"
    "os"
    "time"
)

var Config struct {
    workers int
    rootpaths []string
    summary bool
    humanreadable bool
    maxdepth int
}

func parse() {
    flag.IntVar(&Config.workers,"w", 10, "Number of go workers.")
    flag.BoolVar(&Config.summary,"s", false, "display only a total")
    flag.BoolVar(&Config.humanreadable, "h", false, "human readable")
    flag.IntVar(&Config.maxdepth, "d", math.MaxInt32, "print the total for a directory\nonly if it is N or fewer levels\nbelow the command line argument")

    flag.Usage = func() {
        fmt.Fprintf(os.Stderr, "Usage of %s [OPTION]... [FILE]...\n", os.Args[0])
        fmt.Fprintf(os.Stderr, "Summarize disk usage of the set of FILEs, recursively for directories.\n\n")
        flag.PrintDefaults()
    }

    flag.Parse()

    Config.rootpaths = []string{"."}
    if len(flag.Args()) >= 1 {
        Config.rootpaths = flag.Args()
    }
}

func main() {
    parse()

    start := time.Now()

    ReadDir()

    elapsed := time.Since(start)
    fmt.Println("dugo took", elapsed)
}

