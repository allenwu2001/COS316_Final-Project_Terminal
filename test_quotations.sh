#!/bin/bash

# Path to your Go shell executable
GO_SHELL="./main"

# Function to test a command in your Go shell
test_command() {
    command=$1
    expected_output=$2

    echo "Testing command: $command"
    actual_output=$(echo "$command" < $GO_SHELL)

    if [[ "$actual_output" == "$expected_output" ]]; then
        echo "Test Passed"
    else
        echo "Test Failed: Expected '$expected_output', got '$actual_output'"
    fi
}

# Test case: Command with quotes
# test_command "echo \"Hello, world\"" "Hello, world"

# Test case: Command without quotes
#test_command "echo Hello world" "Hello world"

# Add more test cases as needed
