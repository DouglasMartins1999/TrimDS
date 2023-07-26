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

APPNAME=trimDS
APPLOCATION=./dist/$APPNAME.$GOARCH
SYSOPATH=./syso/res_windows_$GOARCH.syso

if [[ "$GOOS" == "darwin" ]]; then
    export CGO_ENABLED=1 # PRe or Pos sudo?
    sudo go build -ldflags="-s -w" -o bin

    rm -rf $APPLOCATION.app
    mkdir $APPLOCATION.app
    mkdir $APPLOCATION.app/Contents
    mkdir $APPLOCATION.app/Contents/MacOS
    mkdir $APPLOCATION.app/Contents/Resources

    mv ./bin $APPLOCATION.app/Contents/MacOS/$APPNAME
    cp ./assets/info.plist $APPLOCATION.app/Contents
    cp ./assets/icon.icns $APPLOCATION.app/Contents/Resources/$APPNAME.icns
    cp ./assets/liblcl_$GOOS\_$GOARCH.dylib $APPLOCATION.app/Contents/MacOS/liblcl.dylib
    sed -i'' -e "s/__appname__/$APPNAME/g" $APPLOCATION.app/Contents/info.plist
    rm $APPLOCATION.app/Contents/info.plist-e

    chmod +x $APPLOCATION.app
else
    if [[ "$GOOS" == "windows" ]]; then
        if [ ! -f "$SYSOPATH" ]; then
            GOOS=$SYSOS GOARCH=$SYSARCH go build github.com/akavel/rsrc
            ./rsrc -ico assets/icon.ico -manifest assets/app.manifest -arch $GOARCH -o $SYSOPATH
            rm ./rsrc
        fi

        go build -ldflags="-s -w -H windowsgui" -o $APPNAME.exe
        zip -q $APPLOCATION.zip $APPNAME.exe ./assets/liblcl_$GOOS\_$GOARCH.dll
        rm $APPNAME.exe
    else
        export CGO_ENABLED=1
        go build -ldflags="-s -w" -o $APPNAME
        tar -czf $APPLOCATION.tar.gz $APPNAME ./assets/liblcl_$GOOS\_$GOARCH.so
        rm $APPNAME
    fi
fi