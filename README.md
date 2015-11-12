pick
====
A tiny password manager for OS X and Linux.

![demo](https://github.com/bndw/pick/raw/master/demo.gif)

Features
--------
* JSON formatted data
* Environment Variable configuration
* GPG for encryption

Dependencies
------------
* GPG
   ```sh
   # OS X
   brew install gpg
   ```
   ```sh
   # Linux
   sudo apt-get install gnupg
   ```

* xclip (Linux only)
   ```sh
   # Linux
   sudo apt-get install xclip
   ```

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

3. Initialize pick
    ```sh
    pick --init
    ```

Usage
-----
* Write a password (interactive)
    ```sh
    pick -w
    ```

* Read a password by alias
    ```sh
    pick github
    ```

* Read all passwords
    ```sh
    pick all
    ```

Advanced
--------

### Environment Variables

* Don't want to type in the password everytime?
    ```sh
    export PICK_TOKEN=<PASSWORD HERE>
    ```

* Don't ever want your passwords printed to stdout?
    ```sh
    export PICK_CONFIG='{"silent":true}'
    ```

* Want to print additonal metadata along with passwords?
    ```sh
    export PICK_CONFIG='{"verbose":true}'
    ```
