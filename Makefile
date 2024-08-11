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
	@current_version=$$(cat VERSION); \
	echo "Current version: $$current_version"; \
	IFS='.' read -r major minor patch <<< "$${current_version#v}"; \
	read -p "Is this a patch, minor, or major release? " level; \
	case $$level in \
		patch) patch=$$(($$patch + 1));; \
		minor) minor=$$(($$minor + 1)); patch=0;; \
		major) major=$$(($$major + 1)); minor=0; patch=0;; \
		*) echo "Invalid level, use patch/minor/major"; exit 1;; \
	esac; \
	new_version="v$$major.$$minor.$$patch"; \
	echo "$$new_version" > VERSION; \
	echo "Updated VERSION file to $$new_version"; \
	git add VERSION; \
	git commit -m "Bump version to $$new_version"; \
	read -p "Enter the tag message: " tag_message; \
	git tag -a $$new_version -m "$$tag_message"; \
	git push --follow-tags; \
	echo "Releasing the binary..."; \
	goreleaser release --clean

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

# Unbrew will remove the brew installed version for testing
.PHONY: unbrew
unbrew:
	@brew uninstall thejawker/tappen/rennen

# brew will pull the bin from homebrew
.PHONY: brew
brew:
	@brew install thejawker/tappen/rennen

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
