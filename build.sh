#!/bin/bash

if [[ "$GOOS" == "" ]]; then
    if [[ "$OSTYPE" == "linux-gnu"* ]]; then
        export GOOS="linux"
        export CGO_ENABLED=1
    elif [[ "$OSTYPE" == "darwin"* ]]; then
        export GOOS="darwin"
        export CGO_ENABLED=1
    elif [[ "$OSTYPE" == "cygwin" ]]; then
        export GOOS="windows"
    elif [[ "$OSTYPE" == "msys" ]]; then
        export GOOS="windows"
    elif [[ "$OSTYPE" == "win32" ]]; then
        export GOOS="windows"
    fi
fi

if [[ "$GOARCH" == "" ]]; then
    if [[ "$(uname -m)" == "x86_64" ]]; then
        export GOARCH="amd64"
    elif [[ "$(uname -m)" == "i686" ]]; then
        export GOOS="386"
    elif [[ "$(uname -m)" == "i386" ]]; then
        export GOOS="386"
    elif [[ "$(uname -m)" == "aarch64" ]]; then
        export GOOS="arm64"
    elif [[ "$(uname -m)" == "armv8l" ]]; then
        export GOOS="arm64"
    elif [[ "$(uname -m)" == "armv8b" ]]; then
        export GOOS="arm64"
    fi
fi

APPNAME=trimDS
APPLOCATION=./dist/$APPNAME.$GOARCH

if [[ "$GOOS" == "darwin" ]]; then
    sudo go build -ldflags="-s -w" -o bin

    rm -rf $APPLOCATION.app
    mkdir $APPLOCATION.app
    mkdir $APPLOCATION.app/Contents
    mkdir $APPLOCATION.app/Contents/MacOS
    mkdir $APPLOCATION.app/Contents/Resources

    mv ./bin $APPLOCATION.app/Contents/MacOS/$APPNAME
    cp ./assets/info.plist $APPLOCATION.app/Contents
    cp ./assets/icon.icns $APPLOCATION.app/Contents/Resources/$APPNAME.icns
    cp ./assets/liblcl_$GOOS_$GOARCH.dylib $APPLOCATION.app/Contents/MacOS/liblcl.dylib
    sed -i'' -e "s/__appname__/$APPNAME/g" $APPLOCATION.app/Contents/info.plist
    rm $APPLOCATION.app/Contents/info.plist-e

    chmod +x $APPLOCATION.app
else
    echo $APPNAME
fi