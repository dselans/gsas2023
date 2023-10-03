# GSAS 2023 Workshop: Step 5

### Overview

In this step, we will write our second service that will participate in our
architecture example.

This service has _no idea_ about any other services - as far as it is concerned,
it is the only service in the system. Just like the `welcome-svc`, it will
listen to the ISB for messages it is interested in but instead of sending a
welcome email for new users, it will create a charge in Stripe.

### Steps

1. Open the `billing-svc` directory in your preferred IDE
2. Look at that... the service looks identical to `welcome-svc` - and that is
**exactly** what we want. We have established a pattern!
3. Let's spot the difference though - it is in the handler code:
   1. Open up `services/isb/isb_shared.go` and look at `handleSignupEvent`
   2. While slightly different it is _extremely_ similar to the `welcome-svc`
   code - except, it already has the idempotency code.
   3. But there is an important detail here:
      1. We are storing our "idempotency" data in a different bucket than the `welcome-svc`
      2. And this is **BY DESIGN**
      3. We do not want to couple our services together - they should be able to
      operate completely independently of each other!
4. Let's test it all out!
   1. Deploy the service: `cd billing-svc && k apply -f deploy.yaml`
   2. Watch the logs of _both_ `welcome-svc` and `billing-svc` in two separate
   terminals: 
      ```
      TERM 1: k logs -f deployment/billing-svc-deployment
      TERM 2: k logs -f deployment/welcome-svc-deployment 
      ```
   3. Emit a `signup` event:
      ```
      plumber write rabbit \
      --exchange-name events \
      --routing-key events.signup.new \
      --input-file plumber-input.json
      ```
   4. Verify that both services have logged they received the event and handled it
5. Now, let's see what "eventual consistency" is all about:
   1. Let's pretend that the `welcome-svc` has crashed and can't come back
   up automatically: `k delete deployment welcome-svc-deployment`
   2. Emit another `signup` event:
      ```
      plumber write rabbit \
      --exchange-name events \
      --routing-key events.signup.new \
      --input-file plumber-input.json
      ```
   3. Watch the logs of `billing-service`, as expected, it handled the event as
   as it should
   4. Now let's imagine 30 minutes have passed and an SRE finally got `welcome-svc`
   to be able to start again: `cd step-4-write-welcome-svc/welcome-svc-2 && k apply -f deploy.yaml`
   5. Watch the logs of `k logs -f deployment/welcome-svc-deployment` - without us doing anything, it has
   handled the event that was emitted 30 minutes ago!
   6. That is eventual consistency in action!

### Done

And that's it! We have reached the end of our workshop. Let's reflect on what
we've done and what we've learned & leave with some wise words in 
[step 6](../step-6-whats-next/README.md).
