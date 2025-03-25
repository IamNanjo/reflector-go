#!/usr/bin/env bash

NC='\033[0m'
YELLOW='\033[0;33m'
GREEN='\033[0;32m'

SCRIPT_PATH="$(dirname $(realpath $0))"

mkdir -p "$SCRIPT_PATH/dist"

go mod tidy &&

platforms=("linux/amd64" "linux/arm64" "windows/amd64" "windows/arm64" "darwin/amd64" "darwin/arm64")

echo -e "${YELLOW}Building project for ${#platforms[@]} platforms$NC"

for platform in "${platforms[@]}"
do
	platform_split=(${platform//\// })
	GOOS=${platform_split[0]}
	GOARCH=${platform_split[1]}
	
    output_name="$SCRIPT_PATH/dist/reflector-go-$GOOS-$GOARCH"

	if [ $GOOS = "windows" ]; then
		output_name+=".exe"
	fi

	printf "Building $output_name - "

	env GOOS=$GOOS GOARCH=$GOARCH go build -tags release -o $output_name -ldflags "-s -w"

	if [ $? -ne 0 ]; then
   		echo "failed"
		exit 1
	else
		echo -e "${GREEN}ok$NC"
	fi
done
