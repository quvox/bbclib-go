bbclib-go
====
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)
[![Build Status](https://travis-ci.org/quvox/bbclib-go.svg?branch=develop)](https://travis-ci.org/quvox/bbclib-go)
[![Go Report Card](https://goreportcard.com/badge/github.com/quvox/bbclib-go)](https://goreportcard.com/report/github.com/quvox/bbclib-go)
[![Coverage Status](https://coveralls.io/repos/github/quvox/bbclib-go/badge.svg?branch=develop)](https://coveralls.io/github/quvox/bbclib-go?branch=develop)

Golang implementation of bbc1.core.bbclib and bbc1.core.libs modules in beyond-blockchain/bbc1

### Features
* Support serializing/deserializing BBc-1 transaction object
    * transaction version 1.2 or later
* Support sign/verify transaction
* Utility methods for creating transaction are not implemented
    * Need to set information to struct BBcTransaction and its members directly

### dependencies
* https://github.com/beyond-blockchain/libbbcsig


## Install

For linux/mac
```
sh prepare.sh
```

For Amazon Lambda, you need docker and do the following:
```
sh prepare.sh aws
```

After finishing prepare.sh script, you will find libbbcsig.dylib or libbbcsig.so in bbclib/.

