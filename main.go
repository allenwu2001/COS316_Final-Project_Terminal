package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"syscall"
)

func main() {
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("go-shell> ")
		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			continue
		}

		input = strings.TrimSuffix(input, "\n")
		input = strings.TrimSpace(input)

		// Split the input to separate the command and the arguments
		args := strings.Fields(input)

		if len(args) == 0 {
			continue
		}

		switch args[0] {
		case "exit":
			os.Exit(0)
		case "cd":
			// Change Directory
			if len(args) < 2 {
				fmt.Println("expected path")
			} else {
				err := os.Chdir(args[1])
				if err != nil {
					fmt.Fprintln(os.Stderr, err)
				}
			}
		case "setenv":
			// Set Environment Variable
			if len(args) != 3 {
				fmt.Println("expected usage: setenv KEY VALUE")
			} else {
				os.Setenv(args[1], args[2])
			}
		case "unsetenv":
			// Unset Environment Variable
			if len(args) != 2 {
				fmt.Println("expected usage: unsetenv KEY")
			} else {
				os.Unsetenv(args[1])
			}
		case "echo":
			echoCommand(args[1:])
		default:
			// Execute other commands
			executeCommand(args)
		}
	}
}

func echoCommand(args []string) {
	for i, arg := range args {
		if i > 0 {
			fmt.Print(" ")
		}
		fmt.Print(arg)
	}
	fmt.Println()
}

func executeCommand(args []string) {
	// Background process check
	background := false
	if args[len(args)-1] == "&" {
		background = true
		args = args[:len(args)-1] // Remove "&" from arguments
	}

	cmd := exec.Command(args[0], args[1:]...)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin

	if background {
		err := cmd.Start()
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
	} else {
		err := cmd.Run()
		if err != nil {
			// If it's a system call error, exit with the system call's exit status
			if exitError, ok := err.(*exec.ExitError); ok {
				if ws, ok := exitError.Sys().(syscall.WaitStatus); ok {
					os.Exit(ws.ExitStatus())
				}
			}
			fmt.Fprintln(os.Stderr, err)
		}
	}
}
