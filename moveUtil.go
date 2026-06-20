package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"
)

var (
	MUTEX sync.Mutex //MUTEX variable to intiate lock unlock mech for shared resources among Routines to prevent any inconsistency with the shared Resource
) 

func Worker(
	jobs <-chan CopyJob,
	fileSize *int,
	bufSize int,
	wg *sync.WaitGroup,
) {
	defer wg.Done()

	buffer := make([]byte, bufSize*1024) // declare a buffer of size bufSize

	for job := range jobs { // assinging jobs to the workers in this loop (this loop blocks the workers from working until there is a JOB in the job channel)

		jobStat, err := os.Stat(job.Source)
		if err != nil {
			fmt.Println(err)
			continue // this continue here is imp, as if I write Return here this will lead to the death of a worker (as returning means this a particular worker is now quit the worker func)
		}

		// open Sourcefile
		srcHandle, err2 := os.Open(job.Source)
		if err2 != nil {
			fmt.Println(err)
			continue
		}

		// create file at dest
		dstHandle, err3 := os.Create(job.Dest)
		if err3 != nil {
			srcHandle.Close() // if Dest is Faulty then close source File
			fmt.Println(err)
			continue
		}

		// copy data using buffer
		_, err4 := io.CopyBuffer(dstHandle, srcHandle, buffer)

		defer srcHandle.Close()
		defer dstHandle.Close()

		if err4 != nil {
			fmt.Println(err)
			continue
		}

		
		MUTEX.Lock()
		(*fileSize) += int(jobStat.Size()) // add fileSize
		MUTEX.Unlock()
		
		// uncomment the code below to delete the source after the movement of file is done
		
		// err5 = os.Remove(job.Source)
		// if err5 != nil {
		// 	fmt.Println("failed to delete:", job.Source)
		// }
	}
}

func MoveUtil(source, dest string, jobs chan<- CopyJob) {

	sourceStat, err := os.Stat(source) // retrive Info about Source
	if err != nil {
		return
	}

	basePath := filepath.Base(source) // retrive the basePath of source , so that we can create a named file at the dest

	if !sourceStat.IsDir() { // this is the base case of the recurrsive call of MoveUtil func

		// if the Source is a file then create a JOB and push it in the job channel
		jobs <- CopyJob{
			Source: source,
			Dest:   filepath.Join(dest, basePath),
		}
		return
	}

	// if source is a Dir then read its structure

	sourceStructure, _ := os.ReadDir(source)

	// make a named Dir at dest
	newDest := filepath.Join(dest, basePath)
	err = os.Mkdir(newDest, 0755)
	if err != nil {
		fmt.Println(err)
		return
	}

	// loop over all the entries in the Source and recursively call MoveUtil for each entry of source
	for _, entry := range sourceStructure {

		MoveUtil(
			filepath.Join(source, entry.Name()),
			newDest,
			jobs,
		)
	}
}
