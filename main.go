package main

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/chzyer/readline"
)

func main() {
	completer := readline.NewPrefixCompleter(
		readline.PcItem("exit"),
		readline.PcItem("cd"),
		readline.PcItem("setenv"),
		readline.PcItem("unsetenv"),
		readline.PcItem("sleep"),
	)

	rl, err := readline.NewEx(&readline.Config{
		Prompt:          "\033[31mgo-shell> \033[0m",
		AutoComplete:    completer,
		InterruptPrompt: "^C",
		EOFPrompt:       "exit",
	})

	if err != nil {
		panic(err)
	}
	defer rl.Close()

	for {
		line, err := rl.Readline()
		if err != nil {
			break
		}

		input := strings.TrimSpace(line)
		args := tokenize(input)

		if len(args) == 0 {
			continue
		}

		switch args[0] {
		case "exit":
			os.Exit(0)
		case "cd":
			if len(args) < 2 {
				fmt.Println("need path")
			} else {
				err := os.Chdir(args[1])
				if err != nil {
					fmt.Fprintln(os.Stderr, err)
				}
			}
		case "setenv":
			if len(args) != 3 {
				fmt.Println("Invalid input for setenv")
			} else {
				os.Setenv(args[1], args[2])
			}
		case "unsetenv":
			if len(args) != 2 {
				fmt.Println("Invalid input for unsetenv")
			} else {
				os.Unsetenv(args[1])
			}
		case "sleep":
			if len(args) != 2 {
				fmt.Println("Sleeping for " + args[1] + " seconds in background")
			} else {
				duration, err := strconv.Atoi(args[1])
				if err != nil {
					fmt.Fprintln(os.Stderr, "Invalid duration:", err)
					break
				}
				time.Sleep(time.Duration(duration) * time.Second)
			}
		default:
			executeCommand(args)
		}
	}
}

// Tokenize the input
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
		return false
	}
	return !fileInfo.IsDir()
}

func executeCommand(args []string) {
	var background bool = false
	if len(args) > 0 && args[len(args)-1] == "&" {
		background = true
		args = args[:len(args)-1]
	}

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

	if background {
		runInBackground(cmd)
	} else {
		if err := cmd.Run(); err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
	}
}

func runInBackground(cmd *exec.Cmd) {
	if err := cmd.Start(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}
	go func() {
		if err := cmd.Wait(); err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
	}()
}
