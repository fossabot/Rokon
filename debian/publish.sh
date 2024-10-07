#!/usr/bin/env bash

debuild -S -b -i -I
dput ppa:brycensranch/rokon-stable ~/rokon_*.changes
