GSAS 2023 Workshop Materials
============================
This repo contains the materials used in the workshop "Designing a Modern Distributed
System Using EDA" at [GSAS 2023]().

![img.png](img.png)

In this workshop, we will build a distributed system that uses an event-driven
architecture (EDA). We will deploy several components, create a couple of
services and demonstrate how it all works together. We will also bulletproof
our services to be more resilient to failure and be more EDA-friendly.

The presentation associated with the workshop can be found 
[here](https://docs.google.com/presentation/d/189EN83MUDw2KvxkDdlrFFWQcAIJKd9NL/edit?usp=sharing&ouid=116843161829560983453&rtpof=true&sd=true)

The `step-*` directories in this repo correspond to the hands-on steps in the
workshop. Each step contains a `README.md` that contains pertinent information
for that specific step. Start with [step 1](./step-1-setup-kubernetes/README.md)!

### Prerequisites

NOTE: This entire workshop has been written and tested on MacOS. It should work
on Linux and Windows with some _minor_ modifications, but no promises.

You will need to have installed:

- Basic understanding of [Docker](https://docs.docker.com/get-docker/)
- An IDE of some sort (VSCode, IntelliJ, etc.)
- Being comfortable operating in a terminal, working in multiple term tabs, etc.
