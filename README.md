go_service_broker
=================

This is a service broker written in Go Language for Cloud Foundry.

Getting Started
===============
* Clone the git repository and setup go environment, make sure `$GOPATH` is correctly setup.
* Install `godeps`: 

```
$ go get github.com/tools/godep
```

* Save the dependencies:

```
$ godep save
```

* Build your executable `out/broker`

```
bin/build
```

* Run the executable to start the service broker which will listening on port `8001` by default

```
out/broker
```


License
=======
This is under [Apache 2.0 OSS license](https://github.com/xingzhou/go_service_broker/LICENSE).
