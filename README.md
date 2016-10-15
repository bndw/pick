pick
====
[![Build Status](https://travis-ci.org/bndw/pick.svg?branch=master)](https://travis-ci.org/bndw/pick)

A minimal password manager for OS X and Linux.

![demo](https://github.com/bndw/pick/raw/master/demo.gif)

## Install

#### go get
```sh
$ go get github.com/bndw/pick
```

#### Homebrew
```sh
$ brew tap bndw/pick
$ brew install bndw/pick/pick-pass
```

#### The old fashioned way
```sh
$ git clone https://github.com/bndw/pick && cd pick
$ make
$ make install
```

## Usage
```
Usage:
  pick [command]

Available Commands:
  add         Add a credential
  backup      Backup the safe
  cat         Cat a credential
  cp          Copy a credential to the clipboard
  edit        Edit a credential
  export      Export decrypted credentials in JSON format
  ls          List all credentials
  rm          Remove a credential
  version     Print the version number of pick

Use "pick [command] --help" for more information about a command.
```

## Similar software
* [pwd.sh: Unix shell, GPG-based password manager](https://github.com/drduh/pwd.sh)
* [Pass: the standard unix password manager](http://www.passwordstore.org/)
