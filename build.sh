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
elif [[ "$(uname -m)" == "arm64" ]]; then
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

APPNAME=TrimDS
APPLOCATION=./dist/$APPNAME.$GOARCH
SYSOPATH=./assets/res_windows_$GOARCH.syso

if [[ "$GOOS" == "darwin" ]]; then
    export CGO_ENABLED=1
    go build -ldflags="-s -w" -o bin

    rm -rf $APPNAME.app
    mkdir $APPNAME.app
    mkdir $APPNAME.app/Contents
    mkdir $APPNAME.app/Contents/MacOS
    mkdir $APPNAME.app/Contents/Resources

    mv ./bin $APPNAME.app/Contents/MacOS/$APPNAME
    cp ./assets/info.plist $APPNAME.app/Contents
    cp ./assets/icon.icns $APPNAME.app/Contents/Resources/$APPNAME.icns
    cp ./assets/liblcl_${GOOS}_$GOARCH.dylib $APPNAME.app/Contents/MacOS/liblcl.dylib
    sed -i'' -e "s/__appname__/$APPNAME/g" $APPNAME.app/Contents/info.plist
    rm $APPNAME.app/Contents/info.plist-e

    chmod +x $APPNAME.app
    
    rm -rf $APPLOCATION
    mkdir -p $APPLOCATION
    mv $APPNAME.app $APPLOCATION
    hdiutil create "$APPLOCATION.tmp.dmg" -ov -volname "$APPNAME Install" -fs HFS+ -srcfolder "$APPLOCATION" &> /dev/null
    hdiutil convert "$APPLOCATION.tmp.dmg" -format UDZO -o "$APPLOCATION.dmg" &> /dev/null
    rm "$APPLOCATION.tmp.dmg"

else
    if [[ "$1" == "libres" ]]; then
        export LIBRES="-tags tempdll"
        APPNAME=$APPLOCATION

        LIBNAME=liblclbinres
        LIBVER=$(cat go.mod | grep $LIBNAME | grep -Eo 'v\d\S+')

        if [[ "$GOPATH" != "" ]]; then
            BINPATH="$GOPATH"
        else
            BINPATH="$HOME/go"
        fi

        LIBPATH=$BINPATH/pkg/mod/github.com/ying32/$LIBNAME@$LIBVER

        if [[ ! -f "$LIBPATH/liblcl_${GOOS}_$GOARCH.go" ]]; then
            if [[ "$GOOS" == "windows" ]]; then
                LIBEXT="dll"
            else
                LIBEXT="so"
            fi

            GOOS=$SYSOS GOARCH=$SYSARCH go run utils/main.go ./assets/liblcl_${GOOS}_$GOARCH.$LIBEXT ./liblcl_${GOOS}_$GOARCH.go
            sudo mv ./liblcl_${GOOS}_$GOARCH.go $LIBPATH
            go clean -cache
        fi
    fi

    if [[ "$GOOS" == "windows" ]]; then
        if [ ! -f "$SYSOPATH" ]; then
            GOOS=$SYSOS GOARCH=$SYSARCH go build github.com/akavel/rsrc
            ./rsrc -ico assets/icon.ico -manifest assets/app.manifest -arch $GOARCH -o $SYSOPATH
            rm ./rsrc
        fi

        go build $LIBRES -ldflags="-s -w -H windowsgui" -o $APPNAME.exe

        if [[ "$LIBRES" == "" ]]; then
            ln -s ./assets/liblcl_${GOOS}_$GOARCH.dll ./liblcl.dll
            zip -q $APPLOCATION.zip $APPNAME.exe liblcl.dll
            rm $APPNAME.exe liblcl.dll
        fi
    else
        export CGO_ENABLED=1
        go build $LIBRES -ldflags="-s -w" -o $APPNAME

        if [[ "$LIBRES" == "" ]]; then
            ln -s ./assets/liblcl_${GOOS}_$GOARCH.so ./liblcl.so
            tar -czhf $APPLOCATION.tar.gz $APPNAME liblcl.so
            rm $APPNAME liblcl.so
        fi
    fi
fi