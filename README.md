# Sysremot - System Resource Monitoring Tool

Sysremot monitors your Linux/Mac system resource. It gathers the resource metrics and use redis to save those metrics.

What it does:

- Monitor memory and CPU usage
- Monitor disk usage
- Monitor system uptime and load average
- Generate metrics and save it to redis

You can use those metrics to generate your own graph

## Requirements

- You need Linux or MacOS machine
- You need redis server

For now, I only test sysremot under Linux and MacOS system. Other systems maybe supported.

## Installation

You need to build this tool first.
Asuming that you already build the source code :

- Install redis, and configure
- Copy the tool binary to /usr/local/bin
- Refresh your shell, one method will be sourcing your .bashrc
- Run ```sysremot install```
- Modify the config file in /etc/sysremot/sysremot.env

## License
Copyright &copy; 2015 Donny Prasetyobudi.
Sysremot is released under the Apache 2.0 License. See LICENSE.