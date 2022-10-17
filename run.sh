#!/bin/bash

set -e

scriptDir="$(cd "$(dirname "$0")" && pwd)"
cd "${scriptDir}"

# get env from .env
export $(grep -v '^#' .env | xargs)

go run main.go