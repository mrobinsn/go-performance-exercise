# go-performance-exercise
> Exercise to get familiar with profiling and performance tuning in Go

# Goal
In the `data` directory, there are input files that contain a JSON record on each line:
```
{"value":"... some string ..."}
{"value":"... some string ..."}
{"value":"... some string ..."}
{"value":"... some string ..."}
...
```
The objective is for your code to return the number of buffalos found within the record set,
meaning, the number of records that contain the word "buffalo". You can assume that a record
can only ever have a maximum of one buffalo.

You will be given an existing naive implementation. From there you should use performance profiling tools and your intuition to make improvements. You can make tweaks and use the benchmark tool to compare it to a given set of 'golden' implementations (source to be uploaded after completion of this exercise). 

At the end of this exercise, I will go over the 'golden' implementations to discuss the performance characteristics as well how it ties into 'clear vs. clever', premature optimization pitfalls, and the importance of benchmarking.

# Instructions
Before you start, be sure to decompress the test files:
```
gunzip data/*.gz
```

`main.go` contains an implementation as a starting point.

Edit this code to make your performance improvements. 

## Test & Benchmark
> Binaries in the repo are pre-compiled for OSX

The `run` command will automatically build your attempt and run it against the 'golden' implementation,
reporting the run times or failure if your code is incorrect.
```
./run
```

You can also manually test your implementation
```
cat data/raw1000 | go run main.go
```

## Profiling
> So how do I make the naive implementation more performant? 

### `pprof`
> https://golang.org/pkg/runtime/pprof/ & https://godoc.org/net/http/pprof
```
go tool pprof --help
```

`pprof` is a tool that helps with analyzing cpu and memory profiles for Go applications. 

It is comprised of two components: a CLI tool, and a package to generate the profiles.

In a web application you can use the `net/http/pprof` package to add endpoints that return cpu/memory profiles.

In an application without an HTTP interface, you can use the `runtime/pprof` to directly generate and store the profiles. 

[Good blog post on the subject](https://jvns.ca/blog/2017/09/24/profiling-go-with-pprof/)

The prepared implementation in `main.go` already contains the necessary code to create cpu and memory profiles. Which can be generated with the following flag:

```
cat data/raw10000 | go run main.go -profile
```
>You can use whichever data file you want, depending on the performance characteristics of your code, some trends may only appear on larger datasets.

This will generate files `cpu.profile` and `mem.profile` in your local directory.

Caveat -- on OSX, CPU profiling may not function correctly. I've had limited success but have found it easier to simply run your code within a Linux environment and pull profiles from there. Docker makes this super easy for us:
```
GOOS=linux go build -o attempt-linux . && \
docker run -it -v $(pwd):/code -w /code golang \
sh -c "cat data/raw1000000 | ./attempt-linux -profile"
```

Now you can feed those profiles into the `go tool pprof` command to get valuable insights.

Create a SVG displaying memory allocations:
```
go tool pprof -svg mem.profile
```

Create a SVG displaying CPU usage
```
go tool pprof -svg cpu.profile
```

We'll discuss these examples and what insights can be gleamed from the reports.

***

### Appendix
Commands to build the `golden-*` and `run` binaries (uses `upx` for reduced binary size)

```
go build -ldflags="-s -w" -o golden-simple ./cmd/golden-simple && \
go build -ldflags="-s -w" -o golden-complicated ./cmd/golden-complicated && \
go build -ldflags="-s -w" -o run ./cmd/run && \
upx --brute run
```
