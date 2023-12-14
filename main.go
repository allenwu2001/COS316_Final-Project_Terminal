package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"
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
		args := tokenize(input)

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
		// case "echo":
		// echoCommand(args[1:])
		default:
			// Execute other commands
			executeCommand(args)
		}
	}
}

// Updated tokenize function
func tokenize(input string) []string {
	var tokens []string
	var currentToken strings.Builder
	inQuotes := false

	for _, char := range input {
		switch {
		case char == '"' && inQuotes:
			inQuotes = false
			tokens = append(tokens, currentToken.String())
			currentToken.Reset()
		case char == '"' && !inQuotes:
			inQuotes = true
		case (char == '>' || char == '<') && !inQuotes:
			if currentToken.Len() > 0 {
				tokens = append(tokens, currentToken.String())
				currentToken.Reset()
			}
			tokens = append(tokens, string(char))
		case char == ' ' && !inQuotes:
			if currentToken.Len() > 0 {
				tokens = append(tokens, currentToken.String())
				currentToken.Reset()
			}
		default:
			currentToken.WriteRune(char)
		}
	}

	if currentToken.Len() > 0 {
		tokens = append(tokens, currentToken.String())
	}

	return tokens
}

func isFileName(s string) bool {
	fileInfo, err := os.Stat(s)
	if os.IsNotExist(err) {
		return false // File does not exist
	}
	return !fileInfo.IsDir() // Return true if it's not a directory
}

func executeCommand(args []string) {
	// Check for input and output redirection
	var inputFile, outputFile string
	for i, arg := range args {
		if arg == "<" {
			if i+1 < len(args) {
				inputFile = args[i+1]
				if !isFileName(inputFile) {
					fmt.Println("File does not exist")
					return
				}
				args = append(args[:i], args[i+2:]...)
				break
			}
		} else if arg == ">" {
			if i+1 < len(args) {
				outputFile = args[i+1]
				if !(strings.HasSuffix(outputFile, ".txt") || strings.HasSuffix(outputFile, ".md")) {
					fmt.Println("Invalid file type")
					return
				}
				args = append(args[:i], args[i+2:]...)
				break
			}
		}
	}

	cmd := exec.Command(args[0], args[1:]...)

	// Set up input redirection
	if inputFile != "" {
		inFile, err := os.Open(inputFile)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			return
		}
		defer inFile.Close()
		cmd.Stdin = inFile
	} else {
		cmd.Stdin = os.Stdin
	}

	// Set up output redirection
	if outputFile != "" {
		outFile, err := os.Create(outputFile)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			return
		}
		defer outFile.Close()
		cmd.Stdout = outFile
	} else {
		cmd.Stdout = os.Stdout
	}

	cmd.Stderr = os.Stderr

	// Execute the command
	if err := cmd.Run(); err != nil {
		// Handle errors...
	}
}
