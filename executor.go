package main

import (
    "fmt"
    "os"
    "os/exec"
)

// execute decides whether to run a built-in or external command.
func execute(cmd *Command) {
    if cmd == nil {
        return
    }

    switch cmd.Name {
    case "cd", "setenv", "unsetenv", "exit":
        executeBuiltinCommand(cmd)
    default:
        go executeExternalCommand(cmd) // Goroutine for concurrent execution
    }
}

// executeBuiltinCommand executes built-in shell commands.
func executeBuiltinCommand(cmd *Command) {
	switch cmd.Name {
	case "cd":
			if len(cmd.Args) < 1 {
					fmt.Fprintln(os.Stderr, "cd: missing argument")
					return
			}
			err := os.Chdir(cmd.Args[0])
			if err != nil {
					fmt.Fprintln(os.Stderr, "cd:", err)
			}
	case "setenv":
			if len(cmd.Args) < 2 {
					fmt.Fprintln(os.Stderr, "setenv: missing argument")
					return
			}
			os.Setenv(cmd.Args[0], cmd.Args[1])
	case "unsetenv":
			if len(cmd.Args) < 1 {
					fmt.Fprintln(os.Stderr, "unsetenv: missing argument")
					return
			}
			os.Unsetenv(cmd.Args[0])
	case "exit":
			os.Exit(0)
	default:
			fmt.Fprintf(os.Stderr, "Unknown built-in command: %s\n", cmd.Name)
	}
}

// executeExternalCommand executes external commands and scripts.
func executeExternalCommand(cmd *Command) {
	command := exec.Command(cmd.Name, cmd.Args...)

	// Handle input redirection
	if cmd.InputFile != "" {
			inputFile, err := os.Open(cmd.InputFile)
			if err != nil {
					fmt.Fprintln(os.Stderr, "Failed to open input file:", err)
					return
			}
			defer inputFile.Close()
			command.Stdin = inputFile
	}

	// Handle output redirection
	if cmd.OutputFile != "" {
			outputFile, err := os.Create(cmd.OutputFile)
			if err != nil {
					fmt.Fprintln(os.Stderr, "Failed to open output file:", err)
					return
			}
			defer outputFile.Close()
			command.Stdout = outputFile
	} else {
			command.Stdout = os.Stdout
	}

	command.Stderr = os.Stderr

	if cmd.Background {
			// Execute the command in the background
			go func() {
					err := command.Run()
					if err != nil {
							fmt.Fprintln(os.Stderr, err)
					}
			}()
	} else {
			// Execute the command and wait for it to finish
			err := command.Run()
			if err != nil {
					fmt.Fprintln(os.Stderr, err)
			}
	}
}