![GLC (glog cleaner)](./glc.png "GLC (glog cleaner)")

# GLC (glog cleaner)

[![Build Status](https://travis-ci.org/xuri/glc.svg?branch=master)](https://travis-ci.org/xuri/glc)
[![Code Coverage](https://codecov.io/gh/xuri/glc/branch/master/graph/badge.svg)](https://codecov.io/gh/xuri/glc)
[![Go Report Card](https://goreportcard.com/badge/github.com/xuri/glc)](https://goreportcard.com/report/github.com/xuri/glc)
[![GoDoc](https://godoc.org/github.com/xuri/glc?status.svg)](https://godoc.org/github.com/xuri/glc)
[![license](https://img.shields.io/github/license/mashape/apistatus.svg?maxAge=2592000)](https://github.com/xuri/glc/blob/master/LICENSE)
[![Donate](https://img.shields.io/badge/Donate-PayPal-green.svg)](https://www.paypal.me/xuri)

## Overview

GLC (glog cleaner) is a log clear for glog written in Go. This library support for deleting old logs. There are tools which can be run to do the cleanup such as logrotate, but logrotate can't runs on Windows and embedded system, so we need a cross platform library to rotate the log.

## Installation

```go
go get github.com/xuri/glc
```

## Usage

Here is a example usage that will do check every hour and clean the log files that creation time older than 30 minutes.

```go
glc.NewGLC(glc.InitOption{
	Path:     path,
	Prefix:   `glc`,
	Interval: time.Duration(time.Hours),
	Reserve:  time.Duration(time.Minute * 30),
})
```

## Contributing

Contributions are welcome! Open a pull request to fix a bug, or open an issue to discuss a new feature or change.

## Licenses

This program is under the terms of the MIT License. See [LICENSE](https://github.com/xuri/glc/blob/master/LICENSE) for the full license text.