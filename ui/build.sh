#!/bin/sh

cd $(dirname $0)

rm -rf build/

npm run build

mkdir -p build/css/
cp src/css/* build/css/