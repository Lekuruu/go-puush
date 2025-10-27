package email

// NoopEmail is a development helper that discards all outbound email
type NoopEmail struct {
	from string
}

// NewNoopEmail returns a no-op email implementation
func NewNoopEmail(from string) Email {
	return &NoopEmail{from: from}
}

// FromAddress returns the configured sender address
func (n *NoopEmail) FromAddress() string {
	return n.from
}

// Setup does nothing, it's just there to satisfy the Email interface
func (n *NoopEmail) Setup() error {
	return nil
}

// Send validates the message then discards it
func (n *NoopEmail) Send(message *Message) error {
	if err := message.Validate(); err != nil {
		return err
	}

	// TODO: Logging?
	return nil
}
