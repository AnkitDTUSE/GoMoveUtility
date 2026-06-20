package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

type CopyJob struct { // Each job will have these attributes
	Source string
	Dest   string
}

var (
	wg sync.WaitGroup // declaring a waitGroup for syncing the GO routines
)

func main() {
	// accepting arguements
	args := os.Args

	// retrive Source and Dest from args

	Source := args[1]
	Dest := args[2]

	fmt.Println("source is: ", Source)
	fmt.Println("Dest is: ", Dest)

	// checking the valdity of recieved paths
	_, err1 := os.Stat(Source)
	_, err2 := os.Stat(Dest)

	// checking for errors
	if err1 != nil {
		fmt.Println("error while reading Source")
		return
	}
	if err2 != nil {
		fmt.Println("error while reading Dest")
		return
	}

	// intialising filesize and bufSize variables

	fileSize := 0
	bufSize := 1024

	// read Input for bufSize
	fmt.Println("Enter Buffer size in KBs or hit enter for 1MB(1024KB) default buffsize")

	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)
	if input != "" {
		if val, err := strconv.Atoi(input); err == nil {
			bufSize = val
		}
	}

	// making a Buffered Channel (JOB) of 100 jobs at max
	jobs := make(chan CopyJob, 100)

	// initially 15 workers
	workers := 15

	for i := 0; i < workers; i++ { // Initializing each worker
		wg.Add(1)

		go Worker(
			jobs,
			&fileSize,
			bufSize,
			&wg,
		)
	}

	// record entry time
	t1 := time.Now()

	// call MoveUtility
	MoveUtil(Source, Dest, jobs)

	close(jobs) // unneccesary Step ,just to prevent the receivers from asking jobs from an empty job queue

	wg.Wait() // wait for each go routine's completion

	elapsed := time.Since(t1) // record exit time

	// Print Results
	fmt.Printf(
		"\nTotal time and total size in MBs => %v %vMB %vBytes\n",
		elapsed,
		float64(fileSize)/(1024*1024),
		fileSize,
	)
}

// BENCHMARKING -->

// bufSize	TransferTime	DataSize			 					TransferRate

// 5kb 		1m4.8684673s 	1332.5364799499512MB 1397265772Bytes -> 	21MB/s
// 10kb 	8.8855202s 		1406.5087223052979MB 1474831290Bytes -> 	156MB/s
// 5MB  	5.8288255s 		1410.5730619430542MB 1479093059Bytes -> 	235MB/s
// 10MB 	8.752762s 		1945.4457368850708MB 2039947709Bytes ->		233MB/s

// workers 		bufSize	TransferTime	DataSize			TransferRate

// 5     		5MB		61s				1332				22MBps
// 12			5MB		14s				1410				100MBps
// 20       	5MB		17s				1531				90MBps
// 15       	5MB		10s				1406				140MBps
// 10			5MB     13min			46GB				0.058GBps ~58MBps
// 15			5MB		12min			46GB				0.064GBps ~64MBps
// windows  	~NA~    11min			46GB				0.069GBps ~69MBps
// 	CopyPaste
// 17			4MB		13.4min			46GB				0.056GBps ~56MBps

// 15			5MB		8m30s			32.5GB				0.064GBps ~64MBps
// 17			5MB		9m				32.5GB				0.060GBps ~60MBps