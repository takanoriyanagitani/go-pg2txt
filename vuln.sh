#!/bin/sh

verbose=${ENV_VERBOSE_VULN}

govulncheck ${verbose} ./...
