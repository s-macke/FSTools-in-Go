package utils

import (
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"sync"
	"regexp"
)

func AddPath(path string, name string) string {
	var str strings.Builder
	str.WriteString(path)
	str.WriteString("/")
	str.WriteString(name)
	return str.String()
}

func FormatNumber(size int64, humanreadable bool) string {
	if !humanreadable {
		return strconv.FormatInt(size, 10)
	}
	if size < 4000 {
		return strconv.FormatInt(size, 10)
	}
	if size < 4000000 {
		return strconv.FormatInt(size>>10, 10)+"K"
	}
	if size < 4000000000 {
		return strconv.FormatInt(size>>20, 10)+"M"
	}
	return strconv.FormatInt(size>>30, 10)+"G"
}

// this is either a directory, file or symlink
type NodeInfo struct {
	Depth  uint // used by du
	Size int64 // used by du
}

func WildcardToRegex(pattern string) (*regexp.Regexp, error) {
	pattern = regexp.QuoteMeta(pattern)
	pattern = strings.Replace(pattern, "\\*", ".*", -1)
	pattern = strings.Replace(pattern, "\\?", ".", -1)
	pattern = pattern + "$"
	return regexp.Compile(pattern);
}


func ForEachDirEntry(
	path string,
	dir *NodeInfo,
	guard chan struct{},
	worker func(path string, node *NodeInfo, wg *sync.WaitGroup)) {

	guard <- struct{}{}
	f, err := os.Open(path)
	<-guard

	if err != nil {
		fmt.Println(err)
		return
	}

	var wg sync.WaitGroup

	for {
		names, err := f.Readdirnames(200)

		for _, name := range names {
			//fmt.Println(name)
			wg.Add(1)
			go worker(
				AddPath(path, name),
				dir,
				&wg)
		}

		if err == io.EOF {
			break
		}

		if err != nil {
			fmt.Println(err)
			break
		}

		if len(names) == 0 {
			break
		}
	}
	f.Close()
	wg.Wait()
}
