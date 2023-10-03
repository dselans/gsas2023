# GSAS 2023 Workshop: Step 1

### Overview

In this step, we will install K8S (kubernetes) on our local machine.

We will use `minikube` for this, following the install steps in the 
[minikube docs](https://minikube.sigs.k8s.io/docs/start/).

### Why K8S?

To maximize the benefits of an EDA, it makes sense to have a modern, intelligent
orchestration layer that can be adapted to a variety of different use cases.

Since K8S (k + 8 letters (ubernete) + s) is the current de-facto standard for 
running containerized applications, we are going with that for this workshop.

But ultimately, it is perfectly fine if you do not wish to run and operate your
own K8S cluster - there are many managed K8S services out there and some other
hosted container platforms that might fit your needs.

### How does K8S benefit EDA?

One of the core principles in EDA is solid domain isolation. Because of this, 
you will probably have many services and they will all need to be deployed and
running somewhere.

Rather than spending time on building out a complex deployment pipeline, it is
much simpler to include a default K8S deployment yaml config in every service's
repo and establish a standard way to deploy your applications.

In addition, services will inevitably crash, get restarted, get redeployed, etc.
Optimally, every single one of these events should be handled automatically,
without human intervention. If you are on bare metal, juggling infra, juggling
deployment, juggling jenkins, etc. WHILE writing code and maintaining your EDA
architecture is a lot to handle at once and will ultimately slow you and your
team down.

With that said, there is no requirement for K8S (or any orchestration system) in
EDA, it will just make your and your devops/sre team's lives MUCH easier.

### Alternatives

K8S is fairly complex and has a fairly steep learning curve. If your devops or
SRE team does not already know K8S, requires learning it and do not have spare
cycles - it might be a better idea to use a hosted solution instead. Similarly,
if you are a small team and do not have a dedicated devops or SRE team, choosing
a hosted solution is probably the right choice.

Some alternatives to operating your own K8S:

1. [AWS EKS](https://aws.amazon.com/eks/) - AWS's managed K8S offering
   1. Real K8S, but you do not have to manage the actual K8S cluster
   2. You will still need to deploy EC2 for compute
   3. There is little to no vendor lock-in and will be able to move to your own 
   K8S cluster with minimal effort
2. [AWS ECS](https://aws.amazon.com/ecs/) - AWS's managed container orchestration
   1. Similar to K8S but not K8S; uses Docker Engine under the hood and requires
   _significantly_ less upfront container & orchestration knowledge.
   2. Unlike EKS - this is a proprietary solution and will lock you in on using
   AWS and due to this, migrating away from ECS will not be a simple process.
3. [Google's GKE](https://cloud.google.com/kubernetes-engine) - Google's managed
   K8S service
   1. Similar to EKS but on Google Cloud
   2. Same PROs and CONs as AWS; same little to no lock-in
   3. Originators of K8S, so you can expect a good experience
   
Practically all of the major cloud providers (AWS, GCP, Azure, DigitalOcean, 
Linode, Heroku, etc.) have some sort of managed K8S offering. You should probably
pick the vendor you already have other services with (and have a _really_ good
reason for why you decided to pick a different vendor instead).

### Steps

From minikube's docs for MacOS:

1. `curl -LO https://storage.googleapis.com/minikube/releases/latest/minikube-darwin-arm64`
2. `sudo install minikube-darwin-arm64 /usr/local/bin/minikube`
3. `minikube start`

Once installed, we will verify that everything is working as expected, by 
running a few commands:

<sub>To save on typing, run `alias k="minikube kubectl --"` and use `k` instead
of `minikube kubectl -- ...`.</sub>

1. `k get nodes` -- get all the nodes in our cluster
   1. Should see something like:
    ```
        ❯ k get nodes
        NAME       STATUS   ROLES           AGE   VERSION
        minikube   Ready    control-plane   70s   v1.27.4
    ```
2. `k get pods --all-namespaces` -- get all the pods in our cluster
   1. Should see something like: 
    ```
    ❯ k get pods --all-namespaces
    NAMESPACE     NAME                               READY   STATUS    RESTARTS        AGE
    kube-system   coredns-5d78c9869d-kbc55           1/1     Running   0               5m17s
    kube-system   etcd-minikube                      1/1     Running   0               5m29s
    kube-system   kube-apiserver-minikube            1/1     Running   0               5m29s
    kube-system   kube-controller-manager-minikube   1/1     Running   0               5m30s
    kube-system   kube-proxy-pkshb                   1/1     Running   0               5m17s
    kube-system   kube-scheduler-minikube            1/1     Running   0               5m29s
    kube-system   storage-provisioner                1/1     Running   1 (4m57s ago)   5m28s
    ```
3. `k get services --all-namespaces` -- get all the services in our cluster
    1. Should see something like:
    ```
    ❯ k get services --all-namespaces
    NAMESPACE     NAME         TYPE        CLUSTER-IP   EXTERNAL-IP   PORT(S)                  AGE
    default       kubernetes   ClusterIP   10.96.0.1    <none>        443/TCP                  7m15s
    kube-system   kube-dns     ClusterIP   10.96.0.10   <none>        53/UDP,53/TCP,9153/TCP   7m14s
    ```
4. `minikube dashboard` -- open the k8s dashboard in our browser
5. Now we will also install [`helm`](https://helm.sh/docs/intro/install/) which
we will use to install several components in later steps:
   1. `curl -fsSL -o get_helm.sh https://raw.githubusercontent.com/helm/helm/main/scripts/get-helm-3`
   2. `chmod 700 get_helm.sh`
   3. `./get_helm.sh`

Finally, we will verify that we can deploy a sample service to the cluster and 
are able to access it:

1. `k apply -f deploy.yaml` -- deploy a sample echo service
2. `k get pods | grep example` -- make sure you see the example svc in the pod list
3. `k port-forward service/example-service 8080:8080` -- open a port-forward for
localhost:8080 -> to the example service (**do NOT close/stop this command until
done accessing example service**)
4. Open a new terminal tab and test that we are able to access the service:: `curl http://localhost:8080/test`
    1. Should look something like: 
    ```
    ❯ curl http://localhost:8080/test
    Request served by example-service-deployment-d79c588db-gvkgd
   
    HTTP/1.1 GET /test
   
    Host: localhost:8080
    Accept: */*
    User-Agent: curl/7.79.1
    ```
   
Cool! Everything works! On top [step 2](../step-2-deploy-rabbit/README.md) we go!
