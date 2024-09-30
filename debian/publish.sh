#!/usr/bin/env bash

debuild -S -i -I
dput ppa:brycensranch/rokon-stable ~/rokon_*.changes
