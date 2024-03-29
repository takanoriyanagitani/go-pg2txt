#!/bin/sh

out="./cov.txt"
oht="./cov.html"

go \
	test \
	-race \
	-coverprofile="${out}" \
	-covermode=atomic \
	-failfast \
	./...

go \
	tool \
	cover \
	-html="${out}" \
	-o "${oht}"
