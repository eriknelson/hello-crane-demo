# hello-crane-demo

This demo will walk through a basic use-case with Crane 2, a tool for rehosting
cloud workloads for Kubernetes. It seeks to highlight some of Crane's primary
goals:


* Compatible with vanilla k8s out-of-the-box

> NOTE: As a first pass attempt to get something demo-able, we're going to be
> using an openshift workload and tooling, and will aim to update this to remove
> the requirement.

* No requirement of elevated privileges; application owners can migrate their
own applications
* Extraction and discovery of all resources from source namespaces
* User readable / tweakable transforms applicable to the resources before
deployment to target cluster
* Onramp into gitops managed target deployments

## Getting started

Crane exists as two primary repos at the moment:

* https://github.com/konveyor/crane - The cli tool, effectively a wrapper exposing
the reusable logic found in crane-lib
* https://github.com/konveyor/crane-lib - Resuable library housing the core crane logic

First thing you need to do is ensure your workstation is set up to build go
projects. Currently Crane is built from source, although we plan to release
pre-build binaries in the future.

First, checkout crane into your `GOPATH`:

``
mkdir -p $GOPATH/konveyor
git clone https://github.com/konveyor/crane.git $GOPATH/konveyor/crane
cd $GOPATH/konveyor/crane
``

The following command will build the crane binary, which can be moved to a user's
PATH for regular usage:

`go build -o crane main.go`

## Workflow

Now that we have the Crane binary, we're ready to start migrating our app.

For this demo, we're going to use a very simple nginx example application. Be
sure you're logged into your source cluster and let's deploy the example:

`oc create -f ./nginx-example`


