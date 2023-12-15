#!/bin/bash

# Path to your Go shell executable
GO_SHELL_EXECUTABLE="./main"

# Function to test output redirection
test_output_redirection() {
    test_command=$1
    expected_output=$2
    output_file=$3

    # Create a temporary file for commands
    temp_command_file=$(mktemp)
    
    # Write the test command to the temporary file
    echo "$test_command" > "$temp_command_file"
    
    # Execute your Go shell with the temporary file as input
    cat "$temp_command_file" | $GO_SHELL_EXECUTABLE

    # Check if the output file was created and contains expected output
    if [[ -f "$output_file" ]]; then
        actual_output=$(cat "$output_file")
        if [[ "$actual_output" == "$expected_output" ]]; then
            echo "Output Redirection Test Passed"
        else
            echo "Output Redirection Test Failed: Expected '$expected_output', got '$actual_output'"
        fi
        rm "$output_file"  # Clean up the output file
    else
        echo "Output Redirection Test Failed: Output file '$output_file' not created"
    fi

    rm "$temp_command_file"  # Clean up the temporary command file
}

# Test case: Output redirection
test_output_redirection "echo Hello, World! > test_output.txt" "Hello, World!" "test_output.txt"

# You can add more test cases here as needed
