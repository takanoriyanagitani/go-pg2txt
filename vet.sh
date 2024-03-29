#!/bin/sh

verbose=${ENV_VERBOSE}

go \
	vet \
	-race \
	"${verbose}" \
	./...
