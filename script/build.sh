#!/bin/bash

cd "$( dirname "${BASH_SOURCE[0]}" )/.."

ROOT=`pwd`
RELEASE_DIR="$ROOT/release"
RELEASE_BIN_DIR="$RELEASE_DIR/bin"

version=$(if [ "`git describe --tags --abbrev=0 2>/dev/null`" != "" ];then git describe --tags --abbrev=0; else git log --pretty=format:'%h' -n 1; fi)
prog=$(grep "^module .*$" go.mod | sed -r "s/^module .*\/(.*)$/\1/")

function pre_build() {
    # clean dist, release
    rm -rf dist release
    mkdir -p $RELEASE_BIN_DIR
}

function build_frontend() {
    yarn build
}

function build() {
    local GOOS=$1
    local GOARCH=$2
    local PLATFORM=$3

    name=$prog-$PLATFORM-$version
    
    echo "building $name"

    CGO_ENABLED=0 GOOS=$GOOS GOARCH=$GOARCH go build -o $RELEASE_BIN_DIR || exit 1
    
    archive="$name.tar.gz"

    if [[ $1 == "windows" ]]; then
        tar -zcf "$RELEASE_DIR/$archive" -T package.list -C $RELEASE_BIN_DIR "$prog.exe"
        rm -f "$RELEASE_BIN_DIR/$prog.exe"
    else
        tar -zcf "$RELEASE_DIR/$archive" -T package.list -C $RELEASE_BIN_DIR $prog
        rm -f "$RELEASE_BIN_DIR/$prog"
    fi
}

function post_build() {
    rm -rf $RELEASE_BIN_DIR
}

pre_build

build_frontend
build windows amd64 win64
build linux amd64 linux64

post_build