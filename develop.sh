#!/bin/bash
if [[ "$OSTYPE" == "linux-gnu"* ]]; then
    export SYSOS="linux"
elif [[ "$OSTYPE" == "darwin"* ]]; then
    export SYSOS="darwin"
elif [[ "$OSTYPE" == "cygwin" ]]; then
    export SYSOS="windows"
elif [[ "$OSTYPE" == "msys" ]]; then
    export SYSOS="windows"
elif [[ "$OSTYPE" == "win32" ]]; then
    export SYSOS="windows"
fi

if [[ "$(uname -m)" == "x86_64" ]]; then
    export SYSARCH="amd64"
elif [[ "$(uname -m)" == "i686" ]]; then
    export SYSARCH="386"
elif [[ "$(uname -m)" == "i386" ]]; then
    export SYSARCH="386"
elif [[ "$(uname -m)" == "aarch64" ]]; then
    export SYSARCH="arm64"
elif [[ "$(uname -m)" == "armv8l" ]]; then
    export SYSARCH="arm64"
elif [[ "$(uname -m)" == "armv8b" ]]; then
    export SYSARCH="arm64"
fi

if [[ "$GOOS" == "" ]]; then
    export GOOS="$SYSOS"
fi

if [[ "$GOARCH" == "" ]]; then
    export GOARCH="$SYSARCH"
fi

ln -s ./assets/liblcl_$GOOS\_$GOARCH.* ./

for filename in ./liblcl*; do 
    mv "$filename" "${filename//_$GOOS\_$GOARCH/}"
done

go run main.go
rm ./liblcl.*