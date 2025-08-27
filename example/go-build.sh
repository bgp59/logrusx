#! /bin/bash --noprofile

set -e
case "$0" in
    /*|*/*) this_dir=$(cd $(dirname $0) && pwd);;
    *) this_dir=$(cd $(dirname $(which $0)) && pwd);;
esac

set -x
cd $this_dir
mkdir -p bin
go build -o bin/logger-example

  