# FSTools-in-Go
Parallelized versions of du and rm in Go

This repository contains simplified versions of the filesystem tools du and rm

* du - estimate file space usage
* rm - remove files or directories

They have been specially developed for network file systems such as NFS or CIFS, which contain an unusually large number of files. The usual POSIX utilities are not parallelized. Therefore every file system command needs a roundtrip on the network, which makes these tools very inefficient.

Both tools provided here can mitigate these problems by parallel access to the file system. 
