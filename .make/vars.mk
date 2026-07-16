# .make/vars.mk - Common variables
# Included by main Makefile

# Go commands
GO := go
GOBUILD := $(GO) build
GOTEST := $(GO) test
GOFMT := $(GO) fmt
GOVET := $(GO) vet
GOMOD := $(GO) mod

# Test settings
COVERAGE_OUT := coverage.out
COVERAGE_HTML := coverage.html
TEST_TIMEOUT := 5m
RACE_FLAG := -race
