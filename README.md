# cmpdb

cmpdb (**c**o**mp**ilation **d**ata**b**ase) is a go rewrite of [compiledb](https://github.com/nickdiego/compiledb).

This tool is suitable for make-based projects and is used to generate a JSON
Compilation Database file[1] for programs such as the popular LSP client, clangd[2].

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
To install and run the unit tests, `make` is required. Simply invoke `make test`
to run the unit tests and `make install` to install it to your system. You will
need to have root access to install the program.

Note: The benchmark is somewhat of a fork bomb as it calls `make` continuously
in order to measure the total run time of what cmpdb's code functionality does.

## Bugs
There is no man page for this tool.

## References
1: [LLVM.org - JSON Compilation Database Format Specification](https://clang.llvm.org/docs/JSONCompilationDatabase.html)
2: [LLVM.org - Clangd](https://clangd.llvm.org/)
