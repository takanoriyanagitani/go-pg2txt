#!/bin/sh

which pkgsite | fgrep -q pkgsite || exec sh -c 'echo pkgsite missing.; exit 1'

addr=0.0.0.0
port=3298

pkgsite \
	-http "${addr}:${port}"
