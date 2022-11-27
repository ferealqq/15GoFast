# README

## About

15Puzzle game solver

## Documentation
If you want to see the documentation of the application you'll have to have docker installed on your computer. 


To start the documentation server use script `run_documentation.sh`
Example: 
```terminal
chmod +x run_documentation.sh
./run_documentation.sh
```

After you have successfully started the documentation server you can find the documentation via this [path](http://localhost:6060/pkg/github.com/ferealqq/15GoFast/)


## Live Development

To run in live development mode, run `wails dev` in the project directory. This will run a Vite development
server that will provide very fast hot reload of your frontend changes. If you want to develop in a browser
and have access to your Go methods, there is also a dev server that runs on http://localhost:34115. Connect
to this in your browser, and you can call your Go code from devtools.

## Building

To build a redistributable, production mode package, use `wails build`.


## Testing

To test the core logic of the application run
```terminal
go test . 
```https://github.com/ferealqq/15GoFast/blob/main/documentation/weekly/week1.md

## Reports

[week 1](https://github.com/ferealqq/15GoFast/blob/main/documentation/weekly/week1.md)
[week 2](https://github.com/ferealqq/15GoFast/blob/main/documentation/weekly/week2.md)
