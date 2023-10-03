# GSAS 2023 Workshop: Step 3b - Deploy etcd

## NOTE: This step is optional - the workshop uses `NATS` instead of `etcd`

### Why is there a 3a and 3b?

Both `etcd` and `NATS` are a good choice for a distributed key/value store. I
have used both in production and they both have their own PROs and CONs.

I started this workshop choosing `etcd` but decided to include `NATS` as well as
that is what I have been using most recently. NOTE: The code examples use NATS.

### Overview

In this step, we will deploy [`etcd`](https://etcd.io), a popular, distributed
key/value store, to our K8S cluster.

`etcd` will be used by our services for storing state, configs and whatever 
other data that you'd want your application to have at startup, runtime, for 
recovery purposes and so on.

### Why etcd?

We choose `etcd` because it is distributed by default, has a simple API, is
*highly* battle tested (as it is used in Kubernetes for storage) and has many
[concurrency & synchronization primitives](https://etcd.io/docs/v3.5/dev-guide/api_concurrency_reference_v3/).
These are _extremely_ useful when building software that uses distributed system
concepts like _leader election_ and _distributed locks/mutexes_.

But, there are some downsides to `etcd` as well:

1. Max `2GB` total storage
   1. Can be configured to `8GB` but beware
2. OK not great read perf
   1. Under perfect conditions, possible to achieve 30K reads/sec
   2. Realistic: <10K reads/sec
   3. In comparison to Redis being able to achieve 300K reads/sec
3. OK not great write perf
   1. Under perfect conditions (<1KB key/value) - possible to achieve 30K writes/sec
   2. Realistic: <1K writes/sec

For this reason, `etcd` is best used for storing not-frequently accessed data,
such as state/config info that is used by applications at startup.

I would not recommend using `etcd` as a `cache` - there are better solutions for
that use-case.

### Example Usage of `etcd`

#### At Startup

1. Service is deployed to K8S
2. Service begins start up
3. Service reads its state from global `etcd` store
4. Service becomes available

#### Runtime w/ Distributed Locks

1. User request results in a long-running, involved task
2. Service keeps track of state of task and places a distributed lock on it
   (with a TTL) in global `etcd` store
3. If service crashes AND once lock expires, another instance of the service 
can pick up task where previous service left-off

#### Runtime w/ Leader Election

1. Have 3+ replicas of the same service, task comes in that requires coordination
between all replicas. Who gets the results? Who is the "leader"?
2. Answer: all services attempt to become leader using `etcd` - only one replica
will succeed and become leader.
3. Rest of replicas now know who the leader is and can send results to it.

#### Compare-and-Swap

1. Have 3+ replicas of the same service - all services get an *almost* 
simultaneous request that will cause the service to write state to `etcd`.
2. Problem: a replica of the service could overwrite "newer" state info that
was created by another service. All services should first check the age of the
state - but that is not an atomic operation...
3. Answer: `etcd` supports a `compare and swap` operation that allows the user
to only write/update the value IF the value is not different from the value the
user thinks it is.

### Steps

1. You should have `helm` already installed from step 1 (if not, go to
[step 1](../step-1-setup-kubernetes/README.md) and follow the instructions there)
2. Install & deploy `etcd`
   1. `helm install my-release oci://registry-1.docker.io/bitnamicharts/etcd`
3. Test that `etcd` works
   1. Export etcd root password: `export ETCD_ROOT_PASSWORD=$(kubectl get secret --namespace default my-release-etcd -o jsonpath="{.data.etcd-root-password}" | base64 -d)`
   2. Exec into one of etcd pods so we can run `etcdctl`: `minikube kubectl -- exec -it my-release-etcd-0 -- bash`
   3. In exec session, write something to etcd: `etcdctl --user root:$ETCD_ROOT_PASSWORD put /message Hello`
      1. You should see "OK"
   4. In `exec -it` session, read it back: `etcdctl --user root:$ETCD_ROOT_PASSWORD get /message`
      1. You should see "Hello"
   
### Done

Sweet - all of our dependencies are deployed and working. Let's double check 
that everything we deployed is still running and jot down all of the ports and
addresses that our applications will need.

1. `k get pods`
   1. It should look something like this:
   ```
   ❯ k get pods
   NAME                                         READY   STATUS    RESTARTS   AGE
   example-service-deployment-d79c588db-9njjs   1/1     Running   0          31h
   example-service-deployment-d79c588db-gvkgd   1/1     Running   0          31h
   my-release-etcd-0                            1/1     Running   0          14m
   my-release-etcd-client                       1/1     Running   0          10m
   rabbitmq-0                                   1/1     Running   0          8h
   ```
2. `k get services`
   1. Jot down the names of each service + ports - we will reference them as ENV
   vars in our services.
   2. It should look something like this:
   ```
   ❯ k get services
   NAME                       TYPE           CLUSTER-IP       EXTERNAL-IP   PORT(S)             AGE
   example-service            ClusterIP      10.101.235.30    <none>        8080/TCP            31h
   kubernetes                 ClusterIP      10.96.0.1        <none>        443/TCP             32h
   my-release-etcd            ClusterIP      10.106.149.159   <none>        2379/TCP,2380/TCP   14m
   my-release-etcd-headless   ClusterIP      None             <none>        2379/TCP,2380/TCP   14m
   rabbitmq-management        LoadBalancer   10.103.130.239   <pending>     8888:30797/TCP      8h
   rabbitmq-service           ClusterIP      10.111.145.196   <none>        5672/TCP            8h
   ```

OK, everything looks good. Let's move on to [step 5](../step-4-write-welcome-svc/README.md)
where we will write our first service!
