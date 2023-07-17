# altv cli [![Test & Release](https://github.com/Timo972/altv-cli/actions/workflows/test-release.yml/badge.svg?branch=main)](https://github.com/Timo972/altv-cli/actions/workflows/test-release.yml)[![License][license-src]][license-href]

Incredibly flexible and easy to use altv server manager. Install or update only necessary files, reducing the bandwidth usage and time spent to a minimum.
Supports every official module and continues working even on module file renamings by respecting their manifest.json.
Extendable to support custom modules, feel free to open a pull request to add compatibility.

## Table of Contents
- [Motivation](#motivation)
    - [Features](#features)
    - [Planned Features](#planned-features)
- [Installation](#installation)
- [Usage](#usage)
    - [Example Makefile](#makefile)
    - [Example package.json](#packagejson)

## <a name="motivation"></a>Motivation

There are several altv server updater libraries and scripts out there, but none of them is as flexible and resilient as this one. I was annoyed by the fact that I always had to download the whole server files every time an update released. This is especially annoying if your internet connection poor. Additionally I wanted to have a tool that is able to update the server files even if the module files are renamed. See for example the js-module node library name.

### <a name="features"></a>Features

- ⚡ &nbsp;Fast
- 🔀 &nbsp;Flexible
- 💎 &nbsp;Resilient
- 🏅 &nbsp;Supports every official module
- 🛠 &nbsp;Supports custom modules
- 📉 &nbsp;Reduces bandwidth usage to a minimum

### <a name="planned-features"></a>Planned Features
- 🔨 &nbsp;Workspace configs for use of cli without having to set flags every time: `altv init -p ./server -b dev -t 30`
- ⚙ &nbsp;&nbsp;JSON configuration for cdn's
- 🤖 &nbsp;CI integrations

## <a name="installation"></a>Installation
If you have Go installed, you can install the cli with the following command:
```bash
go install github.com/timo972/altv-cli/cmd/altv@latest
```
Otherwise you can use the prebuilt binaries from github [releases]() or npm in the future.

## <a name="usage"></a>Usage

Type `altv --help` to get a list of all available commands and flags.<br />
Consider using Makefile or a package.json script to simplify the usage.<br />

Example `Makefile`:
<a name="makefile"></a>
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
<a name="packagejson"></a>
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

<!-- badges -->
[license-src]: https://img.shields.io/npm/l/@timo792/altv-cli.svg?style=flat&colorA=18181B&colorB=28CF8D
[license-href]: https://npmjs.com/package/@timo972/altv-cli

[npm-version-src]: https://img.shields.io/npm/v/@timo792/altv-cli/latest.svg?style=flat&colorA=18181B&colorB=28CF8D
[npm-version-href]: https://npmjs.com/package/@timo792/altv-cli