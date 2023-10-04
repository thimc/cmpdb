# cmpdb

cmpdb (**c**o**mp**ilation **d**ata**b**ase) is a go rewrite of [compiledb](https://github.com/nickdiego/compiledb).

This tool is suitable for make-based projects and is used to generate a
[JSON Compilation Database file](https://clang.llvm.org/docs/JSONCompilationDatabase.html)
for programs such as the popular LSP client, [clangd](https://clangd.llvm.org/).

## Usage
Before you invoke cmpdb, make sure your current working directory is the root
of your project.

cmpdb does not require any any flags. Granted, the output will be printed to
stdout rather than to a file, so to store the output you will need to redirect
the output of cmpdb like so `cmpdb > compile_commands.json`. Another solution
is to simply pass the `-w` flag which will write the file on disk.

If you wish to alter the behaviour of cmpdb, these flags are available:
```
-c  output the compilation command as a single string instead of a list of arguments
-d  sets the working directory
-f  expands the compiler executable path
-h  displays a help message
-i  sets the json indentation (default "  ")
-m  adds the compilers predefined macros to the argument list
-w  writes the content to "compile_commands.json" (default behaviour is to write to stdout)
```

## Installation

    go build -o cmpdb .
    cp ./cmpdb /usr/local/bin/

### Unit tests

    go test -v ./...

### Benchmarks:

    go test -v -bench=. -run=^#


Note: The benchmark is somewhat of a fork bomb as it spawns a new process
continuously in order to measure the total run time of what cmpdb's code
functionality does

## Bugs
There is no man page for this tool.
