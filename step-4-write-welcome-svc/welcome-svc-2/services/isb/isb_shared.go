package isb

import (
	"context"
	"crypto/sha1"
	"encoding/json"
	"fmt"

	"github.com/nats-io/nats.go"
	"github.com/pkg/errors"
	amqp "github.com/rabbitmq/amqp091-go"

	"github.com/dselans/welcome-svc/types"
)

const (
	BucketName = "welcome-svc"
)

// SharedConsumeFunc will receive rabbitmq messages on only one running instance of this service
func (i *ISB) SharedConsumeFunc(msg amqp.Delivery) error {
	if err := msg.Ack(false); err != nil {
		i.log.Errorf("Error acknowledging message: %s", err)
		return nil
	}

	//i.log.Infof("Received message: %s", msg.Body)

	// Try to decode it so we can figure out event type
	event := &types.Event{}

	if err := json.Unmarshal(msg.Body, event); err != nil {
		i.log.Errorf("Unable to unmarshal message: %s", err)
		return nil
	}

	switch event.Type {
	case types.SignupEventType:
		signupEvent := &types.SignupEvent{}
		if err := json.Unmarshal(event.Data, signupEvent); err != nil {
			i.log.Errorf("Unable to unmarshal signup event: %s", err)
			return nil
		}

		i.log.Debugf("Forwarding signup event for '%s' to handler", signupEvent.Email)

		i.handleSignupEvent(signupEvent, event.Data)
	default:
		i.log.Errorf("Unknown event type: %v", event.Type)
	}

	return nil
}

// handleSignupEvent w/ idempotence
func (i *ISB) handleSignupEvent(event *types.SignupEvent, eventData []byte) {
	i.log.Debugf("handling signup event for email '%s'", event.Email)

	isKnown, err := i.knownEmail(event.Email, eventData)
	if err != nil {
		// Raise alarm - something went wrong!
		i.log.Errorf("Unable to determine if email '%s' is known: %s", event.Email, err)
		return
	}

	if isKnown {
		i.log.Debugf("email '%s' is known - ignoring", event.Email)
	} else {
		i.log.Debugf("email '%s' is unknown - sending welcome email", event.Email)
	}
}

func (i *ISB) knownEmail(email string, eventData []byte) (bool, error) {
	emailSha := fmt.Sprintf("%x", sha1.Sum([]byte(email)))

	_, err := i.NATS.Get(context.Background(), BucketName, emailSha)
	// Ran into an error with NATS - abort + raise alarm!
	if err != nil && err != nats.ErrKeyNotFound {
		return false, errors.Wrapf(err, "unable to get email '%s' from k/v store", email)
	}

	known := true

	// Key does not exist - create it but still tell caller that email is unknown
	if err != nil && err == nats.ErrKeyNotFound {
		i.log.Debugf("first time seeing email '%s' - creating entry in k/v", email)

		if err := i.NATS.Put(context.Background(), BucketName, emailSha, eventData); err != nil {
			return false, errors.Wrapf(err, "unable to put email '%s' into k/v store", email)
		}

		known = false
	}

	return known, nil
}
