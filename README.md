[![Join the chat at https://gitter.im/cloudfoundry-samples/go_service_broker](https://badges.gitter.im/Join%20Chat.svg)](https://gitter.im/cloudfoundry-samples/go_service_broker?utm_source=badge&utm_medium=badge&utm_campaign=pr-badge&utm_content=badge)

go_service_broker
=================

This is a service broker written in Go Language for Cloud Foundry. This service broker supports creating VMs on AWS or SoftLayer. Since VMs take some to spin up creating them is done asynchronously. 

This broker also supports creating service keys. This is a new feature added to the Service Broker v2.6 APIs. In this broker we implement service keys by creating SSH keys and adding them to the VM.

Finally, this broker also supports arbitrary service parameters. For AWS this is done by allowing the user to pass the `ami-id` to use when spinning up the VM. And for SoftLayer you can specify the Virtual Guest Device Block Device Template Group ID.

NOTE that for AWS, in this implementation, service bind operation will generate a keypair and inject the public key into that EC2 instance and return the corresponding credentials including private key, user name and public IP address information which can be used to ssh login that EC2 instance. The service unbind operation will revoke that public key from the EC2 instance.

Presentations
=============

This sample project has been presented in CF Summit 2015:

* [Presentation on CloudFoundry Summit 2015](https://www.youtube.com/watch?v=MrSy4iZZPZE)

Videos for each of the main features of the brokers are in these Youtube videos:

* [Service Async Demo video](https://www.youtube.com/watch?v=Ij5KSKrAq9Q)
* [Service Key Demo video](https://www.youtube.com/watch?v=V5uzLcPQPmo)
* [Service Arbitrary Parameters Demo video](https://www.youtube.com/watch?v=Qc3bZljGscs)

The following blog post on IBM's OpenTech web site covers the broker in much details:

* [CloudFoundry Services Keys and Sample Go Service Broker](https://developer.ibm.com/opentech/2015/07/09/cloudfoundry-services-keys-and-sample-go-service-broker/)

Getting Started
===============

Get Latest Executable: go_service_broker
----------------------------------------

Assuming you have a valid [Golang 1.4.2](https://golang.org/dl/) or [later](https://golang.org/dl/) installed for your system, you can quickly build and get the latest `go_service_broker` executable by running the following `go` command:

```
$ go get github.com/cloudfoundry-samples/go_service_broker
```

This will build and place the `go_service_broker` executable built for your operating system in your `$GOPATH/bin` directory.


Building From Source
--------------------

Clone this repo and build it. Using the following commands on a Linux or Mac OS X system:

```
$ mkdir -p go_service_broker/src/github.com/cloudfoundry-samples
$ export GOPATH=$(pwd)/go_service_broker:$GOPATH
$ cd go_service_broker/src/github.com/cloudfoundry-samples
$ git clone https://github.com/cloudfoundry-samples/go_service_broker.git
$ cd go_service_broker
$ ./bin/build
```

NOTE: if you get any dependency errors, then use `go get path/to/dependency` to get it, e.g., `go get github.com/onsi/ginkgo` and `go get github.com/onsi/gomega`

The executable output should now be located in: `out/go_service_broker`. Place it wherever you want, e.g., `/usr/local/bin` on Linux or Mac OS X.

Dependencies
------------

Install `godep`.

```
$ go get github.com/tools/godep
```

Download and install packages with dependencies by using godep.

```
$ cd -
$ godep get ./...
```

Save the dependencies by godep.

```
$ godep save ./...
```

Build your executable `out/go_service_broker`.

```
$ bin/build
```

Configuring for AWS
-------------------

Before running the service broker, you need to configure your AWS accpunt's credentials. If you do not have AWS account, then you can get one for [free here](https://aws.amazon.com/free).

As a best practice, we recommend creating an IAM user that has access keys rather than relying on root access keys. You can login into your AWS account to create a new user 'service_broker' with the option to generate an access key for this user. 

Once you get a Access Key ID and Secret Access Key, copy and save it into `~/.aws/credentials` file, which might look like:

```
[default]
aws_access_key_id = YOUR-AWS-ACCESS-KEY-ID
aws_secret_access_key = YOUR-AWS-SECRET-ACCESS-KEY
```

Configuring for SoftLayer
-------------------------

For SoftLayer the configuration requires you to supply your SL user name and API key. If you do not have an SL account, please get one for [free here](http://www.softlayer.com/promo/freeCloud/freecloud). The API key can be requested once you login to your account.

You need to setup two environment variables with your SL credentials as follows.

```
export SL_USERNAME=your-softlayer-username@your-company.com
export SL_API_KEY=YOUR-SOFTLAYER-API-KEY
```

These two environment variables must exist where you run your broker. Locally, in a VM or server process, or whithin CloudFoundry. See below on details on how to run broker in CF or locally.

Running Broker
==============

The broker can be ran in one of two modes: locally or as an app in a CF environment.

Locally
-------

Run the executable to start the service broker which will listening on port `8001` by default.

```
$ out/go_service_broker --cloud AWS
```

This will run the broker in `AWS` mode. You can also specify `SoftLayer` mode with:

```
$ out/go_service_broker --cloud SoftLayer
```

If no argument is passed to the `--cloud` flag then AWS mode is assumed/


In CF
-----

When running the broker in a CF environment (including [BOSH lite](https://github.com/cloudfoundry/bosh-lite)). You simply need to:

```
$ git clone https://github.com/cloudfoundry-samples/go_service_broker.git
$ cd go_service_broker
$ cf push
```

You, of course, need to have the [CF CLI](https://github.com/cloudfoundry/cli) installed into your system. Also, you can edit the `Procfile` if you want to specify a different mode (AWS or SoftLayer) as well as any additional optional parameters to the CF Golang buildpacks.

Using Broker
============

TODO

License
=======
This is under [Apache 2.0 OSS license](https://github.com/cloudfoundry-samples/go_service_broker/LICENSE).
