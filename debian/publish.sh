#!/usr/bin/env bash

debuild -S -i -I
debuild -b -i -I
dput ppa:brycensranch/rokon-stable ~/rokon_*.changes
