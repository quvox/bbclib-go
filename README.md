bbclib-go
====
Golang implementation of bbclib.py in beyond-blockchain/bbc1

### Features
* Support serializing/deserializing BBc-1 transaction object
    * transaction version 1 or later (original one in bbc1 v.1.0 was 0)
    * only BSON format is supported (ZLIB compression is also available)
* Support sign/verify transaction
* Utility methods for creating transaction are not implemented
    * Need to set information to struct BBcTransaction and its members directly

### dependencies
* libbbcsig.so used in bbc1 is needed


## Install

For linux/mac
```
sh prepare.sh
```

For Amazon Lambda, you need docker and do the following:
```
sh prepare.sh aws
```

After finishing prepare.sh script, you will find libbbcsig.so in bbclib/.

