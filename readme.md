# GoMove ⚡

A high-performance concurrent file transfer utility built with **Go**, designed to recursively copy/move files and directories using a configurable worker pool and buffered I/O.

The utility leverages **Go's goroutines, channels, mutexes, and wait groups** to achieve transfer speeds comparable to native operating system file managers such as Windows Explorer's Copy/Paste and Move operations.

---

## Features

* 🚀 Concurrent file transfers using worker goroutines
* 📂 Recursive directory traversal
* 🔄 Preserves source directory structure at destination
* ⚙️ Configurable buffer size for performance tuning
* 🧵 Worker-pool architecture using channels
* 📊 Transfer statistics and benchmarking
* 🔒 Thread-safe shared state management using mutexes
* 💾 Memory-efficient buffered copying with `io.CopyBuffer`
* 🏗 Cross-platform (Windows, Linux, macOS)

---

## Architecture

The utility consists of three major components:

### 1. Directory Scanner (`MoveUtil`)

Recursively traverses the source directory structure.

* Detects files and directories
* Creates destination directories
* Generates copy jobs
* Pushes jobs into a buffered channel

### 2. Job Queue

A buffered channel stores transfer jobs:

```go
type CopyJob struct {
    Source string
    Dest   string
}
```

This decouples file discovery from file transfer execution.

### 3. Worker Pool

Multiple worker goroutines consume jobs concurrently.

Each worker:

1. Opens source file
2. Creates destination file
3. Copies data using a configurable buffer
4. Updates transferred size statistics
5. (Optional) Deletes source file after successful transfer

---

## How It Works

```text
Source Directory
       │
       ▼
  MoveUtil()
       │
       ▼
   Job Channel
       │
 ┌─────┼─────┐
 ▼     ▼     ▼
Worker Worker Worker
  1      2      N
       │
       ▼
 Destination
```

---

## Installation

Clone the repository:

```bash
git clone https://github.com/yourusername/gomove.git

cd gomove
```

Build the executable:

```bash
go build
```

Or run directly:

```bash
go run .
```

---

## Usage

```bash
go run . <source_path> <destination_path>
```

Example:

```bash
go run . "D:\Movies" "E:\Backup"
```

The utility will then prompt for a buffer size:

```text
Enter Buffer size in KBs or hit enter for 1MB(1024KB) default buffsize
```

Examples:

```text
1024
```

for 1 MB buffer

```text
5120
```

for 5 MB buffer

or simply press Enter to use the default.

---

## Configuration

### Worker Count

Current default:

```go
workers := 15
```

You can tune this value based on:

* CPU cores
* Storage type (SSD/HDD/NVMe)
* File sizes
* System memory

---

### Buffer Size

Default:

```go
bufSize := 1024 // KB
```

Larger buffers generally improve throughput for large files while consuming more memory.

---

## Optional Move Mode

Currently, the utility performs a copy operation.

To enable true move behavior, uncomment:

```go
err5 = os.Remove(job.Source)
if err5 != nil {
    fmt.Println("failed to delete:", job.Source)
}
```

inside the worker function.

This deletes the source file after a successful copy.

---

## Performance Benchmarks

### Buffer Size Tests

| Buffer Size | Transfer Time | Data Size | Transfer Rate |
| ----------- | ------------- | --------- | ------------- |
| 5 KB        | 1m 4s         | 1332 MB   | ~21 MB/s      |
| 10 KB       | 8.8s          | 1406 MB   | ~156 MB/s     |
| 5 MB        | 5.8s          | 1410 MB   | ~235 MB/s     |
| 10 MB       | 8.7s          | 1945 MB   | ~233 MB/s     |

---

### Worker Count Tests

| Workers | Buffer | Data Size | Transfer Speed |
| ------- | ------ | --------- | -------------- |
| 5       | 5 MB   | 1332 MB   | ~22 MB/s       |
| 12      | 5 MB   | 1410 MB   | ~100 MB/s      |
| 15      | 5 MB   | 1406 MB   | ~140 MB/s      |
| 20      | 5 MB   | 1531 MB   | ~90 MB/s       |

---

### Large Dataset Benchmark

| Method                      | Dataset | Time    |
| --------------------------- | ------- | ------- |
| GoMove (15 Workers)         | 46 GB   | ~12 min |
| Windows Explorer Copy/Paste | 46 GB   | ~11 min |

Results demonstrate that GoMove achieves performance highly comparable to native Windows file transfer operations while providing full control over concurrency and buffering.

---

## Concurrency Concepts Used

This project demonstrates practical use of:

* Goroutines
* Channels
* Buffered Channels
* WaitGroups
* Mutexes
* Worker Pools
* Recursive File Traversal
* Buffered I/O
* Synchronization Primitives

---

## Future Improvements

* Progress bar
* Transfer rate monitoring
* Resume interrupted transfers
* Dynamic worker scaling
* File integrity verification (checksums)
* Symbolic link support
* Configuration file support
* CLI flags for workers and buffer size
* Logging and error reports

---

## Example Output

```text
source is: D:\Movies
Dest is: E:\Backup

Enter Buffer size in KBs or hit enter for 1MB(1024KB) default buffsize

5120

Total time and total size in MBs =>
10.2s 1406MB 1474831290Bytes
```

---

## Why This Project?

This project was built as a practical exploration of Go's concurrency model and demonstrates how goroutines and channels can be used to build a real-world, high-performance file transfer utility capable of handling large directory trees efficiently.

---

## License

This project is currently not licensed.
