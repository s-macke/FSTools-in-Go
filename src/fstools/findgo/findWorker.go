package main

import (
	"fstools/utils"
	"fmt"
	"strings"
	"sync"
	"os"
	"sync/atomic"
)

var guard = make(chan struct{}, 10)

func MatchFileStep2(name string, path string) {
	if (Config.nameregex != nil && Config.nameregex.MatchString(name)) {
		fmt.Println(path)
	} else
	if (Config.inameregex != nil && Config.inameregex.MatchString(strings.ToLower(name))) {
		fmt.Println(path)
	} else
	if (Config.nameregex == nil && Config.inameregex == nil) {
		fmt.Println(path)
	}
}

func MatchFileStep1(fip os.FileInfo, path string) {
	if (Config.typ == "d" && fip.IsDir()) {
		MatchFileStep2(fip.Name(), path)
	} else
	if (Config.typ == "f" && !fip.IsDir()) {
		MatchFileStep2(fip.Name(), path)
	} else
	if (Config.typ == "") {
		MatchFileStep2(fip.Name(), path)
	}
}

func HandleDir(path string, dir *utils.NodeInfo) {
	if dir.Depth > Config.maxdepth {
		return
	}
	utils.ForEachDirEntry(path, dir, guard, ForEachNodeWorker)
}

func ForEachNodeWorker(path string, node *utils.NodeInfo, wg *sync.WaitGroup) {
	defer wg.Done()
	guard <- struct{}{}

	fip, lerr := os.Stat(path)
	if lerr != nil {
		fmt.Println(lerr)
		<-guard
		return
	}
	<-guard

	MatchFileStep1(fip, path)

	if fip.IsDir() {
    	nodenew := utils.NodeInfo{Depth: node.Depth + 1}
		HandleDir(path, &nodenew)
		return
	}

	atomic.AddInt64(&node.Size, fip.Size())
}

func ReadDir() {
	guard = make(chan struct{}, Config.workers)

	var wg sync.WaitGroup
	for _, rootpath := range Config.rootpaths {
		wg.Add(1)
		node := utils.NodeInfo{Depth: 0}
		go ForEachNodeWorker(rootpath, &node, &wg)
	}
	wg.Wait()
}