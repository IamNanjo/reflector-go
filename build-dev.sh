#!/usr/bin/env bash

NC='\033[0m'
YELLOW='\033[0;33m'
GREEN='\033[0;32m'
RED='\033[0;31m'

SCRIPT_PATH="$(dirname $(realpath $0))"

mkdir -p "$SCRIPT_PATH/dist"

go mod tidy &&

echo -e "${YELLOW}Building go packages$NC"
go build -tags debug -o "$SCRIPT_PATH/dist" -v &&
echo -e "${GREEN}All built$NC"

dist/reflector-go -save "$SCRIPT_PATH/dist/test"
