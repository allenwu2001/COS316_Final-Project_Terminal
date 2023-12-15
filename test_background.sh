#!/bin/bash

# Path to your Go shell executable
go_shell_executable="./go-shell" # Replace with the actual path to your Go shell executable

# Command to test background execution
background_command="sleep 5 &"
#background_command="sleep 5 "
test_command="echo Next command after sleep"

# Start the Go shell
$go_shell_executable << EOF
$background_command
$test_command
exit
EOF

echo "If this message displays immediately, background execution is successful."
