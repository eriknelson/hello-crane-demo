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
custom plugins ([RFC document](https://github.com/konveyor/enhancements/pull/41)).

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

For ease of use, this project wraps these steps and their details within helper
scripts: `export.sh`, `transform.sh`, and `apply.sh`. This illustrates the
pipeline pattern, and is clearly demonstrated by the `pipeline.sh` script, what
will run each task in sequence.

### Export

The first thing you want to do with crane is to export everything within your
application namespace. Crane will discover all the API resources in-use using
the the k8s discovery API, and will export them out of the API server to your
local disk:

**export.sh**
```
#!/bin/bash
_dir="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
./bin/crane export --namespace=nginx-example --export-dir=export
```


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
to allow for total customization. The transform command accepts a "plugins"
directory argument. Each plugin is an executable that has a well defined stdin/out
interface, so you can easily write your own, or install and use those that have
been published and are generically useful. We're going to use several plugins to
help mutate the data into an agnostic state:

* **hc-whiteout**: This is a custom plugin written to whiteout (read: skip),
certain kinds of resources that we know aren't gonig to be needed in the target.
* **pvc**: Stripping PVCs of their environment specifics to ensure they can be
satisfied in an independent manner on the target side.
* **route**: This is an OpenShift specific plugin that handles removing host
specific information.
* **service**: Similar to the OpenShift plugin, but for k8s services.
* **skip-owned**: This plugin eliminates resources that are derivatives of their
owner. This is a very common pattern in Kubernetes, i.e. Deployments spawn Pods.
We only want to restore the master resource and allow it to recreate it's owned
Pods.
* **spec-ns**: This plugin illustrates the acceptance of plugin arguments. In this
case, I'm going to provide a *new* destination namespace where I want my resources
to be created, and this plugin will handle that.
* **status-removal**: Finally, we're going to strip all of the status from the
objects, since `Status` is really only relevant when a resource as been instantiated.

It's likely Crane 2 will ship with most of the shared functionality that we know
many people are likely going to need as part of a "Core" package of functionality
that will be available out of the box.

> NOTE: Much of the logic found in these plugins comes from Crane 1.x and the
> Velero plugins that were used to handle these types of corner cases, although
> sometimes there were some limitations around the data that we were able to
> use to make decisions. Here, we have full access to the input resources.

Think of the transform command as a function that accepts the set of exported
resources you discovered initially, plus a set of plugins that are applied to
each of those resources. Its output is going to be a diretory with a set of
"transform" files that describe the mutation that should be applied to the
original resources before their final import. These mutations are expressed with
the [JSONPatch](https://jsonpatch.com) format, are human readable, and easily
hackable.

Let's run the transform command to generate our transform files:

**transform.sh**

```
#!/bin/bash
_dir="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
source $_dir/var.sh # Source some configuration

./bin/crane transform \
  --export-dir=$_dir/export/resources \
  --plugin-dir=$_dir/bin/plugins \
  --transform-dir=$_dir/transform \
  --optional-flags="dest-namespace=$DEST_NAMESPACE"
```

Looking at the output directory, we see a directory structure organized by
namespace, with a set of the transform files to be applied. Taking a look at
the route transform:

```
[{"op":"remove","path":"/spec/host"},{"op":"remove","path":"/status", /* SNIP */}]
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

**transform.sh**
```
#!/bin/bash
_dir="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

./bin/crane apply \
  --export-dir=$_dir/export/resources \
  --transform-dir=$_dir/transform \
  --output-dir=$_dir/output
```

Comparing the resources found in the `export` directory to those found in the
`output` directory, you'll see resources skipped, along with all of the mutations
applied to the input  resources. The host has been stripped from the `Route`,
and the `Status` has beben removed, for example.

### Gitops Integration

At this point, Crane 2 has done its job, and we have an environment agnostic
bundle of k8s resources that we could use `kubectl apply -f ./output/nginx-example`
on to instantiate in whatever cluster that our `KUBECONFIG` points to. However,
let's take the opportunity to get this into a git repository, and then get that
integrated with a CD tool like Argo CD. This is how Crane 2 can help onboard
users into best practes using Gitops that unlocks the ability to be truly fluid
and portable.

In my cluster, I have OpenShift Gitops installed (OpenShift's flavor of ArgoCD),
and we have an `argo/` directory with a couple of resources inside of it. The first
is the destination namespace `nginx-example-foo` so we have a place to create our
objects, and the second is an Argo `Application`. This resource tells Argo about
my git repository and grants it access to be able to monitor my repo, and deploy
it to my desired location.

I'm going to create a github repo and commit the resources that Crane 2 output to
that repository in an `app` directory that the `Application` is configured to
look for.

Before pushing to that repository, I'm going to `oc apply -f argo` to create
that `Application`, and you'll see the `Application` get created in the Argo UI.
It will try to reconcile this repo, but of course, since there's nothing in the
repo yet, it's going to fail to deploy. However, upon pushing `nginx-example-foo`,
ArgoCD will recognize this, reconcile the contents of the resository, and you'll
see your app deployed into the target cluster.
