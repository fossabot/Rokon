  <p align="center">
      <img src="resources/icon.png" width="240" />
      <h1 align="center" >Rokon (Roku Remote for your computer) </h1>
  </p>
  <p align="center">
    <b> Control your Roku from your Desktop or Laptop or whatever can output a display. Forget the batteries.</b>
  </p>

> NUCLEAR POWERED BY GOLANG!

> **Note:** This project is still in development and is not yet ready for general use. Please check back later for updates.

> This application was rewritten from Electron to Go for performance and stability reasons.

> This application is not affiliated with Roku, Inc. in any way.
> All product names, logos, and brands are property of their respective owners. All company, product, and service names used in this website are for identification purposes only.

> Also, none of the features listed below are implemented yet. This is a roadmap for the future.

This application provides a remote control interface for Roku devices, utilizing React for the frontend.

## Features

- Control your Roku device remotely with a sleek interface.
- Supports various functions such as navigation, volume control, input selection, typing, and so much more.
- **Speed**, go faster than any Roku remote could dream of, all without the noise.
- Purely use your keyboard to control your TV (Neovim Mode)
- Automatic Roku Discovery via [SSDP](https://www.pcmag.com/encyclopedia/term/ssdp) (You can manually input your Roku IP)
- Search your installed Roku apps and channels and quickly launch them. (Roadmap)
- DiscordRPC integration, display what you're doing on your Roku on Discord!
- [ActivityWatch](https://activitywatch.net/) integration (Roadmap)
- Option to run on startup and optionally turn on your Roku
- Use open source LLM to convert speech to text for voice control (Roadmap)
- Use your Xbox or PlayStation controller to control your Roku (Roadmap)
- Scripting functionatlity (Roadmap)
- CLI (Roadmap)
- Run actions such as auto scanning at a certain time (Roadmap)
- Webhook support (Roadmap)
- Installing channels (Roadmap)
- Launching things like YouTube with a video (Roadmap)

## Screenshots

Below is an example screenshot of the application:

![Example Screenshot](screenshots/example.png)

_(Screenshot taken on March 11, 2024)_

## Installation

To install the app, simply download the appropriate installer for your platform from the [releases page](https://github.com/BrycensRanch/Rokon/releases) and follow the installation instructions.

## Usage

Once installed, launch the application, and you'll be greeted with a remote control interface. Use the buttons to control your Roku device.

## Roku ECP API Integration

The application communicates with Roku devices using the Roku External Control Protocol (ECP) API. This allows for seamless control and interaction with Roku devices.

## Undocumented API Calls

Additionally, the app leverages some undocumented API calls to gain an edge over the competition, providing enhanced functionality and a better user experience.