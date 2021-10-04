Terraform plusserver Provider
==================

- Website: https://www.terraform.io

<img src="https://cdn.rawgit.com/hashicorp/terraform-website/master/content/source/assets/images/logo-hashicorp.svg" width="600px">

Requirements
------------

-	[Terraform](https://www.terraform.io/downloads.html) 0.12.x
-	[Go](https://golang.org/doc/install) 1.17 (to build the provider plugin)

Building The Provider
---------------------

Clone repository to

```sh
$ git clone git@github.com:plusserver/terraform-provider-plusserver.git
```

building the provider

```sh
$ make build
```

install the provider for local development
```sh
$ make install
```

Using the provider
----------------------

Please see the documentation at [terraform.io](#).

Or you can browse the documentation within this repo [here](#).

Developing the Provider
---------------------------

Testing the Provider
--------------------
