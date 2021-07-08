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
projects. Currently Crane and its plugins are build from source, but in the
future we expect to distribute the binaries themselves, and have a mechanism
that will easily allow users to install community plugins, as well as their own
custom plugins.

To get started, you can simply run the `./prep.sh` script. This script will
create two directories: a `build/` dir with build artifacts, and a `bin/` directory
that contains the crane cli tool, as well as the plugns that we'll be using
for this demonstration.

Now that we have built Crane and our plugins, we'll want a demo application to
showcase the crane workflow.

For this demo, we're going to use a very simple nginx example application. Be
sure you're logged into your source cluster and let's deploy the example:

`oc create -f ./nginx-example`

## Workflow

Crane's workflow is phased, providing for significant transparency in the process.
These phases are `export`, `transform`, and `apply`.

You can think of these phases as an idempotent pipeline. The commands do
not alter their inputs, so that you may run a command, verify its output, and if
anything does not look correct, back up a single step and rerun the command
after making some tweaks. This idempotency is incredibly useful when performing
large scale migrations. It's directly something we've learned thanks to Crane 1.0.

> NOTE: By default, crane will use your active kube context. Be sure to configure
> your context such that you're authenticated with your source cluster before
> continuing.

## Export

The first thing you want to do with crane is to export everything within your
application namespace. Crane will discover all the API resources in-use using
the the k8s discovery API, and will export them out of the API server to your
local disk:

`./bin/crane export --namespace=nginx-example --export-dir=export`

You'll find two directories within the export dir. The first is a `failures/`
directory for you to be able to inspect any errors that occurred during the
export. The second is a `resources/` directory with all the k8s resources that
were discovered within the namespace.
