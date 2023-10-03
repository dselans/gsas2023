# GSAS 2023 Workshop: Step 6

### Overview

OK, so what's next? To review:

1. We have deployed K8S, RabbitMQ, NATS and two services "welcome-svc" and "billing-svc"
2. We know that idempotence is important and we have implemented it in both services
3. We know that eventual consistency is easily achievable by just having a sane
messaging foundation

So what's next?

### What's next?

There are still a lot of things that _should_ be done for us to be able to 
sleep easy at night. Some things you might want to explore on your own time:

1. Observability is *extremely* important, especially when dealing with event
driven systems. Use an application performance monitoring solution, use debug
tools like [`streamdal`](https://github.com/streamdal/streamdal) and
[`plumber`](https://github.com/streamdal/plumber).
2. Protect your services against DoS attacks - add a caching layer to (most)
network calls - and cache the response data in memory for 100ms + auto-reap old
cache data.
3. Learn about the saga pattern to deal with distributed transactions.
4. Implement a dead-letter strategy to deal with bad or incorrectly handled events.
5. Establish a common schema - the stricter, the better! Use protobuf!
6. Establish an event versioning strategy - diffing events is a pain - comparing
event version strings is _infinitely_ easier.
7. Make sure that your service shutdown behaves correctly, as in, it properly
stops consuming messages from the ISB and drains any in-flight messages.
8. Establish a standard service template - the more consistent your services are,
the faster your org will be able to move!
9. Establish _excellent_ monitoring for your global components - busses, K/V
stores, caches, etc. Your entire EDA architecture depends on these.
10. Load test your systems - due to the complexity of EDA, it is much harder to
identify potential bottlenecks by just looking at diagrams or code.

### Done

🎉 You have completed the workshop! 🎉

The setup and "scenarios" in this workshop are not "theoretical" - I have 
personally worked with several, almost identical platforms and software 
architectures. This architecture can _easily_ scale to hundreds of services
and not have it all devolve into complete chaos.

Some last moment notes:

1. The above `welcome-svc` and `billing-svc` are both based on a Go code template
that our company has used to power ~30 services in production for the last two
years, maintained by a team of 3 backend engineers.
2. The services are usually handling on avg ~5,000 internal, protobuf encoded
events/s.
3. Event driven is not for everyone - but it is a common enough eventuality for 
companies that have exhausted traditional scaling techniques. This workshop
has aimed to prepare you for that scenario.

If you have any questions or comments, please feel free to ping me at:
daniel [at] streamdal.com

Thank you!
