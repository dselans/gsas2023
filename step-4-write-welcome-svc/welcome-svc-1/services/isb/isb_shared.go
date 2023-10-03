package isb

import (
	"encoding/json"

	amqp "github.com/rabbitmq/amqp091-go"

	"github.com/dselans/welcome-svc/types"
)

// SharedConsumeFunc will receive rabbitmq messages on only one running instance of this service
func (i *ISB) SharedConsumeFunc(msg amqp.Delivery) error {
	if err := msg.Ack(false); err != nil {
		i.log.Errorf("Error acknowledging message: %s", err)
		return nil
	}

	i.log.Debugf("Received message: %s", msg.Body)

	// Try to decode it so we can figure out event type
	event := &types.Event{}

	if err := json.Unmarshal(msg.Body, event); err != nil {
		i.log.Errorf("Unable to unmarshal message: %s", err)
		return nil
	}

	switch event.Type {
	case types.SignupEventType:
		i.log.Debugf("Signup event data: %s", event.Data)

		signupEvent := &types.SignupEvent{}
		if err := json.Unmarshal(event.Data, signupEvent); err != nil {
			i.log.Errorf("Unable to unmarshal signup event: %s", err)
			return nil
		}

		i.log.Debugf("Forwarding signup event for '%s' to handler", signupEvent.Email)

		i.handleSignupEvent(signupEvent)
	default:
		i.log.Errorf("Unknown event type: %v", event.Type)
	}

	return nil
}

// handleSignupEvent v1
func (i *ISB) handleSignupEvent(event *types.SignupEvent) {
	i.log.Debugf("handling signup event for email '%s'", event.Email)

	i.log.Debugf("sending welcome email to '%s'", event.Email)
}
