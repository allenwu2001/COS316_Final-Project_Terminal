package main

import (
    "bufio"
    "fmt"
    "os"
    "strings"
)

func main_3step() {
    reader := bufio.NewReader(os.Stdin)

    for {
        fmt.Print("GoShell> ")
        input, _ := reader.ReadString('\n')
        input = strings.TrimSpace(input)

        if input == "exit" {
            break
        }

        tokens := lex(input)
        cmd := parse(tokens)
        execute(cmd)
    }
}