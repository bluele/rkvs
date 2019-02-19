# rkvs

## Build

```
$ dep ensure
$ make build
```

## Getting Started

First create a new cluster.
```
$ ./build/rkvs start --id=first --address=127.0.0.1:10000 --bootstrap=true
# In another terminal
$ ./build/rkvs start --id=second --address=127.0.0.1:10001 --join=first@127.0.0.1:10000
```

After creating the cluster, you can execute client commands.
```
$ ./build/rkvs servers
suffrage=0 id=first address=127.0.0.1:10000
suffrage=0 id=second address=127.0.0.1:10001
$ ./build/rkvs kvs write key1 value1
$ ./build/rkvs kvs read key1
value1
```

## Author

**Jun Kimura**

* <http://github.com/bluele>
* <junkxdev@gmail.com>
