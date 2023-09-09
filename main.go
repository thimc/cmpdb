package main

import (
	"bufio"
	"container/list"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
	"regexp"
	"runtime"
	"strings"
)

var (
	cc    = regexp.MustCompile(`^.*-?[(g?|t|em)cc|clang]-?[0-9.]*$`)
	cpp   = regexp.MustCompile(`^.*-?[clang|gc|em]\+\+-?[0-9.]*$`)
	dir   = regexp.MustCompile(`^(g?make): (Entering|Leaving) directory ['"]([^'"]+)['"].*$`)
	shell = regexp.MustCompile(`[\$\(|\x60](.*?)[\)|\x60]`)

	commandFlag  = flag.Bool("c", false, "output the compilation command as a single string instead of a list of arguments")
	workDirFlag  = flag.String("d", "", "sets the working directory")
	fullPathFlag = flag.Bool("f", false, "expands the compiler executable path")
	indentFlag   = flag.String("i", "  ", "sets the json indentation")
	macroFlag    = flag.Bool("m", false, "adds the compilers predefined macros to the argument list")
	writeFlag    = flag.Bool("w", false, "writes the content \"compile_commands.json\" (default behaviour is to write to stdout)")

	makeFlags     = []string{"-Bknw"}
	jsonFile      = "compile_commands.json"
	make_entering = "Entering"
	make_leaving  = "Leaving"
)

func main() {
	flag.Parse()

	var workingDir string = *workDirFlag
	if workingDir == "" {
		wd, err := os.Getwd()
		if err != nil {
			fail(err)
		}
		workingDir = wd
	}

	reader, err := runMakeCmd(workingDir, makeFlags)
	if err != nil {
		fail(err)
	}
	defer reader.Close()
	db, err := scanProcOutput(workingDir, reader)
	if err != nil {
		fail(err)
	}

	data, err := json.MarshalIndent(db, "", *indentFlag)
	if err != nil {
		fail(err)
	}

	if *writeFlag {
		file, err := os.OpenFile(jsonFile, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
		if err != nil {
			fail(err)
		}
		defer file.Close()
		file.Write(data)
	} else {
		fmt.Printf("%s\n", string(data))
	}
}

func getMakeCmd() string {
	switch runtime.GOOS {
	case "openbsd":
		return "gmake"
	default:
		return "make"
	}
}

func runMakeCmd(baseDir string, args []string) (io.ReadCloser, error) {
	cmd := exec.Command(getMakeCmd(), strings.Join(args, " "))
	cmd.Dir = baseDir
	out, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}
	if err := cmd.Start(); err != nil {
		return nil, err
	}

	return out, nil
}

func scanProcOutput(baseDir string, reader io.ReadCloser) (*[]JSONFile, error) {
	var (
		db      = make([]JSONFile, 0)
		scanner = bufio.NewScanner(reader)
	)
	directories := list.New()
	directories.PushFront(baseDir)

	for scanner.Scan() {
		var (
			line   = scanner.Bytes()
			tokens = strings.Split(string(line), " ")
			entry  = JSONFile{
				Arguments: removeEmpty(tokens),
				File:      tokens[len(tokens)-1],
				Directory: directories.Front().Value.(string),
			}
		)
		if dir.Match(line) {
			groups := dir.FindSubmatch(line)
			if len(groups) > 3 {
				dirMode := string(groups[3])
				switch string(groups[2]) {
				case make_entering:
					directories.PushFront(dirMode)
				case make_leaving:
					directories.Remove(directories.Front())
				default:
					panic(fmt.Sprintf("Unhandled case: %s\n", string(line)))
				}
				entry.Directory = directories.Front().Value.(string)
			}

		} else if cc.Match([]byte(tokens[0])) || cpp.Match([]byte(tokens[0])) {
			if shell.Match(line) {
				groups := shell.FindAllStringSubmatch(string(line), 2)
				for _, group := range groups {
					args := strings.Split(group[1], " ")
					expandCmd, err := exec.Command(args[0], args[1:]...).Output()
					if err != nil {
						panic(err)
					}
					expand := string(expandCmd[:len(expandCmd)-1])
					tokenCmd := strings.Replace(strings.Join(tokens, " "), group[0], expand, -1)
					tokens = strings.Split(tokenCmd, " ")
				}

				if *fullPathFlag {
					exePath, err := exec.LookPath(tokens[0])
					if err == nil {
						tokens[0] = exePath
					}
				}

				if *macroFlag {
					macros, err := getCompilerMacros(tokens[0])
					if err != nil {
						//TODO
						continue
					}
					if len(macros) > 0 {
						tokens = append(tokens, macros...)
					}
				}

			}

			entry.Arguments = removeDuplicates(removeEmpty(tokens))

			if *commandFlag {
				entry.Command = strings.Join(entry.Arguments, " ")
				entry.Arguments = nil
			}

			if entry.File == "" {
				continue
			}
			_, err := os.Stat(baseDir + "/" + entry.File)
			if err != nil {
				continue
			}
			db = append(db, entry)
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	if len(db) < 1 {
		_, err := os.Stat("./Makefile")
		if err != nil {
			return nil, fmt.Errorf("Makefile not found")
		}
		return nil, fmt.Errorf("Couldn't parse the output from make")
	}

	return &db, nil
}

func getCompilerMacros(compiler string) ([]string, error) {
	var (
		out          []string
		compilerArgs = []string{"-x", "-c", "-dM", "-E", "-"}
	)

	macroCmd := exec.Command(compiler, compilerArgs...)
	outPipe, err := macroCmd.StdoutPipe()
	if err != nil {
		return out, err
	}
	if err := macroCmd.Start(); err != nil {
		return out, err
	}

	scanner := bufio.NewScanner(outPipe)
	for scanner.Scan() {
		tokens := strings.Split(scanner.Text(), " ")
		decl := fmt.Sprintf("-D%s=%s", tokens[1], strings.Join(tokens[2:], " "))
		out = append(out, decl)
	}
	if err := scanner.Err(); err != nil {
		return out, err
	}

	return out, nil
}

func removeEmpty(array []string) []string {
	var out []string
	for _, item := range array {
		if item != "" {
			out = append(out, item)
		}
	}

	return out
}

func removeDuplicates(array []string) []string {
	var out []string
	m := make(map[string]bool)
	for _, item := range array {
		if _, exist := m[item]; !exist {
			m[item] = true
			out = append(out, item)
		}
	}

	return out
}

func fail(err error) {
	fmt.Fprintf(os.Stderr, "error: %s\n", err)
	os.Exit(1)
}
