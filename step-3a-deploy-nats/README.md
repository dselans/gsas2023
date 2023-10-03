# GSAS 2023 Workshop: Step 3a - Deploy NATS

### Why is there a 3a and 3b?

Both `etcd` and `NATS` are a good choice for a distributed key/value store. I
have used both in production and they both have their own PROs and CONs.

I started this workshop choosing `etcd` but decided to include `NATS` as well as
that is what I have been using most recently. NOTE: The code examples use `NATS`.

### Overview

In this step, we will deploy [`NATS`](https://nats.io), a popular, distributed,
all-in-one streaming, messaging, key/value and pubsub system to our K8S cluster.

We will use `NATS` specifically for it's [key/value](https://docs.nats.io/nats-concepts/jetstream/key-value-store)
functionality but there is nothing stopping you from also using it for its other
features.

### Why NATS?

NATS is extremely light-weight, very fast and has an extremely simple API.

It has *really* solid client libraries (in Go), is a breeze to install and has
a really good CLI tool that you would use for debug/troubleshooting/management.

While its libraries do not have support for the same distributed-system primitives
like `etcd` (distributed locks, leader election), there is nothing preventing you
from implementing those features yourself.

I have done this in the past, having written a library specifically for executing
bits of code as a cluster leader: [natty](https://github.com/streamdal/natty).

### Why Not Consolidate?

Some solutions, such as `NATS` or `Pulsar`, have support for messaging AND
streaming (and key/value, in the case of `NATS`) functionality. That is cool but
I think you should still have a separate message bus and a separate data bus
(and/or distributed store).

1. **Scaling is simpler**
   1. Just because your internal message bus needs to be bigger does not
      automatically mean that you need to also increase your data bus (which is
      almost certainly going to cost a lot more).
2. **Smaller blast radius**
   1. When issues arise (and they will) with one of the busses, you don't have
      to take everything down to perform a maintenance. Ie. If there's an issue
      with the ISB, your services that accept and produce data to the data bus
      won't be affected while you service the ISB.
3. **Separation of concerns**
   1. Event driven architecture is all about decoupling, reducing dependencies,
      and increasing autonomy. Not everyone will need to use the data bus so why
      expose it unnecessarily?

With that said, this is not a hard rule - do what makes sense for your use-case.

### Steps

1. We will use `helm` to deploy `NATS` to our local K8S cluster
   1. You should have helm installed from [step 1](../step-1-setup-kubernetes/README.md)
2. Add the `nats` helm repo
   1. `helm repo add nats https://nats-io.github.io/k8s/helm/charts/`
3. Update the repo
   1. `helm repo update`
4. Deploy `NATS` using the provided `values.yaml`
   1. `helm install -f values.yaml gsas nats/nats`
   2. NOTE: If this the second+ time you are installing NATS, you might get a
   `INSTALLATION FAILED: cannot re-use a name that is still in use` error. To
   get around it, run `helm upgrade --install -f values.yaml gsas nats/nats`
5. Verify that NATS is running
   1. `kubectl get pods | grep nats`; you should see something like:
   ```
   ❯ kubectl get pods | grep nats
    gsas-nats-0                                  2/2     Running            0               31m
    gsas-nats-box-5c6f7b945b-sq88h               1/1     Running            0               31m
   ```
6. Let's verify that it is working by doing key/value tests from a helper pod
   1. `kubectl exec -it deployment/gsas-nats-box /bin/sh`
   2. Inside the pod:
   ```
   ~ # nats kv ls
   No Key-Value buckets found
   
   ~ # nats kv add test
   Information for Key-Value Store Bucket test created 2023-10-02T21:56:12Z
    
   Configuration:
    
              Bucket Name: test
             History Kept: 1
            Values Stored: 0
       Backing Store Kind: JetStream
              Bucket Size: 0 B
      Maximum Bucket Size: unlimited
       Maximum Value Size: unlimited
              Maximum Age: unlimited
         JetStream Stream: KV_test
                  Storage: File
    
      Cluster Information:
    
                     Name:
                   Leader: gsas-nats-0
   
   ~ # nats kv put test foo bar
   bar
   
   ~ # nats kv get test foo
   test > foo revision: 1 created @ 02 Oct 23 21:56 UTC
    
   bar
   ~ # exit
   ```
   
### Done!

Alright! We have deployed `NATS` to our K8S cluster and verified that it is working.
We can now move on to [step-4-write-welcome-svc](../step-4-write-welcome-svc/README.md).
   
