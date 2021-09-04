# rk-common
The common library mainly used by rk-entry.

[![build](https://github.com/rookie-ninja/rk-common/actions/workflows/ci.yml/badge.svg)](https://github.com/rookie-ninja/rk-common/actions/workflows/ci.yml)
[![codecov](https://codecov.io/gh/rookie-ninja/rk-common/branch/master/graph/badge.svg?token=7YYDHG1JVL)](https://codecov.io/gh/rookie-ninja/rk-common)
[![Go Report Card](https://goreportcard.com/badge/github.com/rookie-ninja/rk-common)](https://goreportcard.com/report/github.com/rookie-ninja/rk-common)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)

<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->
**Table of Contents**  *generated with [DocToc](https://github.com/thlorenz/doctoc)*

- [Installation](#installation)
- [Quick Start](#quick-start)
  - [common](#common)
  - [flags](#flags)
    - [Usage of rkboot:](#usage-of-rkboot)
    - [Usage of rkset:](#usage-of-rkset)
  - [strvals](#strvals)
- [Contributing](#contributing)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

## Installation
`go get -u github.com/rookie-ninja/rk-common`

## Quick Start
### common
Utility functions mainly used with rk series of packages, including rk-gin, rk-grpc, rk-entry and etc.

Please don't add any dependency cycle with above packages.

### flags
pflag.FlagSet which contains **rkboot** and **rkset** as key.

#### Usage of rkboot:
Receives path of boot config file path, can be either absolute of relative path.
If relative path was provided, then current working directory would be attached in front of provided path.

example:
```bash
./your_compiled_binary --rkboot <your path to config file>`
```

#### Usage of rkset:
Receives flattened boot config file(YAML) keys and override them in provided boot config.

We follow the way of HELM does while overriding keys, refer to https://helm.sh/docs/intro/using_helm/
example:

Lets assuming we have boot config YAML file (example-boot.yaml) as bellow:
```yaml
gin:
  - port: 1949
    commonService:
      enabled: true
```

We can override values in example-boot.yaml file as bellow:
```bash
./your_compiled_binary --rkboot example-boot.yaml --rkset "gin[0].port=2008,gin[0].commonService.enabled=false"
```

Basic rules:
- Using comma(,) to separate different k/v section.
- Using [index] to access arrays in YAML file.
- Using equal sign(=) to distinguish key and value.
- Using dot(.) to access map in YAML file.

### strvals
Copied from https://github.com/helm/helm/blob/main/pkg/strvals/parser.go with some utility function.

## Contributing
We encourage and support an active, healthy community of contributors &mdash;
including you! Details are in the [contribution guide](CONTRIBUTING.md) and
the [code of conduct](CODE_OF_CONDUCT.md). The rk maintainers keep an eye on
issues and pull requests, but you can also report any negative conduct to
lark@rkdev.info.

Released under the [Apache 2.0 License](LICENSE).
