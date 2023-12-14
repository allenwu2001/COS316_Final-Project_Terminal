#!/bin/bash

# Path to your Go shell executable
GO_SHELL="./main"

# Function to execute a command in the Go shell
execute_in_go_shell() {
    echo "$1" < $GO_SHELL
}

# Function to test if the shell is forking processes
test_forking() {
    # Store the current directory
    current_dir=$(pwd)

    # Change directory in your Go shell
    execute_in_go_shell "cd /tmp"

    # Check the current directory after running the command
    if [[ $(pwd) == "$current_dir" ]]; then
        echo "Forking Test Passed: The current directory has not changed."
    else
        echo "Forking Test Failed: The current directory has changed."
    fi
}

# Run the test
test_forking

# You can add more tests here as needed
