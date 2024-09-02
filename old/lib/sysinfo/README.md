# Welcome to the Sysinfo library directory.

This was created because I couldn't find a good cross platform library that exposes system information in a simple way. I wanted to be able to get the number of CPUs, the amount of memory used, I also wanted to be able to get the number of CPUs and the amount of memory on the running system.

This library simply wraps the sysinfo rust crate and exposes the information in a simple way to golang with FFI and cgo.
Expected to dramatically slow down the build process, but it's a small price to pay for the convenience of having a simple cross platform library that exposes system information.
