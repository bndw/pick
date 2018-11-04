pick
====
[![Build Status](https://travis-ci.org/bndw/pick.svg?branch=master)](https://travis-ci.org/bndw/pick)
[![Go Report Card](https://goreportcard.com/badge/github.com/bndw/pick)](https://goreportcard.com/report/github.com/bndw/pick)

A secure and easy-to-use password manager for macOS and Linux.

![demo](https://user-images.githubusercontent.com/4248167/29298817-564f4f54-811f-11e7-9a54-934afa1374df.gif)

## Features

* Strong, modern encryption with sensible defaults (ChaCha20-Poly1305, AES-GCM, OpenPGP)
* Configurable safe storage (file, AWS S3)
* Secure notes
* Built-in password generator
* Clipboard support
* Automatic backups
* Sync multiple safes
* Export accounts to JSON
* No external dependencies

## Install

#### go get

```sh
go get -u github.com/bndw/pick
```

#### Homebrew

```sh
brew install bndw/pick/pick-pass --build-from-source
```


#### From source

```sh
git clone https://github.com/bndw/pick && cd pick
make
make install
```

## Getting started

If you haven't used `pick` before, first initialize your safe to set a master
password:
```sh
pick init
```
Make your master password strong, unique, and don't forget it! You'll need your
master password to access your safe. Without it your safe can not be unlocked.

#### Add an account

```sh
pick add work/email
```

#### List accounts
 
```sh
pick ls
```

#### View an account

```sh
pick cat work/email
```

#### Copy a password to the clipboard

```sh
pick cp work/email
```

*For all commands, please refer to the [Usage](#usage) section with `pick --help`.*

## Usage

```
Usage:
  pick [command]

Available Commands:
  add             Add a credential
  cat             Cat a credential
  cp              Copy a credential to the clipboard
  edit            Edit a credential
  help            Help about any command
  init            Initialize pick
  ls              List all credentials
  mv              Rename a credential
  note            Create a note
  pass            Generate a password without storing it
  rm              Remove a credential
  safe            Perform operations on safe
  version         Print the version number of pick

Use "pick [command] --help" for more information about a command.
```

## Security

`pick` is focused on security and to this end it is _essential_ to only run the 
`pick` binary on a trusted computer. Conversely, you don't necessarily need to 
trust the computer or server storing the pick safe (e.g. `Amazon S3`). This is
because the pick safe is encrypted and authenticated and cannot by decrypted or
unnoticeably modified without the master password.

If you've found a vulnerability or a potential vulnerability in pick please
email us at pick-security@bndw.co. We'll send a confirmation email to
acknowledge your report, and we'll send an additional email when we've
identified the issue positively or negatively.

## Similar software
* [pwd.sh: Unix shell, GPG-based password manager](https://github.com/drduh/pwd.sh)
* [Pass: the standard unix password manager](https://www.passwordstore.org/)
