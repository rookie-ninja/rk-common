<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->
**Table of Contents**  *generated with [DocToc](https://github.com/thlorenz/doctoc)*

- [rk-common](#rk-common)
  - [Installation](#installation)
  - [Quick Start](#quick-start)
    - [Context](#context)
  - [Contributing](#contributing)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

# rk-common
The common library mainly used by rk-boot

## Installation
`go get -u rookie-ninja/rk-common`

## Quick Start

### Context
A struct called AppContext witch contains RK style application metadata.

| Element | Description | Default |
| ------ | ------ | ------ |
| application | name of running application | empty |
| startTime | application start time | 0001-01-01 00:00:00 +0000 UTC |
| loggers | loggers with a name as a key | empty map |
| viperConfigs | viper configs with a name as a key | empty map |
| rkConfigs | rk style configs with a name as a key | empty map |
| logger | default logger whose name is "default" | zap.NewNoop() |
| eventFactory | event data factory | standard event factory which logs to stdout |
| shutdownSig | a channel receiving shutdown signals | empty channel |
| customValues | custom k/v store | empty map |

## Contributing
We encourage and support an active, healthy community of contributors &mdash;
including you! Details are in the [contribution guide](CONTRIBUTING.md) and
the [code of conduct](CODE_OF_CONDUCT.md). The rk maintainers keep an eye on
issues and pull requests, but you can also report any negative conduct to
dongxuny@gmail.com. That email list is a private, safe space; even the zap
maintainers don't have access, so don't hesitate to hold us to a high
standard.

<hr>

Released under the [MIT License](LICENSE).