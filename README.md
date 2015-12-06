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
##### GPG
Used to encrypt all data

* OS X: `brew install gpg` 
* Linux: `sudo apt-get install gnupg`

##### xclip (Linux only)
Used to copy passwords to the clipboard on Linux systems

* Linux: `sudo apt-get install xclip`

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
* Don't want to type in the password everytime?
    ```sh
    export PICK_TOKEN=<PASSWORD HERE>
    ```
