package main

import (
    "strings"
)

// lex takes an input string and breaks it into tokens, handling special characters and quoted strings.
func lex(input string) []string {
    var tokens []string
    var currentToken strings.Builder
    var inQuotes bool

    for _, char := range input {
        switch {
        case char == '"' && inQuotes: // End of quoted string
            inQuotes = false
            tokens = append(tokens, currentToken.String())
            currentToken.Reset()
        case char == '"' && !inQuotes: // Start of quoted string
            inQuotes = true
        case !inQuotes && (char == ' ' || char == '\t'): // Token separator
            if currentToken.Len() > 0 {
                tokens = append(tokens, currentToken.String())
                currentToken.Reset()
            }
        default: // Regular character
            currentToken.WriteRune(char)
        }
    }

    if currentToken.Len() > 0 {
        tokens = append(tokens, currentToken.String())
    }

    return tokens
}