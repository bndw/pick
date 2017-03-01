pick
====
[![Build Status](https://travis-ci.org/bndw/pick.svg?branch=master)](https://travis-ci.org/bndw/pick)

A minimal password manager for OS X and Linux.

![demo](./demo.gif)

## Install

#### go get
```sh
$ go get -u github.com/bndw/pick
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

## Getting started

If you haven't used `pick` before, you first need to initialize your pick safe. This is straightforward:
```sh
$ pick init
```
Running `pick init` will ask you for a master password. Your master password is used to encrypt your pick safe. As this is the only password you need to remember to access all passwords and notes stored in your pick safe, make this a strong and unique one! Use `pick pass` to generate a strong password if you think you're not creative enough :).

### Adding a credential

Once `pick` has been initialized, adding a new credential is easy:
```sh
$ pick add github
```
This will ask you for your master password first which is required to store something in the pick safe.
Then type in your username which should be used for the `github` credential.
`pick` will now ask you if you already have a password for `github` or if should create a new one for you.
Done. Credential added.

### Listing your credentials

```sh
$ pick ls
```

### Copy a credential's password to your clipboard

```sh
$ pick cp github
```

For all commands, please refer to [Usage](#usage).

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
  mv          Rename a credential
  note        Create a note
  pass        Generate a password without storing it
  rm          Remove a credential
  sync        Sync current safe with another pick safe
  version     Print the version number of pick

Use "pick [command] --help" for more information about a command.
```

## Threat model

Although `pick` is focussed on security, once an adversary somehow gains write access to your computer where the `pick` binary is stored, he can simply exchange it and record your master password to decrypt your pick safe.
If you decide to store the pick safe on a remote drive (e.g. a remote server), the system will be secure even if an adversary can modify the pick safe. As this safe is encrypted and authenticated, he can not modify or decrypt it.

## Similar software
* [pwd.sh: Unix shell, GPG-based password manager](https://github.com/drduh/pwd.sh)
* [Pass: the standard unix password manager](https://www.passwordstore.org/)
