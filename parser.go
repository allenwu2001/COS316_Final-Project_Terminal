package main

// Command represents a parsed command with its arguments and redirection info.
type Command struct {
    Name        string
    Args        []string
    InputFile   string
    OutputFile  string
    Background  bool
}

// parse interprets a slice of tokens into a Command struct.
func parse(tokens []string) *Command {
    if len(tokens) == 0 {
        return nil
    }

    cmd := &Command{}
    for i, token := range tokens {
        switch {
        case token == "<" && i+1 < len(tokens): // Input redirection
            cmd.InputFile = tokens[i+1]
        case token == ">" && i+1 < len(tokens): // Output redirection
            cmd.OutputFile = tokens[i+1]
        case token == "&": // Background execution
            cmd.Background = true
        default: // Command or argument
            if cmd.Name == "" {
                cmd.Name = token
            } else if cmd.InputFile == "" && cmd.OutputFile == "" {
                cmd.Args = append(cmd.Args, token)
            }
        }
    }

    return cmd
}