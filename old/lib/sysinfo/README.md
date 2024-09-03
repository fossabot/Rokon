# Welcome to the Sysinfo library directory.

> Important: This library must be compiled before you build Rokon's main go program. This library is written in Rust and must be compiled with the `cargo build --release` command. Then the resulting files must be moved to the `old/lib/sysinfo` directory. The files that need to be moved are `librokon_rust_sysinfo.a` and `librokon_rust_sysinfo.h`. The go program will not build without these files.

This was created because I couldn't find a good cross platform library that exposes system information in a simple way. I wanted to be able to get the number of CPUs, the amount of memory used, I also wanted to be able to get the number of CPUs and the amount of memory on the running system.

This library simply wraps the sysinfo rust crate and exposes the information in a simple way to golang with FFI and cgo.
Expected to dramatically slow down the build process, but it's a small price to pay for the convenience of having a simple cross platform library that exposes system information.
