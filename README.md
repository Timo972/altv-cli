# altv cli

Incredibly flexible and easy to use altv server manager. Install or update only necessary files, reducing the bandwidth usage and time spent to a minimum.<br />
Supports every official module and continues working even on module file renamings by respecting their manifest.json.<br />
Extendable to support custom modules, feel free to open a pull request to add compatibility.

### Installation
```bash
go install github.com/timo972/altv-cli/cmd/altv@latest
```

### Motivation
There are several altv server updater libraries and scripts out there, but none of them is as flexible and resilient as this one. I was annoyed by the fact that I always had to download the whole server files every time an update released. This is especially annoying if your internet connection poor. Additionally I wanted to have a tool that is able to update the server files even if the module files are renamed. See for example the js-module node library name.

### Features
- [x] Fast
- [x] Flexible
- [x] Resilient
- [x] Supports every official module
- [x] Supports custom modules
- [x] Reduces bandwidth usage to a minimum
