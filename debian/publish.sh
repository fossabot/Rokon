#!/usr/bin/env bash

debuild -S -i -I -d
dput ppa:brycensranch/rokon-stable ~/rokon_*.changes
