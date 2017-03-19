gocoreutils
===

Reimplementation of UNIX core utilities by Golang

## Description
Reimplementation of UNIX core utilities, such as tail, pwd..etc by Golang.
These commands are effortly implementing to according to POSIX 1003.1 specification,and does not support GNU/BSD original options or something special actions. 

Currently not all command has been supported yet. 
See Currently Supported Command section below.

This is my first 'decent' Go language project for practicing.
If you found any mistakes or bugs in command, please tell me.

## Installation

```
$ go install github/necomeshi/gocoreutil
```

## Usage
To use the command,

``` 
$ ln -sv coreutils COMMAND_NAME
$ ./COMMAND_NAME
```

Replace ```COMMAND_NAME``` as a actual command name.
For example, to use the command 'pwd',

```
$ ln -sv coreutils pwd
$ ./pwd
```
## Currently Supported Command
basename, dirname, head, md5sum, pwd, tail, wc

## FAQ
1. Why you reimplemented ?
 Because for my practicing golang.
 
1. Why xx command has not been implemented ? When will you implement it ?
 Sometime when I need it. Or sometime when others give me an early Xmas present.


## Author
C.Yoshimura