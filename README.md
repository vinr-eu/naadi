# naadi

Event notification broker in Go

## Decision record

1. Notifier will use API endpoint to push events in JSON format. JSON is chosen as it is being most widely used.
2. Receiver will use API endpoint to fetch messages in JSON format. Again reason mentioned in point 1.
3. Pull mechanism is chosen to avoid credential of receiver getting stored in the broker.
4. Pull mechanism is also chosen so that we don't need to wait for the receiver side processing while pushing.
