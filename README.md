<img alt="The riffraff logo" src="logo_withname.svg" height="192px" />

[![Travis (.org)](https://img.shields.io/travis/mre/riffraff/master.svg?style=flat-square)](https://travis-ci.org/mre/riffraff)
[![Codecov branch](https://img.shields.io/codecov/c/github/mre/riffraff/master.svg?style=flat-square)](https://codecov.io/gh/mre/riffraff)
[![GoDoc](https://img.shields.io/badge/godoc-reference-5272B4.svg?style=flat-square)](https://godoc.org/github.com/mre/riffraff)
[![Go Report Card](https://goreportcard.com/badge/github.com/mre/riffraff?style=flat-square)](https://goreportcard.com/report/github.com/mre/riffraff)

![usage](usage.gif)

A commandline interface for Jenkins.

## Features

* Queries the current status of jobs in parallel.
* Can trigger Jenkins builds from the commandline.
* Visualizes the status of jobs and nodes.
* Can diff the output two runs.

```Shell
riffraff is a commandline interface for Jenkins

Usage:
  riffraff [command]

Available Commands:
  build       Trigger build for all matching jobs
  diff        Print a diff between two builds of a job
  help        Help about any command
  log         Show the logs of a job
  nodes       Show the status of all Jenkins nodes
  open        Open a job in the browser
  queue       Show the queue of all matching jobs
  status      Show the status of all matching jobs

Flags:
  -h, --help      help for riffraff
      --salt      Show failed salt states
  -v, --verbose   Verbose mode. Print full job output

Use "riffraff [command] --help" for more information about a command.
```

## Installation

```Shell
go get github.com/mre/riffraff
```

...or download a static binary from the [releases page](https://github.com/mre/riffraff/releases).

## Getting started

You need to set the following environment variables:

```Shell
export JENKINS_URL="http://example.com/"
export JENKINS_USER="username"
export JENKINS_PW="password"
```

You might want to put those into your `~/.bashrc`, `~/.zshrc` or equivalent.

## Usage

```Shell
riffraff status jenkins-job-name
```

This will print the current status of all Jenkins jobs matching the given pattern (`jenkins-job-name` in this case).
You can use any regular expression for that, e.g.:

```Shell
riffraff status "^application-.*-unittests$"
```

You can get the full output of each last job matching the pattern with

```Shell
riffraff status -v "^application-.*-unittests$"
```

## Development

* Install golang version 1.11 or later for [go modules](https://github.com/golang/go/wiki/Modules) support
* Clone this repository to a directory in your `$GOPATH/src` tree (recommended) or use`go get -u github.com/mre/riffraff.git` (uses https not SSH)
* In the source folder run `go run main.go` to install modules and run `riffraff`
* If you don't have a Jenkins server you can run it by using its [`.war` file](https://jenkins.io/doc/pipeline/tour/getting-started/) (recommended) or installing its [`.deb` file](https://jenkins.io/doc/book/installing)

## OBTW

The tool is named after the [butler from the Rocky Horror Picture Show](https://en.wikipedia.org/wiki/The_Rocky_Horror_Picture_Show:_Let%27s_Do_the_Time_Warp_Again), and not the rapper with the same name ;-).

## Credits

Logo design: [Franziska Böhm noxmoon.de](http://noxmoon.de) ([CC-BY-SA](https://creativecommons.org/licenses/by-sa/4.0/))
