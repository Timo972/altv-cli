# altv cli

Incredibly flexible and easy to use altv server manager. Install or update only necessary files, reducing the bandwidth usage and time spent to a minimum.<br />
Supports every official module and continues working even on module file renamings by respecting their manifest.json.<br />
Extendable to support custom modules, feel free to open a pull request to add compatibility.

### Installation

```bash
go install github.com/timo972/altv-cli/cmd/altv@latest
```

### Usage

Type `altv --help` to get a list of all available commands and flags.<br />
Consider using Makefile or a package.json script to simplify the usage.<br />

Example `Makefile`:

```makefile
.PHONY: install update verify

dir = ./server
branch = dev
timeout = 30

install:
    altv install -p $dir -b $branch -m server -m data-files -m csharp-module -m js-module -t $timeout
update:
    altv update -p $dir -b $branch -m server -m data-files -m csharp-module -m js-module -t $timeout
verify:
    altv verify -p $dir -b $branch -m server -m data-files -m csharp-module -m js-module -t $timeout
```

Example `package.json`:

```json
{
  "scripts": {
    "altv-install": "altv install -p ./server -b dev -m server -m data-files -m csharp-module -m js-module -t 30",
    "altv-update": "altv update -p ./server -b dev -m server -m data-files -m csharp-module -m js-module -t 30",
    "altv-verify": "altv verify -p ./server -b dev -m server -m data-files -m csharp-module -m js-module -t 30"
  }
}
```

This way you can use `make install`, `make update` or `make verify` to install, update or verify your server files.<br />
If you prefer using npm, you can use `npm run altv-install`, `npm run altv-update` or `npm run altv-verify` instead.<br />

### Motivation

There are several altv server updater libraries and scripts out there, but none of them is as flexible and resilient as this one. I was annoyed by the fact that I always had to download the whole server files every time an update released. This is especially annoying if your internet connection poor. Additionally I wanted to have a tool that is able to update the server files even if the module files are renamed. See for example the js-module node library name.

### Features

- [x] Fast
- [x] Flexible
- [x] Resilient
- [x] Supports every official module
- [x] Supports custom modules
- [x] Reduces bandwidth usage to a minimum

### Planned Features
- [ ] Workspace configs for use of cli without having to set flags every time: altv init -p ./server -b dev -t 30
- [ ] JSON configuration for cdn's
- [ ] CI integrations