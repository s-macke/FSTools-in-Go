# FSTools-in-Go
Parallelized versions of du and rm in Go

This repository contains simplified versions of the filesystem tools du and rm

* du - estimate file space usage
* rm - remove files or directories

They have been specially developed for network file systems such as NFS or CIFS, which contain an unusually large number of files. The usual POSIX utilities are not parallelized. Therefore every file system command needs a roundtrip on the network, which makes these tools very inefficient.

Both tools provided here can mitigate these problems by parallel access to the file system. 

## Usage of dugo
```
  -d int
    	print the total for a directory
    	only if it is N or fewer levels
    	below the command line argument (default 2147483647)
  -h	human readable
  -s	display only a total
  -w int
    	Number of go workers. Default is 10 (default 10)
      A high number of workers floods your kernel with commands and might block other network tasks.
```

## Usage of rmgo
```  
  -r	Remove directories and their contents recursively
  -v	Verbose
  -w int
    	Number of go workers. Default is 10 (default 10)
      A high number of workers floods your kernel with commands and might block other network tasks.
```
