#!/bin/bash
# list_pkgs <directory>
# Lists all Go packages in the given directory.

directory="$1"
if [ -z "$directory" ]; then
    directory="."
fi
go list $directory/... | sed -e "s/^github.com\/alpine-hodler\/gidari/./"

