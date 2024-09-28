#!/bin/bash

# Attribution https://github.com/getseabird/seabird/blob/main/build/darwin/icns.sh


for size in 16 32 48 128 256; do
  inkscape -o $size.png -w $size -h $size ../assets/Rokon.svg
done
png2icns icon.icns 16.png 32.png 48.png 128.png 256.png
rm *.png
