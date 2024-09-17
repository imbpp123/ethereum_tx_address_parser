# Ethereum transaction parser

This parser watchs Ethereum network for specific address and returns it's transaction information.

## How to run

Run project (address is hardcoded in main.go):

```
make run
```

Run tests:
```
make tests
```

## Solution Architecture

Solution is built using DDD approach. There are following layers in application:

* ethereum/data - contains data objects
* ethereum/rpc - clients for ethereum network communication
* ethereum/storage - storages for data objects, repositories
* ethereum/domain - services for domains: block, transaction, address

Everything is combined in Parser. 
 
### Tests

Each package has unit test coverage
Parser has 1 integration test

### Logging

Domain package write logs because it is business logic of application. 
Log level can be changed in main.go
