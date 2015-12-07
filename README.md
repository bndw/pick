pick
====
A minimal password manager for OS X and Linux.

![demo](https://github.com/bndw/pick/raw/master/demo.gif)

Features
--------
* GPG for encryption
* JSON formatted data
* Password generation
* Environment Variable configuration

Dependencies
------------
#### GPG
* **OS X**: `brew install gpg` 
* **Linux**: `sudo apt-get install gnupg`

#### xclip (Linux only)
* **Linux**: `sudo apt-get install xclip`

Installation
------------
1. Clone the repository
    ```sh
    git clone https://github.com/bndw/pick.git && cd pick
    ```

2. Copy the `pick` executable into your PATH
    ```sh
    cp pick /usr/local/bin
    ```

Commands
--------
```
add [ALIAS] [USERNAME] [PASSWORD]     Add a credential to the safe
cat ALIAS                             Print a credential to STDOUT
cp  ALIAS                             Copy a credential's password to the clipboard
ls                                    List credentials
rm  ALIAS                             Remove a credential
```

Config
------

* Override the safe location (default: ~/.pick.safe)
    ```sh
    export PICK_SAFE=/path/to/pick.safe
    ```
