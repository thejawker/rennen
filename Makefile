# Variables
BINARY_NAME=ren
BUILD_DIR=build
LOG_FILE=ren_debug.log
MAIN_FILE=cmd/rennen/main.go

# Default command is start
.PHONY: start
start: run

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

# go is just too fucking fast building this
	@sleep 1

# fullpath from root
	$(eval currentdir := $(shell pwd))
	$(eval fullpath := $(currentdir)/$(path))

	@echo "\n✨  Binary located at $(fullpath)"

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


# Install the binary to a system-wide location
.PHONY: install
install: build
	@echo "Installing the binary..."
	@if [ "$(shell uname)" = "Linux" ] || [ "$(shell uname)" = "Darwin" ]; then \
		install_path="/usr/local/bin"; \
		sudo mkdir -p "$$install_path"; \
	elif [ "$(shell uname | grep -i mingw)" ]; then \
		install_path="/c/Program\ Files/$(BINARY_NAME)"; \
		mkdir -p "$$install_path"; \
	else \
		echo "Unsupported OS"; exit 1; \
	fi; \
	sudo install -m 755 $(BUILD_DIR)/$(BINARY_NAME) "$$install_path/$(BINARY_NAME)"; \
	echo "Installed $(BINARY_NAME) to $$install_path"

# Uninstall the binary from the system-wide location
.PHONY: uninstall
uninstall:
	@echo "Uninstalling the binary..."
	@if [ "$(shell uname)" = "Linux" ] || [ "$(shell uname)" = "Darwin" ]; then \
		install_path="/usr/local/bin"; \
	elif [ "$(shell uname | grep -i mingw)" ]; then \
		install_path="/c/Program\ Files/$(BINARY_NAME)"; \
	else \
		echo "Unsupported OS"; exit 1; \
	fi; \
	sudo rm -f "$$install_path/$(BINARY_NAME)"; \
	echo "Uninstalled $(BINARY_NAME) from $$install_path"
