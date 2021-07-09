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

### Export

The first thing you want to do with crane is to export everything within your
application namespace. Crane will discover all the API resources in-use using
the the k8s discovery API, and will export them out of the API server to your
local disk:

`./bin/crane export --namespace=nginx-example --export-dir=export`

You'll find two directories within the export dir. The first is a `failures/`
directory for you to be able to inspect any errors that occurred during the
export. The second is a `resources/` directory with all the k8s resources that
were discovered within the namespace.

### Transform

Frequently when migrating workloads between one environment to another, you'll
encounter a need to change something about the resources before they're imported
into your target environment. A few examples include:

* Stripping the status information about a resources that isn't relevant after
the resource has been serialized out of a cluster.
* Adjusting resource quotas to fit your destination environment.
* Altering node selectors to match your new environment, if the node labels don't
match your source environment.
* Applying custom labels or annotation to your resources during the migration

A lot of these reasons are specific to your environment, so crane is designed
to allow to be totally customizable. The transform command accepts a "plugins"
directory argument. Each plugin is an executable that has a well defined stdin/out
interface, so you can easily write your own, or install and use those that have
been published and are generically useful. In our case, we're going to use
three plugins: an `openshift` plugin that will handle OpenShift specific details
like clearing the host on a `Route` resource so that it can be regenerated in
the target environment, a `whiteout-pods` plugin that will remove the pods from
the resources that will be imported (we want our higher level resources to recreate
the pods naturally in the environment, since they're spawned and owned by a
Deployment), and finally we have a `status-removal` that will strip the status
from the resources since it's not desired on import.

Think of the transform command as a function that accepts the set of exported
resources you discovered initially, plus a set of plugins that are applied to
each of those resources. It's output is going to be a diretory with a set of
"transform" files that describe the mutation that should be applied to the
original resources before their final import. These mutations are expressed with
the [JSONPatch](https://jsonpatch.com) format.

Let's run the transform command to generate our transform files:

`./bin/crane transform --export-dir=./export/resources --plugin-dir=./bin/plugins --transform-dir=./transform`

Looking at the output directory, we see a directory structure organized by
namespace, with a set of the transform files to be applied. Taking a look at
the route transform:

```
[{"op":"remove","path":"/spec/host"},{"op":"remove","path":"/status"}]
```

We can see a couple of mutations to be applied, as determined by the plugins.
First the host will be removed from the spec of the Route, which as we pointed
out before is important for an OpenShift route so the host is regenerated in the
new environment (which is going to be a new host). Next, it's going to remove
the status, derived from the `status-removal` plugin.

It's important to note that because all of these transforms are simple files on
disk, they can be tweaked and versioned themselves in a git repo. None of these
are destructive, so they can always be rerun should an output be unexpected
or an error arise. You can always inspect your inputs for ease of diagnosis.

### Apply

Finally, we want to apply the mutations that are described by our transform files
and generate the set of k8s resources that we'll ultimately want to deploy to
our target cluster. These resources could be imported into a gitops pipeline
and deployed via a tool like argo, or `kubectl create -f` directly.

`./bin/crane apply --export-dir=./export/resources --transform-dir=./transform --output-dir=./output`

Note the absence of a Pod, since our `whiteout-pods` plugins removed them from
the output resources. Similarly, the host has been snipped from the spec in
the openshift Route, and the status has been stripped from all the resources.

### Gitops Integration

TODO
