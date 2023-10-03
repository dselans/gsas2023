# welcome-svc-2

This is a Golang microservice that sends welcome notifications to new users.

This version includes an idempotency check in the event handler to deal with
duplicate events.

To run: `go run main.go -d`
To build src, build docker & push docker img: `make build && make docker/build && make docker/push`
