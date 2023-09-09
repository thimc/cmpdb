package main

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"strings"
	"testing"
)

func TestRemoveDuplicates(t *testing.T) {
	tests := []struct {
		name   string
		input  []string
		expect []string
	}{
		{
			name:   "Letters",
			input:  []string{"a", "a", "b", "b", "c", "c"},
			expect: []string{"a", "b", "c"},
		},
		{
			name:   "Numbers",
			input:  []string{"1", "", "2", "", "0", "3"},
			expect: []string{"1", "", "2", "0", "3"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			out := removeDuplicates(tt.input)
			if !reflect.DeepEqual(out, tt.expect) {
				t.Fatalf("expected \"%s\" got \"%s\"", tt.expect, out)
			}
		})
	}
}

func TestRemoveEmpty(t *testing.T) {
	test := struct {
		input  []string
		expect []string
	}{
		input:  []string{"h", "", "e", "", "", "", "l", "", "", "l", "", "o", "", "!"},
		expect: []string{"h", "e", "l", "l", "o", "!"},
	}

	out := removeEmpty(test.input)
	if !reflect.DeepEqual(out, test.expect) {
		t.Fatalf("expected \"%s\" got \"%s\"", test.expect, out)
	}
}
func TestRegexCompileC(t *testing.T) {
	tests := []struct {
		name   string
		log    []string
		expect bool
	}{
		{
			name:   "cc",
			log:    strings.Split("cc -c -std=c99 file.c", " "),
			expect: true,
		},
		{
			name:   "gcc",
			log:    strings.Split("gcc -c file.c", " "),
			expect: true,
		},
		{
			name:   "egcc",
			log:    strings.Split("egcc -c file.c", " "),
			expect: true,
		},
		{
			name:   "clang",
			log:    strings.Split("clang -c file.c", " "),
			expect: true,
		},
		{
			name:   "tcc",
			log:    strings.Split("tcc -c file.c", " "),
			expect: true,
		},
		{
			name:   "emcc",
			log:    strings.Split("emcc -c file.c", " "),
			expect: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if cc.Match([]byte(tt.log[0])) != tt.expect {
				t.Fatalf("expected \"%s\" to match", tt.log)
			}
		})
	}
}

func TestRegexCompileCpp(t *testing.T) {
	tests := []struct {
		name   string
		log    []string
		expect bool
	}{
		{
			name:   "c++",
			log:    strings.Split("c++ -c file.cpp", " "),
			expect: true,
		},
		{
			name:   "g++",
			log:    strings.Split("g++ -c file.cpp", " "),
			expect: true,
		},
		{
			name:   "clang++",
			log:    strings.Split("clang++ -c file.cpp", " "),
			expect: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if cpp.Match([]byte(tt.log[0])) != tt.expect {
				t.Fatalf("expected \"%s\" to match", tt.log)
			}
		})
	}
}

func TestRegexDir(t *testing.T) {
	expectedArrLen := 3
	tests := []struct {
		name string
		log  string
		mode string
	}{
		{
			name: "Enter gmake",
			mode: "Entering",
			log:  "gmake: Entering directory '/path/to/directory'",
		},
		{
			name: "Leave gmake",
			mode: "Leaving",
			log:  "gmake: Leaving directory '/path/to/directory'",
		},
		{
			name: "Enter make",
			mode: "Entering",
			log:  "make: Entering directory '/path/to/directory'",
		},
		{
			name: "Leave make",
			mode: "Leaving",
			log:  "make: Leaving directory '/path/to/directory'",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var (
				matches [][]byte
				line    = []byte(tt.log)
			)

			if dir.Match(line) {
				matches = dir.FindSubmatch(line)
			} else {
				t.Fatalf("expected \"%s\" to match", tt.log)
			}

			if len(matches) < expectedArrLen {
				t.Fatalf("expected length of submatches array to be %d, got %d", expectedArrLen, len(matches))
			}

			if string(matches[2]) != tt.mode {
				t.Fatalf("expected %s, got %s", string(matches[2]), tt.mode)
			}
		})
	}
}

func BenchmarkGenerate(b *testing.B) {
	wd, err := os.Getwd()
	if err != nil {
		b.Fatal(err)
	}
	workingDir := fmt.Sprintf("%s/test/", wd)
	for i := 0; i < b.N; i++ {
		// TODO: Find a better alternative rather than fork bombing the system
		reader, err := runMakeCmd(workingDir, []string{"-Bnkw"})
		if err != nil {
			b.Fatal(err)
		}

		defer reader.Close()
		db, err := scanProcOutput(workingDir, reader)
		if err != nil {
			b.Fatal(err)
		}
		_, err = json.MarshalIndent(db, "", *indentFlag)
		if err != nil {
			b.Fatal(err)
		}

	}
}
