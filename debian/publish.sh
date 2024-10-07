#!/usr/bin/env bash

debuild -S -i -I -d
debuild -b -i -I -d
dput ppa:brycensranch/rokon-stable ~/rokon_*.changes
