# Variables
BINARY_NAME=ren
BUILD_DIR=build
LOG_FILE=ren_debug.log
MAIN_FILE=cmd/rennen/main.go

# Default command is start
.PHONY: start
start: stop run watch

# Run the Go server
.PHONY: run
run:
	@echo "Starting the server..."
	@go run $(MAIN_FILE)

# Build the Go binary
.PHONY: build
build:
	$(eval path := $(BUILD_DIR)/$(BINARY_NAME))
	@echo "Building the server into $(path)..."

	@mkdir -p $(BUILD_DIR)
	@go build -ldflags="-s -w" -o $(path) $(MAIN_FILE)

	@echo "Making the binary executable..."
	@chmod +x $(path)

# Watch for file changes and restart the server using reflex
.PHONY: watch
watch:
	@echo "Watching for changes with reflex..."
	@reflex -r '\.go$$' -- make restart

# Restart the server
.PHONY: restart
restart: stop run

# Stop the server
.PHONY: stop
stop:
	@echo "Stopping the server..."
	@pkill -f "$(MAIN_FILE)" || true

# Show server logs
.PHONY: logs
logs:
	@tail -f $(LOG_FILE)