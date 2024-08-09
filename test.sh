#!/bin/bash

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

    local type=${types[$RANDOM % ${#types[@]}]}
    local message=${messages[$RANDOM % ${#messages[@]}]}
    local timestamp=$(date +"%Y-%m-%d %H:%M:%S")

    echo "[$timestamp] [$type] $message (PID: $$)"
}

while true; do
    generate_random_log
    sleep $(awk 'BEGIN{print 0.1 + (2 - 0.1) * rand()}')
done