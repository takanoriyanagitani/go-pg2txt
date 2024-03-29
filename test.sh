#!/bin/sh

verbose=${ENV_VERBOSE}

go \
	test \
	-failfast \
	${verbose} \
	./...
