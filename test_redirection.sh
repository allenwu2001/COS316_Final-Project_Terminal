#!/bin/bash

# Function to test a command in the Go shell
test_command() {
    command=$1
    expected_output=$2
    output_file=$3

    echo "Testing command: $command"
    echo $command < ./go-shell

    # Check if output matches the expectation
    if [[ -f "$output_file" ]]; then
        actual_output=$(cat "$output_file")
        if [[ "$actual_output" == "$expected_output" ]]; then
            echo "Test Passed"
        else
            echo "Test Failed: Expected '$expected_output', got '$actual_output'"
        fi
    else
        echo "Test Failed: Output file '$output_file' not found"
    fi
}

# Test case 1: Output redirection
test_command "echo hello world > test_output.txt" "hello world" "test_output.txt"

# Test case 2: Input redirection
# First create a test input file
echo "input text" > test_input.txt
test_command "cat < test_input.txt" "input text" "test_output.txt"

# Clean up test files
rm test_input.txt test_output.txt

# Add more test cases as needed
