go_service_broker
=================

This is a service broker written in Go Language for Cloud Foundry.

* The service broker supports asynchronously creating EC2 instances with arbitrary parameters. The state of the instance (in progress/succeeded/failed) can be retrieved later.
* The service broker also support service key feature which means application is not required when binding a service instance. The credentials of the service instance will be created.
* In this implementation, service bind operation will generate a keypair and inject the public key into that EC2 instance and return the corresponding credentials including private key, user name and public IP address information which can be used to ssh login that EC2 instance. The service unbind operation will revoke that public key from the EC2 instance.

This sample project has been presented in CF Summit 2015. Please refer to the presentation and demo videos on youtube:

* [Presentation on CloudFoundry Summit 2015](https://www.youtube.com/watch?v=MrSy4iZZPZE)
* [Service Async Demo](https://www.youtube.com/watch?v=Ij5KSKrAq9Q)
* [Service Key Demo](https://www.youtube.com/watch?v=V5uzLcPQPmo)
* [Service Arbitrary Parameters Demo](https://www.youtube.com/watch?v=Qc3bZljGscs)

Getting Started
===============
* Clone the git repository and setup go environment, make sure `$GOPATH` is correctly setup.

* Install `godep`.

```
$ go get github.com/tools/godep
```

* Build godep.

```
$ cd ../../tools/godep/
$ go build
$ export PATH=$PATH:$GOPATH/bin
```

* Download and install packages with dependencies by using godep.

```
$ cd -
$ godep get ./...
```

* Save the dependencies by godep.

```
$ godep save ./...
```

* Build your executable `out/broker`.

```
$ bin/build
```

* Before running the service broker, you need to configure your AWS credentials. As a best practice, we recommend creating an IAM user that has access keys rather than relying on root access keys. You can login into your AWS account to create a new user 'service_broker' with the option to generate an access key for this user. Once you get a Access Key ID and Secret Access Key, copy and save it into ~/.aws/credentials file, which might look like:

```
[default]
aws_access_key_id = AKID1234567890
aws_secret_access_key = MY-SECRET-KEY
```

* Run the executable to start the service broker which will listening on port `8001` by default.

```
$ out/broker
```


License
=======
This is under [Apache 2.0 OSS license](https://github.com/cloudfoundry-samples/go_service_broker/LICENSE).
