use std::ffi::CStr;
use std::ffi::CString;
extern crate sysinfo;

use std::io::{self, Write};
use std::ptr;
// use std::str::FromStr;
use sysinfo::{Components, Disks, Networks, Pid, Signal, System, Users};

#[no_mangle]
pub extern "C" fn hello(name: *const libc::c_char) {
    // Convert the input C string to a Rust &str safely
    let name_cstr = unsafe { CStr::from_ptr(name) };
    let name = match name_cstr.to_str() {
        Ok(s) => s,
        Err(e) => {
            eprintln!("Invalid UTF-8 string: {}", e);
            return;
        }
    };

    println!("Hello {}!", name);
}

#[no_mangle]
pub extern "C" fn whisper(message: *const libc::c_char) {
    let message_cstr = unsafe { CStr::from_ptr(message) };
    let message = message_cstr.to_str().unwrap();
    println!("({})", message);
}

#[no_mangle]
pub extern "C" fn sysinfo() -> *const libc::c_char {
    // Collect system information with error handling
    let system_name = sysinfo::System::name().unwrap_or_else(|| "<unknown>".to_owned());
    let kernel_version =
        sysinfo::System::kernel_version().unwrap_or_else(|| "<unknown>".to_owned());
    let os_version = sysinfo::System::os_version().unwrap_or_else(|| "<unknown>".to_owned());
    let long_os_version =
        sysinfo::System::long_os_version().unwrap_or_else(|| "<unknown>".to_owned());

    let info = format!(
        "System name:              {}\n\
         System kernel version:    {}\n\
         System OS version:        {}\n\
         System OS (long) version: {}",
        system_name, kernel_version, os_version, long_os_version
    );

    // Convert Rust String to CString and return a pointer to it
    match CString::new(info) {
        Ok(c_string) => c_string.into_raw(),
        Err(_) => ptr::null(),
    }
}

// This is present so it's easy to test that the code works natively in Rust via `cargo test`
#[cfg(test)]
pub mod test {

    use super::*;
    use std::ffi::CString;

    // This is meant to do the same stuff as the main function in the .go files
    #[test]
    fn simulated_main_function() {
        hello(CString::new("world").unwrap().into_raw());
        whisper(CString::new("this is code from Rust").unwrap().into_raw());
        println!("{}", sysinfo());
    }
}
