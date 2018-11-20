package main

import (
	"fstools/utils"
	"fmt"
	"sync"
	"os"
	"sync/atomic"
)

var guard = make(chan struct{}, 10)

func HandleDir(path string, dir *utils.NodeInfo) {
	utils.ForEachDirEntry(path, dir, guard, ForEachNodeWorker)

	if (Config.summary && dir.Depth<=1) || (!Config.summary && dir.Depth<=Config.maxdepth+1) {
		fmt.Println(
			utils.FormatNumber(
				atomic.LoadInt64(&dir.Size),
				Config.humanreadable) +
				"\t" + path)
	}
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

	if fip.IsDir() {
    	nodenew := utils.NodeInfo{Depth: node.Depth + 1, Size: fip.Size()}
		HandleDir(path, &nodenew)
		atomic.AddInt64(&node.Size, nodenew.Size)
		return
	}

	atomic.AddInt64(&node.Size, fip.Size())
}

func ReadDir() {
	guard = make(chan struct{}, Config.workers)
	//guarddir = make(chan struct{}, Config.workers)

	var wg sync.WaitGroup
	for _, rootpath := range Config.rootpaths {
		wg.Add(1)
		node := utils.NodeInfo{Depth: 0, Size: 0}
		go ForEachNodeWorker(rootpath, &node, &wg)
	}
	wg.Wait()
}