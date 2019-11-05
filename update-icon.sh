#!/bin/sh

INPUT=$(pwd)/Luxafor.app/Contents/Resources/icon.png
OUTPUT=$(dirname $INPUT)/icon.icns
WORKDIR=$(mktemp -d)
ICONSET=$WORKDIR/icon.iconset

mkdir $ICONSET

# Normal screen icons
for SIZE in 16 32 64 128 256 512; do
sips -z $SIZE $SIZE $INPUT --out $ICONSET/icon_${SIZE}x${SIZE}.png ;
done

# Retina display icons
for SIZE in 32 64 256 512; do
sips -z $SIZE $SIZE $INPUT --out $ICONSET/icon_$(expr $SIZE / 2)x$(expr $SIZE / 2)x2.png ;
done

# Make a multi-resolution Icon
iconutil -c icns -o $OUTPUT $ICONSET
rm -rf $WORKDIR
