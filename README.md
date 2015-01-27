# Sysremot - System Resource Monitoring tool

Sysremot monitors your linux system. It gathers the resource metrics and use redis to save those metrics.

What it does:

- Monitor memory and CPU usage
- Monitor disk usage
- Monitor system uptime and load average
- Generate metrics and save it to redis

You can use those metrics to generate your own graph

## Requirements

For now, the target system is Linux and Darwin

## Installation

You need to build this tool first.
Asuming that you already have the builed tool :

- Install redis, and configure
- Copy the tool binary to /usr/local/bin
- Refresh your shell, one method will be sourcing your .bashrc
- Run ```sysremot install```
- Modify the config file in /etc/sysremot/sysremot.env

## Author

Sysremot is written by Netzumo Ninja.