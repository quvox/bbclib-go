bbclib-go
====
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)
[![Build Status](https://travis-ci.org/quvox/bbclib-go.svg?branch=develop)](https://travis-ci.org/quvox/bbclib-go)
[![Go Report Card](https://goreportcard.com/badge/github.com/quvox/bbclib-go)](https://goreportcard.com/report/github.com/quvox/bbclib-go)
[![Coverage Status](https://coveralls.io/repos/github/quvox/bbclib-go/badge.svg?branch=develop)](https://coveralls.io/github/quvox/bbclib-go?branch=develop)
[![Maintainability](https://api.codeclimate.com/v1/badges/0c523f5a3d71b77aad46/maintainability)](https://codeclimate.com/github/quvox/bbclib-go/maintainability)

This repo is now obsoleted, and move to https://github.com/beyond-blockchain/bbclib-go.git

Golang implementation of bbc1.core.bbclib and bbc1.core.libs modules in https://github.com/beyond-blockchain/bbc1

### Features
* Support most of features of bbclib in https://github.com/beyond-blockchain/bbc1
    * BBc-1 version 1.2
    * transaction header version 1 only
* Go v1.10 or later

### dependencies
* https://github.com/beyond-blockchain/libbbcsig

## Usage

```bash
go get -u github.com/quvox/bbclib-go
```

Building an external library is also required.
When "go get" is done, you will find github.com/quvox/bbclib-go/ directory in ${GOPATH}/src.
Then, execute the following commands:
```
cd ${GOPATH}/src/github.com/quvox/bbclib-go
bash prepare.sh
```

If you want to use this module in an AWS environment (EC2 or Lambda), do as follows:
```
cd ${GOPATH}/github.com/quvox/bbclib-go
bash prepare.sh aws
```
The preparation script (prepare.sh) produces libbbcsig.a, which is a static link library for signing/verifying a transaction.
Building libbbcsig.a takes long time, so be patient.

After finishing the compilation, you are ready for "go install".

```
go install github.com/quvox/bbclib-go
```

NOTE: [example/](./example) directory includes a sample code for this module. There are a document and a preparation script. 

## Prepare for development (module itself)

For linux/mac
```
sh prepare.sh
```

For Amazon Lambda, you need docker and do the following:
```
sh prepare.sh aws
```

After finishing prepare.sh script, you will find libbbcsig.a and libbbcsig.h, which are used by keypair.go for signing/verifying.

