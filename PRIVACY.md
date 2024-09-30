# Rokon's Privacy Policy

> Created on 9/11/2024

By default, Rokon collects telemetry data about how my application is performing in two ways:

- Application usage data (Tells us how many people are using my application on what systems)
- Error reporting and performance data. (Automatically raises an internal issue when you experience a crash with telemetry on)
- Your Roku's local IP ie `10.0.0.32` will be collected to better understand what local IP addresses Rokon should be looking out for.
- What Operating System you're running and which version you're using. ie `Windows 10 Home` or `macOS Sonana` or `Fedora Linux 42 (KDE Plasma) x86_64 Linux 6.10.9`
- Application version ie `rokon_1.0.0+44e1612`
- Whether or not the application is being ran through a translation layer like [WINE](https://winehq.org) or [Rosetta](https://en.wikipedia.org/wiki/Rosetta_(software)).
- What CPU your computer has ie `i7-11800H`
- What GPU your computer has ie `RTX 3060 Laptop GPU`
- How much memory the application is using (Helps to identify memory leaks)
- How much overall memory your computer has to better understand the "headroom space" left.
- Application environment ie running it on Linux with AppImage, Flatpak, Snap. Or running the portable version of my application on Windows.
- Your Roku's model number ie `C177X`
- Your Roku's operating system ie `Roku/13.1.4 UPnP/1.0 Roku/13.1.4`
- General Region ie `United States`
- How many times you've ran the application ie `19`, minutesActive ie `320`
- [SSDP](https://en.wikipedia.org/wiki/Simple_Service_Discovery_Protocol) broadcast address.

All this data helps to improve Rokon as it is [Free software](https://www.gnu.org/licenses/agpl-3.0.en.html).

### Definitions

telemetry - Modern, dynamic distributed systems require comprehensive monitoring to understand software behavior in various situations. Developers face challenges tracking the software’s performance in the field and responding to various modifications. To keep up with continuously changing requirements, it’s essential to have a simple way to collect data from systems the application is running.

SSDP - On devices and PCs that support SSDP, this feature can be enabled, disabled, or paused. When SSDP is enabled, devices communicate information about themselves and the services they provide to any other UPnP client. Using SSDP, computers connected to the network also provide information about available services.

anonymous - not identified by name; of unknown name.

## What I will not do

- Sell your data
- Spy on what you're doing on your TV
- Collect non-anonymous data such as your name, your TV's friendly name (Brycen's Living Room Roku), etc.
- Be evil

For error tracking and performance data, we use Sentry.io and their privacy policy can be found [here](https://sentry.io/privacy).

For application analytics, we use Aptabase.com and their privacy policy can be found [here](https://aptabase.com/legal/privacy).

## In terms of fingerprinting

- My application turns your machineid into a sha256 value that isn't traceable to you yet essentially unique to your operating system's installation.
- This sha'd machineid helps me to understand user flow and how they could've ran into issues and helps to discern one device from another.
- Although your *public* IP address is naturally processed by these network services. It is not saved and thus discarded.

## How do I disable telemetry?

Go into the application's settings menu or config file (~/.config/rokon on Linux)

Set `telemetry` to `false` and then no data should be sent for telemetry or application analytics.

## How do I request for my personal data to be removed?

There is no personal data. All data collected is anonymous.

## Final notes

I made Rokon because I wanted to learn more about coding. I decided to make it free and open source for you, and this is really all I expect to learn. I doubt I'll even get any donations for my work. This is all I ask, **keep telemetry on**. Help me ***improve*** Rokon.
