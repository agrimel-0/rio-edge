<p align="center">
    <img src="https://agrimel.org/wp-content/uploads/2022/07/Agrimel-Inc.-Main-Logo-2000x1500_cropped_white-337x68.png" width="337" height="68"/>
</p>   

<p align="center">
    <a href="https://go.dev/" target="_blank">
    <img src="https://img.shields.io/badge/go-%2300ADD8.svg?style=for-the-badge&logo=go&logoColor=white" alt="golang">
    <a href="./LICENSE" target="_blank">
    <img src="https://img.shields.io/github/license/SirFed/rio-server?style=for-the-badge" alt="license">
    <a href="https://discord.gg/E5XaRkMamn" target="_blank">
    <img src="https://img.shields.io/badge/Agrimel-5865F2?style=for-the-badge&logo=discord&logoColor=white" alt="discord">
    <a href="https://donorbox.org/agrimel-inc" target="_blank">
    <img src="https://img.shields.io/badge/sponsor-30363D?style=for-the-badge&logo=GitHub-Sponsors&logoColor=#EA4AAA" alt="donate"></a>
</p>

# rio-server
RIO stands for Remote-IO. This rio-grpc repository shows a fully fledged edge server in accordance to the AGR-101 standard.

This is currently in alpha stage!

## Dependencies

Since this is written in go, dependencies are resolved automatically by the `go.mod` file!

## Installation

Installation really depends on the hardware you are building for. You will have to set the correct environment flags for your go build script, this will include setting the correct `GOOS`, `GOARCH`, and potentially `GOARM`. There are several resources available online, but [this list](https://gist.github.com/asukakenji/f15ba7e588ac42795f421b48b8aede63) by [asukakenji](https://gist.github.com/asukakenji) is pretty awesome!

For a Raspberry Pi 4 running a 32 bit install of Arch Linux Arm this was proven to work:
```
env GOOS=linux GOARCH=arm GOARM=7 go build -o bin/server-arm64 main.go
```

## Configuration

Configuration is handled by a `.yaml` file. In this repo you can find a sample configuration file in the directory [config](./config/). There is also a configuration for the Raspberry Pi 4, and as project progresses more and more prebuilt configuration files will be added!

The sample configuration file `.server.config.yaml` is shown below:

```yaml
# Server related config
server:
    port: 51232
    alias: testserver

# Hardware IO related config
pins:
  - io:
      gpiochip: gpiochip0
      lineOffset: 5
      alias: GPIO5
      value: 0
      output: false
  - io:
      gpiochip: gpiochip0
      lineOffset: 11
      alias: pin1
      value: 1
      output: true

```

The `server:` paragraph defines the various settings for the server.
 - **port**: is the network port over which the gRPC service will be served on.
 - **alias**: is a friendly name you want your gRPC service to go by.

The `pins:` paragraphs is composed of several `io:` entries. Each of these are hardware GPIO definitions which will be exported when the service starts. Filling out this information can be tricky, especially if you unfamiliar with Linux. Checkout [libgpiod](https://github.com/brgl/libgpiod), an awesome command line utility for discovering and exercising GPIO!
 - **gpiochip**: is the gpiochip that manages the GPIO line you are interested in
 - **lineOffset**: is the line offset of the IO pin you are interested in as reported by *gpioinfo*. 
 - **alias**: this is a friendly alias.
 - **value**: default value at which you want it to start as. Usually 0 is low and 1 is high, but check your hardware first.
 - **output**: wether the pin is going to be an output or an input, **value** will be ignored if this is set to *false*.

## Usage
```
Usage:
  rio-server [command]
join us
  serve       Start serving the remote-io server

Flags:
  -c, --config string      config file (default "./config/.server.config.yaml")
  -h, --help               help for rio-server
  -t, --toggle             Help message for toggle
  -v, --verbosity string   Log level (debug, info, warn, error, fatal, panic (default "warning")

Use "rio-server [command] --help" for more information about a command.
```

Serve command:
```
  rio-server serve [flags]

Flags:
  -h, --help   help for serve

Global Flags:
  -c, --config string      config file (default "./config/.server.config.yaml")
  -v, --verbosity string   Log level (debug, info, warn, error, fatal, panic (default "warning")
```

## Known issues

Checkout the [issues tab](https://github.com/agrimel-0/rio-server/issues) for any known issues.

## Getting help & Getting involved

We have a Discord server! Feel free to [join us](https://discord.gg/E5XaRkMamn)!

If you want to contribute, feel free to start an issue, create a fork with a branch following the issue, and submit a merge request when ready! 

----

## Open source licensing info
[LICENSE](LICENSE)

<p align="center">
    <img src="https://agrimel.org/wp-content/uploads/2022/07/Agrimel-Inc.-Main-Logo-2000x1500_cropped_white-337x68.png" width="337/2" height="68/2"/>
</p>   

