## Go solution (Part 2)

Use a priority queue for buy and sell queues, to avoid sorting.

## Build tools

Install [go](https://golang.org)

## Building and running the Go solution

```
% go build && gzip -dc ../../orders-100K.txt.gz | time ./Exchange > trades.txt
```

## Profiling the Go solution

Import [`github.com/pkg/profile`](https://github.com/pkg/profile) and add this line to the start of `func main()`:

```
defer profile.Start().Stop()
```

Running the Go solution like above will generate a `.pprof` file, which can be viewed with

```
% go tool pprof -http localhost:7999 /path/to/cpu.pprof
```

## Running benchmarks

```
% go test -bench . -benchmem
```
