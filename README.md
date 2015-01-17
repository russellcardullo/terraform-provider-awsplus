# terraform-provider-awsplus#

This project is a [terraform](http://www.terraform.io/) provider for Amazon services not builtin to terraform.

Currently this supports the following resources:

- [Amazon Kinesis](http://aws.amazon.com/kinesis/)
- [Amazon SQS](http://aws.amazon.com/sqs/)

## Build and install ##

### Dependencies ###

You should have a working Go environment setup.  If not check out the Go [getting started](http://golang.org/doc/install) guide.

This relies on the crowdmob fork of [goamz](https://github.com/crowdmob/goamz). To get that: `go get github.com/crowdmob/goamz`.

You'll also need the libraries from terraform.  Check out those docs under [plugin basics](http://www.terraform.io/docs/plugins/basics.html).

### Build ###

Run `go install github.com/russellcardullo/terraform-provider-awsplus`

### Install ###

Add the following to `$HOME/.terraformrc`

```
providers {
    awsplus = "$GOPATH/bin/terraform-provider-awsplus"
}
```

## Usage ##

```
variable "access_key" {}
variable "secret_key" {}
variable "region" {}

provider "awsplus" {
    access_key = "${var.access_key}"
    secret_key = "${var.secret_key}"
    region = "${var.region}"
}

resource "awsplus_kinesis_stream" "kinesis-example" {
    name = "terraform-kinesis"
    shard_count = 2
}

resource "awsplus_sqs_queue" "sqs-example" {
    name = "str-test-terraform-queue"
}
```

Apply with:
```
 terraform apply \
    -var 'access_key=AWS_ACCESS_KEY' \
    -var 'secret_key=AWS_SECRET_KEY' \
    -var 'region=AWS_REGION'
```

## Resources ##

### Kinesis Stream ###

This support create/read/delete methods on a Kinesis stream.  Update operations are not supported yet since those
resharding/merging Kinesis shards are slightly more complicated than other resources.

The following attributes can be set:

**name** - (Required) The name of the stream

**shard_count** - (Required) The number of shards to set in the stream.

All atribute changes will force a new resource.  Since the stream name is also used as Amazon's identifier
for a Kinesis stream you should first delete the existing stream or create one with a new name.

### SQS Queue ###

This support create/read/delete methods on an SQS queue.

The following attributes can be set:

**name** - (Required) The name of the queue
