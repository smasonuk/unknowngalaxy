package comms

import (
	"bytes"
	"testing"
)

func TestMessageBus(t *testing.T) {
	bus := NewMessageBus()

	var probe1Got *Message

	bus.Subscribe("Earth", func(msg Message) {})
	bus.Subscribe("Probe1", func(msg Message) {
		probe1Got = &msg
	})

	payload := []byte{0x01, 0x02, 0x03}
	bus.Send("Earth", "Probe1", payload)

	// Message should not be delivered before Tick
	if probe1Got != nil {
		t.Fatal("Probe1 received message before Tick was called")
	}

	bus.Tick()

	// Message should be delivered after Tick
	if probe1Got == nil {
		t.Fatal("Probe1 did not receive message after Tick")
	}
	if probe1Got.SenderID != "Earth" {
		t.Errorf("expected SenderID %q, got %q", "Earth", probe1Got.SenderID)
	}
	if probe1Got.TargetID != "Probe1" {
		t.Errorf("expected TargetID %q, got %q", "Probe1", probe1Got.TargetID)
	}
	if !bytes.Equal(probe1Got.Payload, payload) {
		t.Errorf("expected payload %v, got %v", payload, probe1Got.Payload)
	}
}
