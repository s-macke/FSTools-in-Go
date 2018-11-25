package main

import (
    "fstools/utils"
    "flag"
    "fmt"
    "math"
    "os"
    "regexp"
    "strings"
    "time"
)

var Config struct {
    workers uint
    rootpaths []string
    maxdepth uint
    typ string

    name string
    nameregex *regexp.Regexp
    iname string
    inameregex *regexp.Regexp
}

func parse() {
    flag.UintVar(&Config.workers,"w", 10, "Number of go workers.")
    flag.UintVar(&Config.maxdepth, "maxdepth", math.MaxInt32, "Descend at most levels (a non-negative integer) levels of directories below the starting-points")
    flag.StringVar(&Config.typ, "type", "", "File is of type file (f) or directory (d)")
    flag.StringVar(&Config.name, "name", "", "Base of file name (the path with the leading directories removed) matches the pattern string.")
    flag.StringVar(&Config.iname, "iname", "", "Like -name, but the match is case insensitive")

    flag.Usage = func() {
        fmt.Fprintf(os.Stderr, "Usage of %s [OPTION]... [FILE]...\n", os.Args[0])
        fmt.Fprintf(os.Stderr, "Search for files in a directory hierarchy.\n\n")
        flag.PrintDefaults()
    }

    flag.Parse()

    if Config.typ != "" && Config.typ != "f" && Config.typ != "d" {
        flag.Usage()
        os.Exit(1)
    }

    if Config.name != "" {
        r, err := utils.WildcardToRegex(Config.name);
        if err != nil {
            fmt.Println(err)
            os.Exit(1)
        }
        Config.nameregex = r
        //fmt.Println(Config.nameregex);
    }

    if Config.iname != "" {
        r, err := utils.WildcardToRegex(strings.ToLower(Config.iname));
        if err != nil {
            fmt.Println(err)
            os.Exit(1)
        }
        Config.inameregex = r
        //fmt.Println(Config.inameregex);
    }

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
    fmt.Println("findgo took", elapsed)
}

