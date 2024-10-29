# Debian package for Rokon

This directory generates a deb package for Rokon from CI.

The stable versions can be found in this PPA

https://launchpad.net/~brycensranch/+archive/ubuntu/rokon-stable

the nightly version will be created once Rokon is stable (it isn't)

to build a debian package do

```bash
debuild -us -uc -b -i -I
```
