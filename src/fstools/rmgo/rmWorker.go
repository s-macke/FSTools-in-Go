package main

import (
	"fstools/utils"
	"fmt"
	"os"
	"sync"
	"sync/atomic"
	"syscall"
	"time"
)

var stat struct {
	nremoveinqueue int32 // number of removes in queue
	nfiles         int32 // number of files/dir found so far
	nremoved       int32 // number of files removed
}

func HandleDir(path string, dir *utils.NodeInfo) {
	utils.ForEachDirEntry(path, dir, ForEachNodeWorker)

	// remove directory
	atomic.AddInt32(&stat.nremoveinqueue, 1)
	err := os.Remove(path)
	atomic.AddInt32(&stat.nremoveinqueue, -1)

	if err != nil {
		fmt.Println(err)
		return
	}
	atomic.AddInt32(&stat.nremoved, 1)
}

var guard = make(chan struct{}, 50)

func ForEachNodeWorker(path string, dir *utils.NodeInfo, wg *sync.WaitGroup) {
	guard <- struct{}{}

	// this is faster than to figure out if this is a file or directory by stat
	atomic.AddInt32(&stat.nremoveinqueue, 1)
	err := os.Remove(path)
	atomic.AddInt32(&stat.nremoveinqueue, -1)

	if err == nil {
		atomic.AddInt32(&stat.nremoved, 1)
		<- guard
		goto done;
	}

	if !Config.recursive {
		if err != nil {
			fmt.Println(err)
		}
		<- guard
		goto done;
	}

	if err.(*os.PathError).Err != syscall.ENOTEMPTY {
		fmt.Println(err)
		<- guard
		goto done;
	}
	<- guard
	HandleDir(path, dir)

done:
	wg.Done()
}

func PrintStat() {
	fmt.Printf("\rremoved / files / queue   %d / %d / %d", stat.nremoved, stat.nfiles, stat.nremoveinqueue)
}

func InitTicker() chan struct{} {
	ticker := time.NewTicker(500 * time.Millisecond)
    quit := make(chan struct{})
    if !Config.verbose {
		return quit
	}

    go func() {
        for {
            select {
            case <- ticker.C:
				PrintStat()
            case <- quit:
                ticker.Stop()
                return
            }
        }
    }()

    return quit
}

func RemoveDir() {
	guard = make(chan struct{}, Config.workers)
	quit := InitTicker()

	var wg sync.WaitGroup

	for _, rootpath := range Config.rootpaths {
		wg.Add(1)
		atomic.AddInt32(&stat.nfiles, 1)
		node := utils.NodeInfo{}
		go ForEachNodeWorker(rootpath, &node, &wg)
	}

	wg.Wait()

	if Config.verbose {
		quit <- struct{}{}
		PrintStat()
	}
}