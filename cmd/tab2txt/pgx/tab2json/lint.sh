#!/bin/sh

verbose=${ENV_VERBOSE}

gci() {
	golangci-lint \
		run \
		${verbose}
}

gci
