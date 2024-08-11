# Variables
BINARY_NAME=ren
BUILD_DIR=build
LOG_FILE=ren.log
MAIN_FILE=cmd/rennen/main.go

# load .env file if it exists
ifneq (,$(wildcard ./.env))
    include .env
    export $(shell sed 's/=.*//' .env)
endif

# Default command is start
.PHONY: start
start: run

# Run the Go server
.PHONY: run
run:
	@echo "Starting the server..."
	@go run $(MAIN_FILE)

# Dist fake
.PHONY: dist-fake
dist-fake:
	@goreleaser release --snapshot --clean

.PHONY: dist
dist:
	@echo "Releasing the binary..."
# export the env GITHUB_TOKEN first before realeasing
	@goreleaser release

# Build the Go binary
.PHONY: build
build:
	$(eval path := $(BUILD_DIR)/$(BINARY_NAME))
	@echo "Building the server into $(path)..."

	@mkdir -p $(BUILD_DIR)
	@go build -ldflags="-s -w" -o $(path) $(MAIN_FILE)

	@echo "Making the binary executable..."
	@chmod +x $(path)

	@ # go is just too fucking fast building this
	@sleep 1

	@ # fullpath from root/ i/
	$(eval currentdir := $(shell pwd))
	$(eval fullpath := $(currentdir)/$(path))

	@echo "\n✨  Binary located at $(fullpath)"

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
