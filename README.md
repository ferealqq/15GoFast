# README

## About

15Puzzle game solver

## Reports

[week 1](https://github.com/ferealqq/15GoFast/blob/main/documentation/week1.md)
[week 2](https://github.com/ferealqq/15GoFast/blob/main/documentation/week2.md)
[week 3](https://github.com/ferealqq/15GoFast/blob/main/documentation/week4.md)
[week 4](https://github.com/ferealqq/15GoFast/blob/main/documentation/week4.md)
[week 5](https://github.com/ferealqq/15GoFast/blob/main/documentation/week5.md)
[week 6](https://github.com/ferealqq/15GoFast/blob/main/documentation/week6.md)

After you have successfully started the documentation server you can find the documentation via this [path](http://localhost:6060/pkg/github.com/ferealqq/15GoFast/)


## Live Development

To run in live development mode, run `wails dev` in the project directory. This will run a Vite development
server that will provide very fast hot reload of your frontend changes. If you want to develop in a browser
and have access to your Go methods, there is also a dev server that runs on http://localhost:34115. Connect
to this in your browser, and you can call your Go code from devtools.

## Building

To build a redistributable, production mode package, use `wails build`.


## Testing

To test the core logic of the application run: 
```terminal
# takes 5s to 25s depending on how hard the random generated puzzles are.
go test . 
```

Some times the tests may exceed the time limit of solving the puzzle board. But exceeding the time limit should be very rare.

## Profiling 


Profiling is done with [pprof](https://pkg.go.dev/runtime/pprof).
Usage:
```terminal
go tool pprof "profile_file" 
```

To generate profile's use the following [guide](https://pkg.go.dev/runtime/pprof).

Generate IDAStar's cpu profile run
```terminal
go clean --cache ; go clean -testcache ; go test -run "TestPerformance" -cpuprofile cpu.prof
go tool pprof search.prof
(pprof) web
```


## Documentation

## GO pkg documentation
If you want to see the documentation of the application you'll have to have docker installed on your computer. 

To start the documentation server use script `run_documentation.sh`
Example: 
```terminal
chmod +x run_documentation.sh
./run_documentation.sh
```