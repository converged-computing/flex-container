# Flex Kubernetes

> Testing fluxion with a container abstraction

We want to be able to schedule to containers that share the same resources, so I'm trying that with a fluxion graph.

## Design

### Overview

The Flux Framework "flux-sched" or fluxion project provides modular bindings in different languages for intelligent,
graph-based scheduling. When we extend fluxion to a tool or project that warrants logic of this type, we call this a flex!
Thus, the project here demonstrates flex-archspec, or using fluxion to match some system request to what we know in the archspec graph. E.g.,

> Do you have an x86 system with this compiler option?

This is a simple use case that doesn't perfectly reflect the OCI container use case, but we need to start somewhere! For this very basic setup we are going to:

1. Load the machines into a JSON Graph (called JGF).
2. Try doing a query against system metadata

There will eventually be a third component - a container image specification, for which we need to include somewhere here. I am starting simple!

### Concepts

From the above, the following definitions might be useful.

 - **[Flux Framework](https://flux-framework.org)**: a modular framework for putting together a workload manager. It is traditionally for HPC, but components have been used in other places (e.g., here, Kubernetes, etc). It is analogous to Kubernetes in that it is modular and used for running batch workloads.
 - **[fluxion](fluxion)**: refers to [flux-framework/flux-sched](https://github.com/flux-framework/flux-sched) and is the scheduler component or module of Flux Framework. There are bindings in several languages, and specifically the Go bindings (server at [flux-framework/flux-k8s](https://github.com/flux-framework/flux-k8s)) assemble into the project "fluence."
 - **flex** is an out of tree tool, plugin, or similar that uses fluxion to flexibly schedule or match some kind of graph-based resources. This project is an example of a flex!

## Usage

### Build

This demonstrates how to build the bindings. You will need to be in the VSCode developer container environment, or produce the same
on your host. 

```bash
make
```
The output is generated in bin:

```bash
$ ls bin/
flex-container
```

### Run

#### 1. Paths

Ensure you have your LD library path set to find flux sched (fluxion) libraries.

```bash
export LD_LIBRARY_PATH=/usr/lib:/opt/flux-sched/resource:/opt/flux-sched/resource/reapi/bindings:/opt/flux-sched/resource/libjobspec
```
This will be set in the VSCode container already.

#### 2. Nodes

I've generated nodes that have two containers at the top level under the cluster (as a quasi replacement for rack) except they are sharing the same nodes.
Here is how to run.

```bash
./bin/flex-container
```

Here is what works (and what doesn't).
Asking for one container with two nodes works:

```bash
 ./bin/flex-container ./examples/jobspec-1.yaml 
```
```console
This is the flex kubernetes prototype
[./examples/jobspec-1.yaml]
 Match policy: first
 Load format: JSON Graph Format (JGF)
Created flex resource graph &{%!s(*fluxcli.ReapiCtx=&{{{}}})}

✨️ Init context complete!
&{0xc0000140c0 first 0xc00010e0e0 1 map[] map[0:0xc000108210] map[]}
  💻️  Request: ./examples/jobspec-1.yaml
        Match: {"graph": {"nodes": [{"id": "2", "metadata": {"type": "node", "basename": "node", "name": "node0", "id": 0, "uniq_id": 2, "rank": -1, "exclusive": true, "unit": "", "size": 1, "paths": {"containment": "/tiny0/container0/node0"}}}, {"id": "3", "metadata": {"type": "node", "basename": "node", "name": "node1", "id": 1, "uniq_id": 3, "rank": -1, "exclusive": true, "unit": "", "size": 1, "paths": {"containment": "/tiny0/container0/node1"}}}, {"id": "1", "metadata": {"type": "container", "basename": "container", "name": "container0", "id": 0, "uniq_id": 1, "rank": -1, "exclusive": true, "unit": "", "size": 1, "paths": {"containment": "/tiny0/container0"}}}, {"id": "0", "metadata": {"type": "cluster", "basename": "tiny", "name": "tiny0", "id": 0, "uniq_id": 0, "rank": -1, "exclusive": false, "unit": "", "size": 1, "paths": {"containment": "/tiny0"}}}], "edges": [{"source": "1", "target": "2", "metadata": {"name": {"containment": "contains"}}}, {"source": "1", "target": "3", "metadata": {"name": {"containment": "contains"}}}, {"source": "0", "target": "1", "metadata": {"name": {"containment": "contains"}}}]}}

        Number: 1
```

Or just two containers:

```bash
 ./bin/flex-container ./examples/jobspec-2.yaml 
```

But not any other cases. I think the issue is the containment path can only hold one - need to talk about this with someone that knows fluxion well.

## License

HPCIC DevTools is distributed under the terms of the MIT license.
All new contributions must be made under this license.

See [LICENSE](https://github.com/converged-computing/cloud-select/blob/main/LICENSE),
[COPYRIGHT](https://github.com/converged-computing/cloud-select/blob/main/COPYRIGHT), and
[NOTICE](https://github.com/converged-computing/cloud-select/blob/main/NOTICE) for details.

SPDX-License-Identifier: (MIT)

LLNL-CODE- 842614
