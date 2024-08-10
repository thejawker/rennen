#!/bin/bash

# Array of ANSI color codes
colors=(31 32 33 34 35 36 37 91 92 93 94 95 96 97)

generate_random_log() {
    local types=("INFO" "WARNING" "ERROR" "DEBUG")
    local messages=(
        "Processing request"
        "Database query executed"
        "API call successful"
        "Cache miss"
        "File not found"
        "Connection timeout"
        "Authentication failed"
        "Data validation error"
        "Memory usage high"
        "Task completed successfully"
    )

    local color=${colors[$RANDOM % ${#colors[@]}]}
    local type=${types[$RANDOM % ${#types[@]}]}
    local message=${messages[$RANDOM % ${#messages[@]}]}
    local timestamp=$(date +"%Y-%m-%d %H:%M:%S")

#    echo "[$timestamp] [$type] $message (PID: $$)"
    echo "[$timestamp] [$type] $message (PID: $$)" | awk -v color=$color '{print "\033[0;"color"m" $0 "\033[0m"}'
}

while true; do
    generate_random_log
    sleep $(awk 'BEGIN{print 0.1 + (2 - 0.1) * rand()}')
# faster
#    sleep $(awk 'BEGIN{print 0.1 + (0.1 - 0.1) * rand()}')
done