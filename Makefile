# ---------------------------------------------------------------------------
# Makefile – Financial Transaction Tracker
# ---------------------------------------------------------------------------
# CGO_ENABLED=1 is required because the SQLite driver (mattn/go-sqlite3)
# is a C extension that must be compiled via cgo.  Make sure gcc and
# libsqlite3-dev are installed on the build machine:
#
#   sudo apt-get install -y gcc libsqlite3-dev
# ---------------------------------------------------------------------------

BINARY      := bin/finance-api
CMD         := cmd/main.go
LDFLAGS     := -ldflags="-s -w"

.PHONY: all build run clean tidy

all: build

## build: compile a native Linux amd64 binary (requires CGO + gcc)
build:
	CGO_ENABLED=1 GOOS=linux GOARCH=amd64 \
		go build $(LDFLAGS) -o $(BINARY) $(CMD)
	@echo "Binary written to $(BINARY)"

## run: run the application locally
run:
	CGO_ENABLED=1 go run $(CMD)

## clean: remove the compiled binary
clean:
	rm -f $(BINARY)

## tidy: tidy and verify Go modules
tidy:
	go mod tidy
	go mod verify
