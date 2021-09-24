package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/fatih/color"
)

type Test struct {
	name   string //name of the test (the directory it is located in)
	path   string //path to the .ddp file to be run
	output string //the output that is expected
	input  string //the input passed to the test
}

func main() {
	defer fmt.Scanln()
	ddpImplementation := os.Args[1]
	ex, err := os.Executable()
	if err != nil {
		log.Fatal()
	}
	path := filepath.Dir(ex)

	directory, err := os.ReadDir(path + "/tests")
	if err != nil {
		log.Println("tests are not present: " + err.Error())
		return
	}
	tests := make([]Test, len(directory))
	for i, dir := range directory {
		if dir.IsDir() {
			out, err := ioutil.ReadFile(path + "/tests/" + dir.Name() + "/" + dir.Name() + "_output.txt")
			if err != nil {
				out = []byte("")
			}
			in, err := ioutil.ReadFile(path + "/tests/" + dir.Name() + "/" + dir.Name() + "_output.txt")
			if err != nil {
				in = []byte("")
			}
			tests[i] = Test{name: dir.Name(), path: path + "/tests/" + dir.Name() + "/" + dir.Name() + ".ddp", output: string(out), input: string(in)}
		}
	}

	failed := 0
	for _, test := range tests {
		cmd := exec.Command(ddpImplementation, test.path)
		cmd.Dir = filepath.Dir(test.path)
		var (
			stdout bytes.Buffer
			stderr bytes.Buffer
		)
		cmd.Stdout = bufio.NewWriter(&stdout)
		cmd.Stderr = bufio.NewWriter(&stderr)
		cmd.Stdin = strings.NewReader(test.input)

		err := cmd.Run()
		if err != nil || len(stderr.String()) != 0 || stdout.String() != test.output {
			color.Set(color.FgRed)
			fmt.Println("Error running the test '" + test.name + "':\n")
			fmt.Println("Stderr:\n" + stderr.String())
			fmt.Println("Expected output:\n" + test.output)
			fmt.Println("Actual output:")

			runes := []rune(test.output)
			for i, r := range []rune(stdout.String()) {
				if i < len(runes) {
					if r == runes[i] {
						fmt.Print(string(r))
					} else {
						color.Set(color.FgCyan)
						fmt.Print(string(r))
						color.Set(color.FgRed)
					}
				} else {
					fmt.Print(string(r))
				}
			}

			color.Unset()
			failed++
		} else {
			color.Set(color.FgGreen)
			fmt.Println("Test '" + test.name + "' ran successful!\n")
			color.Unset()
		}
	}

	if failed > 0 {
		color.Set(color.FgRed)
	} else {
		color.Set(color.FgGreen)
	}
	fmt.Println(fmt.Sprint(len(tests)-failed) + " out of " + fmt.Sprint(len(tests)) + " succeeded!")
	color.Unset()
}
