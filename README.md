bbclib-go
====
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

